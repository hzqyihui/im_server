package grpc_base

import (
	rpc_controller "IM_Server/grpc/controller"
	"IM_Server/grpc/proto_service"

	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	Address = "0.0.0.0:9090"
)

func StartRpc() {

	listen, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// 服务注册
	proto_service.RegisterUserRpcServer(s, &rpc_controller.UserRpcController{})

	log.Println("gRPC listen on " + Address)

	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
