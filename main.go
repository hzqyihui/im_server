package main

import (
	grpc_base "IM_Server/grpc"
	"IM_Server/model"
	"github.com/joho/godotenv"
)

func main() {
	// 从本地读取环境变量
	godotenv.Load()
	model.Init()
	grpc_base.StartRpc()

	//epoll
	//epoll_server.StartEpoll()
}
