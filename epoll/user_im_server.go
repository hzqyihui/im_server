package epoll_server

import (
	"IM_Server/model"
	"IM_Server/util"
	"encoding/json"
	"errors"
	"fmt"
)

/**
用户上线：
1.修改online字段
*/
func UserIMOnline(request IMMessage) model.User {

	tx := model.DB.Begin()
	userInfo := model.User{}
	tx.Where(map[string]interface{}{"project_id": request.ProjectId, "project_uid": request.ProjectUid}).Find(&userInfo)

	//上线
	userInfo.IsOnline = 1

	if err := tx.Save(&userInfo).Error; err != nil {
		fmt.Println("上线更新失败")
	}
	tx.Commit()

	return userInfo
}

//保存消息记录
func UserIMSaveMessage(request IMMessage, isSend bool) error {

	sendUid, _ := model.GetUidByProjectInfo(request.ProjectUid, request.ProjectId)
	sendToUid, _ := model.GetUidByProjectInfo(request.ToProjectUid, request.ToProjectId)

	messageModel := model.Message{
		Uid:         sendUid,
		GroupId:     request.GroupId,
		ToUid:       sendToUid,
		Message:     request.Data,
		MessageType: request.Type,
		Time:        request.Time,
		IsSend:      util.BoolToInt(isSend),
	}

	fmt.Println("插入消息数据：", messageModel)
	if err := model.DB.Create(&messageModel).Error; err != nil {
		return errors.New("UserIMSaveMessageError")
	}

	return nil
}

//发送离线消息
func UserIMSendOfflineMessage(request IMMessage, s *ServerConn) error {

	curUid, _ := model.GetUidByProjectInfo(request.ProjectUid, request.ProjectId)

	offLineDbMessage := []model.Message{}
	model.DB.Where(map[string]interface{}{"to_uid": curUid, "is_send": 0}).Find(&offLineDbMessage)

	for _, v := range offLineDbMessage {
		sendProjectId, sendProjectUid, _ := model.GetProjectByUidInfo(v.Uid)
		imMessage := IMMessage{
			ProjectId:    sendProjectId,
			ProjectUid:   sendProjectUid,
			Time:         v.Time,
			Data:         v.Message,
			Type:         1,
			ToProjectId:  request.ToProjectId,
			ToProjectUid: request.ToProjectUid,
		}
		//发送离线消息
		responseJson, _ := json.Marshal(imMessage)
		s.Write([]byte(responseJson))
		//更新状态
		v.IsSend = 1
		if err := model.DB.Save(&v).Error; err != nil {
			return err
		}
	}
	return nil
}
