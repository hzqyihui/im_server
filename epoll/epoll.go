//go:build linux
// +build linux

package epoll_server

import (
	"IM_Server/model"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"
)

func NewEpollM() *EpollM {
	return &EpollM{
		conn: make(map[int]*ServerConn),
	}
}

type EpollM struct {
	conn map[int]*ServerConn

	socketFd int //监听socket的fd
	epollFd  int //epoll的fd
}

//关闭所有的链接
func (e *EpollM) Close() {
	syscall.Close(e.socketFd)
	syscall.Close(e.epollFd)
	for _, con := range e.conn {
		con.Close()
	}
}

//获取符合条件的用户
func (e *EpollM) GetUserConn(projectId int, projectUid int) *ServerConn {
	for _, v := range e.conn {
		if v.userInfo.ProjectId == projectId && v.userInfo.ProjectUid == projectUid {
			return v
		}
	}
	return nil

}

//获取一个链接
func (e *EpollM) GetConn(fd int) *ServerConn {
	return e.conn[fd]
}

//添加一个链接
func (e *EpollM) AddConn(conn *ServerConn) {
	e.conn[conn.fd] = conn
}

//删除一个链接
func (e *EpollM) DelConn(fd int) {
	delete(e.conn, fd)
}

//开启监听
func (e *EpollM) Listen(ipAddr string, port int) error {
	//使用系统调用,打开一个socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return err
	}
	//设置端口重复使用
	syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

	//ip地址转换
	var addr [4]byte
	copy(addr[:], net.ParseIP(ipAddr).To4())
	net.ParseIP(ipAddr).To4()
	err = syscall.Bind(fd, &syscall.SockaddrInet4{
		Port: port,
		Addr: addr,
	})
	if err != nil {
		return err
	}

	//开启监听
	err = syscall.Listen(fd, 10)
	if err != nil {
		return err
	}
	e.socketFd = fd
	return nil
}

//在死循环中等待client发来的链接
func (e *EpollM) Accept() error {
	for {
		nfd, _, err := syscall.Accept(e.socketFd)
		if err != nil {
			return err
		}
		err = e.EpollAddEvent(nfd)
		if err != nil {
			return nil
		}
		e.AddConn(&ServerConn{
			fd:     nfd,
			epollM: e,
		})

	}
}

//关闭指定的链接
func (e *EpollM) CloseConn(fd int) error {
	conn := e.GetConn(fd)
	if conn == nil {
		return nil
	}
	if err := e.EpollRemoveEvent(fd); err != nil {
		return err
	}
	conn.Close()
	e.DelConn(fd)
	println("关闭客户端连接：", fd)

	//修改在线状态
	if conn.userInfo.ID > 0 {
		dbUser := conn.userInfo
		dbUser.IsOnline = 0
		if err := model.DB.Save(dbUser).Error; err != nil {
			return err
		}
	}
	return nil
}

//创建一个epoll
func (e *EpollM) CreateEpoll() error {
	//通过系统调用,创建一个epoll
	fd, err := syscall.EpollCreate(1)
	if err != nil {
		return err
	}
	e.epollFd = fd
	return nil
}

//处理epoll
func (e *EpollM) HandlerEpoll() error {
	events := make([]syscall.EpollEvent, 100)
	//在死循环中处理epoll
	for {
		//fmt.Println("epoll——run")
		//msec -1,会一直阻塞,直到有事件可以处理才会返回, n 事件个数
		//这里epoll 不支持接受accept 事件，有连接来时 还是一直阻塞的
		n, err := syscall.EpollWait(e.epollFd, events, -1)
		fmt.Println("epoll_event——run:", n)
		if err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			//先在map中是否有这个链接
			conn := e.GetConn(int(events[i].Fd))
			if conn == nil { //没有这个链接,忽略
				continue
			}
			if events[i].Events&syscall.EPOLLHUP == syscall.EPOLLHUP || events[i].Events&syscall.EPOLLERR == syscall.EPOLLERR {
				//断开||出错
				if err := e.CloseConn(int(events[i].Fd)); err != nil {
					return err
				}
			} else if events[i].Events == syscall.EPOLLIN {
				//可读事件
				conn.Read()
			}
		}
	}
}

//向epoll中加事件
func (e *EpollM) EpollAddEvent(fd int) error {
	return syscall.EpollCtl(e.epollFd, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(fd),
		Pad:    0,
	})
}

//从epoll中删除事件
func (e *EpollM) EpollRemoveEvent(fd int) error {
	return syscall.EpollCtl(e.epollFd, syscall.EPOLL_CTL_DEL, fd, nil)
}

type ServerConn struct {
	epollM   *EpollM
	fd       int
	userInfo model.User //im用户id    校验通过，是正常用户
}

//3 心跳  ProjectId  ProjectUid
//2 回执消息  （tudo 后期加个uuid来标识唯一消息）
//1 消息
type IMMessage struct {
	ProjectId  int `json:"ProjectId,omitempty"`
	ProjectUid int `json:"project_uid,omitempty"`
	// Uid int  //为0是服务器发送的
	Time         int    `json:"time,omitempty"` //时间戳
	Data         string `json:"data,omitempty"` //数据
	Type         int    `json:"type,omitempty"` //1普通消息2回执消息,表示已经收到 3心跳[ProjectId  ProjectUid]4.错误提示
	ToProjectId  int    `json:"to_project_id,omitempty"`
	ToProjectUid int    `json:"to_project_uid,omitempty"`
	GroupId      int    `json:"group_id,omitempty"`
}

func (s *ServerConn) PreHandleMessage(data []byte) (IMMessage, error) {
	request := IMMessage{}

	jsonDecodeErr := json.Unmarshal(data, &request)
	if len(data) == 0 || jsonDecodeErr != nil {
		return request, errors.New("解析失败:" + jsonDecodeErr.Error())
	}

	//必须
	if request.ProjectUid < 0 || request.ProjectId < 0 {
		return request, errors.New("缺少必要身份信息参数")
	}

	//补充时间
	request.Time = int(time.Now().Unix())

	return request, nil
}

//读取数据
func (s *ServerConn) Read() {
	data := make([]byte, 10000)
	//通过系统调用,读取数据,n是读到的长度
	n, err := syscall.Read(s.fd, data)
	if n == 0 || err != nil {
		//如果n=0，那这个链接就是客户端异常关闭了
		s.epollM.CloseConn(s.fd)
		return
	}

	//读取消息
	fmt.Printf("%d say: %s \n", s.fd, data)

	//预处理数据 必须符合requestMessage 格式
	//data[:] 注意这个地方要这么写，不然解析不出来,末尾会多个结束符
	requestMessage, requestMessageErr := s.PreHandleMessage(data[:n])
	if requestMessageErr != nil {
		fmt.Println("预处理数据出错：", requestMessageErr.Error())
		return
	}

	//上线操作
	if s.userInfo == (model.User{}) {
		userInfo := UserIMOnline(requestMessage)
		s.userInfo = userInfo

		fmt.Println(s.userInfo.Name, " 上线成功!")
		//上线成功后发送离线消息
		offlineErr := UserIMSendOfflineMessage(requestMessage, s)
		if offlineErr != nil {
			fmt.Println("发送离线消息失败：", offlineErr.Error())
		}

	}

	//读取消息
	if s.userInfo != (model.User{}) && requestMessage.Type == 1 {

		//发送消息socket
		isSend := false
		ToConn := s.epollM.GetUserConn(requestMessage.ToProjectId, requestMessage.ToProjectUid)
		if ToConn != nil {
			fmt.Println(s.userInfo.Name, " 给 ", ToConn.userInfo.Name, " 发送消息：", requestMessage.Data)
			responseJson, _ := json.Marshal(requestMessage)
			ToConn.Write([]byte(responseJson))
			isSend = true
		}
		//存下消息
		UserIMSaveMessage(requestMessage, isSend)

	}

	//回执消息
	responseStruct := IMMessage{
		ProjectUid: 0,
		Time:       int(time.Now().Unix()),
		Type:       2,
	}
	responseJson, _ := json.Marshal(responseStruct)
	s.Write([]byte(responseJson))
	fmt.Println("服务器回执消息：", responseStruct)

}

//向这个链接中写数据
func (s *ServerConn) Write(data []byte) {
	_, err := syscall.Write(s.fd, data)
	if err != nil {
		fmt.Printf("fd %d write error:%s\n", s.fd, err.Error())
	}
}

//关闭这个链接
func (s *ServerConn) Close() {
	err := syscall.Close(s.fd)
	if err != nil {
		fmt.Printf("fd%d close error:%s\n", s.fd, err.Error())
	}
}

func StartEpoll() {
	fmt.Println("epoll 服务启动")
	epollM := NewEpollM()
	epollIp := os.Getenv("EPOLL_IP")
	epollPort, _ := strconv.Atoi(os.Getenv("EPOLL_PORT"))
	//开启监听
	err := epollM.Listen(epollIp, epollPort)
	if err != nil {
		panic(err)
	}

	//创建epoll
	err = epollM.CreateEpoll()
	if err != nil {
		panic(err)
	}

	//异步处理epoll
	go func() {
		err := epollM.HandlerEpoll()
		epollM.Close()
		panic(err)
	}()

	//等待client的连接
	err = epollM.Accept()
	//todo 心跳检测

	epollM.Close()
	panic(err)
}
