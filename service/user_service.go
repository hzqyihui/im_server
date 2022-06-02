package service

import (
	"IM_Server/grpc/proto_service"
	"IM_Server/model"
	"IM_Server/util"
	"errors"
	"fmt"
)

func UserEdit(userRpc *proto_service.User) error {

	tx := model.DB.Begin()
	userModel := model.User{
		Name:       userRpc.Name,
		Avatar:     userRpc.Avatar,
		ProjectId:  int(userRpc.ProjectId),
		ProjectUid: int(userRpc.ProjectUid),
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

func UserDel(userRpc *proto_service.User) error {

	tx := model.DB.Begin()
	userModel := model.User{
		ProjectId:  int(userRpc.ProjectId),
		ProjectUid: int(userRpc.ProjectUid),
	}
	dbUser := &model.User{}
	dbUserCount := 0
	tx.Where(model.User{ProjectId: userModel.ProjectId, ProjectUid: userModel.ProjectUid}).Find(dbUser).Count(&dbUserCount)
	if dbUserCount > 0 {
		//删除
		tx.Delete(dbUser)
	}
	tx.Commit()
	return nil

}

func QueryUserList(req *proto_service.QueryCommonReq) []*proto_service.IMUser {

	tx := model.DB.Begin()
	userList := []model.User{}
	if req.GroupId == 0 {
		tx.Where(model.User{ProjectId: int(req.ProjectId), ProjectUid: int(req.ProjectUid)}).Find(&userList)

	} else {
		//todo 尝试下left join
		groupUserList := []model.UserGroup{}
		tx.Where(model.UserGroup{GroupId: int(req.GroupId)}).Find(&groupUserList)
		userIds, _ := util.ArrayColumn(groupUserList, "Uid")
		tx.Where(map[string]interface{}{"project_id": req.ProjectId}).Where("id in (?)", userIds).Find(&userList)

	}
	IMUsers := []*proto_service.IMUser{}
	for _, v := range userList {
		IMUsers = append(IMUsers, &proto_service.IMUser{
			Id:             int64(v.ID),
			ProjectId:      int64(v.ProjectId),
			ProjectUid:     int64(v.ProjectUid),
			Name:           v.Name,
			IsOnline:       int64(v.IsOnline),
			Avatar:         v.Avatar,
			Pwd:            v.Pwd,
			LastOnlineTime: int64(v.LastOnlineTime),
		})
	}
	tx.Commit()
	return IMUsers

}
