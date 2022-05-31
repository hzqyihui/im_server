package rpc_controller

import (
	"IM_Server/grpc/proto_service"
	"IM_Server/service"
	"fmt"
	"golang.org/x/net/context"
)

type GroupRpcController struct {
	proto_service.IMRpcServer
}

func (c *UserRpcController) GroupAdd(ctx context.Context, in *proto_service.Group) (*proto_service.HandleResponse, error) {
	fmt.Println("收到一个 GroupAdd 请求，请求参数：", in)

	response := &proto_service.HandleResponse{Ok: true, Msg: ""}
	if err := service.GroupEdit(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}

func (c *UserRpcController) GroupEdit(ctx context.Context, in *proto_service.Group) (*proto_service.HandleResponse, error) {
	fmt.Println("收到一个 GroupEdit 请求，请求参数：", in)

	response := &proto_service.HandleResponse{Ok: true, Msg: ""}
	if err := service.GroupEdit(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}

func (c *UserRpcController) GroupDel(ctx context.Context, in *proto_service.Group) (*proto_service.HandleResponse, error) {
	fmt.Println("收到一个 GroupDel 请求，请求参数：", in)

	response := &proto_service.HandleResponse{Ok: true, Msg: ""}
	if err := service.GroupDel(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}
