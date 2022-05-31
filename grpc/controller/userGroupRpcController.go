package rpc_controller

import (
	"IM_Server/grpc/proto_service"
	"IM_Server/service"
	"fmt"
	"golang.org/x/net/context"
)

type UserGroupRpcController struct {
	proto_service.IMRpcServer
}

func (c *UserGroupRpcController) UserGroupEdit(ctx context.Context, in *proto_service.UserGroup) (*proto_service.HandleResponse, error) {
	fmt.Println("收到一个 UserGroupEdit 请求，请求参数：", in)

	response := &proto_service.HandleResponse{Ok: true, Msg: ""}
	if err := service.UserGroupEdit(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}
