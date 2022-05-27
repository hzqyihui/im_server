package grpc_base

import (
	rpc_controller "IM_Server/grpc/controller"
	"IM_Server/grpc/proto_service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func StartRpc() {

	rpcAddress := os.Getenv("RPC_ADDRESS")
	listen, err := net.Listen("tcp", rpcAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// 服务注册
	proto_service.RegisterUserRpcServer(s, &rpc_controller.UserRpcController{})

	log.Println("gRPC listen on " + rpcAddress)

	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
