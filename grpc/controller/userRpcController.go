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

func (c *UserRpcController) UserEdit(ctx context.Context, in *proto_service.User) (*proto_service.Response, error) {
	fmt.Println("收到一个 UserEdit 请求，请求参数：", in)

	response := &proto_service.Response{Ok: true, Msg: ""}
	if err := service.UserEdit(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}

func (c *UserRpcController) GroupAdd(ctx context.Context, in *proto_service.Group) (*proto_service.Response, error) {
	fmt.Println("收到一个 GroupAdd 请求，请求参数：", in)

	response := &proto_service.Response{Ok: true, Msg: ""}
	if err := service.GroupEdit(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}

func (c *UserRpcController) GroupEdit(ctx context.Context, in *proto_service.Group) (*proto_service.Response, error) {
	fmt.Println("收到一个 GroupEdit 请求，请求参数：", in)

	response := &proto_service.Response{Ok: true, Msg: ""}
	if err := service.GroupEdit(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}

func (c *UserRpcController) UserGroupAdd(ctx context.Context, in *proto_service.UserGroup) (*proto_service.Response, error) {
	fmt.Println("收到一个 UserGroupAdd 请求，请求参数：", in)

	response := &proto_service.Response{Ok: true, Msg: ""}
	if err := service.UserGroupEdit(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}

func (c *UserRpcController) UserGroupEdit(ctx context.Context, in *proto_service.UserGroup) (*proto_service.Response, error) {
	fmt.Println("收到一个 UserGroupEdit 请求，请求参数：", in)

	response := &proto_service.Response{Ok: true, Msg: ""}
	if err := service.UserGroupEdit(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}
