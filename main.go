package main

import (
	"IM_Server/cache"
	epoll_server "IM_Server/epoll"
	grpc_base "IM_Server/grpc"
	"IM_Server/model"
	"github.com/joho/godotenv"
)

func main() {
	// 从本地读取环境变量
	godotenv.Load()
	model.Init()
	//开启redis
	cache.Redis()

	//grpc
	go grpc_base.StartRpc()

	//epoll
	go epoll_server.StartEpoll()

	select {}
}
