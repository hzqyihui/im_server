package service

import (
	"IM_Server/grpc/proto_service"
	"IM_Server/model"
	"errors"
	"fmt"
)

func UserEdit(userRpc *proto_service.User) error {

	tx := model.DB.Begin()
	userModel := model.User{
		Name:       userRpc.Name,
		Avatar:     *userRpc.Avatar,
		ProjectId:  int(userRpc.ProjectId),
		ProjectUid: int(userRpc.Uid),
	}
	fmt.Println(userModel)
	dbUser := &model.User{}
	dbUserCount := 0
	tx.Where(model.User{ProjectId: userModel.ProjectId, ProjectUid: userModel.ProjectUid}).Find(dbUser).Count(&dbUserCount)
	if dbUserCount > 0 {
		//编辑
		dbUser.Name = userModel.Name
		dbUser.Avatar = userModel.Avatar
		fmt.Println("编辑", dbUser)
		if err := tx.Save(&dbUser).Error; err != nil {
			return errors.New("更新失败")
		}
	} else {
		//新增
		fmt.Println("新增", userModel)
		if err := tx.Create(&userModel).Error; err != nil {
			return errors.New("注册失败")
		}
	}
	tx.Commit()
	return nil

}
