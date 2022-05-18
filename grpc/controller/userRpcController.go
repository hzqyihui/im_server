package rpc_controller

import (
	"IM_Server/grpc/proto_service"
	"IM_Server/service"
	"fmt"
	"golang.org/x/net/context"
)

type UserRpcController struct {
	proto_service.UserRpcServer
}

func (c *UserRpcController) UserAdd(ctx context.Context, in *proto_service.User) (*proto_service.Response, error) {
	fmt.Println("收到一个 UserAdd 请求，请求参数：", in)

	response := &proto_service.Response{Ok: true, Msg: ""}
	if err := service.UserEdit(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}
