package rpc_controller

import (
	"IM_Server/grpc/proto_service"
	"IM_Server/service"
	"fmt"
	"golang.org/x/net/context"
)

type MessageRpcController struct {
	proto_service.IMRpcServer
}

func (c *UserGroupRpcController) QueryMessageList(ctx context.Context, in *proto_service.QueryMessageReq) (*proto_service.IMMessages, error) {
	fmt.Println("收到一个 QueryMessageList 请求，请求参数：", in)

	imMessages := service.QueryMessageList(in)

	return &proto_service.IMMessages{
		Message: imMessages,
	}, nil
}
