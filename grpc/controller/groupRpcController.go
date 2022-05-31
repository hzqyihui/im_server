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

func (c *GroupRpcController) GroupAdd(ctx context.Context, in *proto_service.Group) (*proto_service.HandleResponse, error) {
	fmt.Println("收到一个 GroupAdd 请求，请求参数：", in)

	response := &proto_service.HandleResponse{Ok: true, Msg: ""}
	if err := service.GroupEdit(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}

func (c *GroupRpcController) GroupEdit(ctx context.Context, in *proto_service.Group) (*proto_service.HandleResponse, error) {
	fmt.Println("收到一个 GroupEdit 请求，请求参数：", in)

	response := &proto_service.HandleResponse{Ok: true, Msg: ""}
	if err := service.GroupEdit(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}

func (c *GroupRpcController) GroupDel(ctx context.Context, in *proto_service.Group) (*proto_service.HandleResponse, error) {
	fmt.Println("收到一个 GroupDel 请求，请求参数：", in)

	response := &proto_service.HandleResponse{Ok: true, Msg: ""}
	if err := service.GroupDel(in); err != nil {
		response.Ok = false
		response.Msg = err.Error()
	}

	return response, nil
}

//查询方法
func (c *GroupRpcController) QueryGroupList(ctx context.Context, in *proto_service.QueryCommonReq) (*proto_service.IMGroups, error) {
	fmt.Println("收到一个 QueryGroupList 请求，请求参数：", in)

	imGroups := service.QueryGroupList(in)

	return &proto_service.IMGroups{
		Group: imGroups,
	}, nil
}
