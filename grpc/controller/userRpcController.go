package rpc_controller

import (
	"IM_Server/grpc/proto_service"
	"IM_Server/service"
	"fmt"
	"golang.org/x/net/context"
)

type UserRpcController struct {
	proto_service.IMRpcServer
}

func (c *UserRpcController) UserAdd(ctx context.Context, in *proto_service.User) (*proto_service.HandleResponse, error) {
	fmt.Println("收到一个 UserAdd 请求，请求参数：", in)

	response := &proto_service.HandleResponse{Ok: true, Msg: ""}
	if err := service.UserEdit(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}

func (c *UserRpcController) UserEdit(ctx context.Context, in *proto_service.User) (*proto_service.HandleResponse, error) {
	fmt.Println("收到一个 UserEdit 请求，请求参数：", in)

	response := &proto_service.HandleResponse{Ok: true, Msg: ""}
	if err := service.UserEdit(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}
func (c *UserRpcController) UserDel(ctx context.Context, in *proto_service.User) (*proto_service.HandleResponse, error) {
	fmt.Println("收到一个 UserDel 请求，请求参数：", in)

	response := &proto_service.HandleResponse{Ok: true, Msg: ""}
	if err := service.UserDel(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}

//查询方法
func (c *UserRpcController) QueryUserList(ctx context.Context, in *proto_service.QueryCommonReq) (*proto_service.IMUsers, error) {
	fmt.Println("收到一个 QueryUserList 请求，请求参数：", in)

	imUsers := service.QueryUserList(in)

	return &proto_service.IMUsers{
		User: imUsers,
	}, nil
}
