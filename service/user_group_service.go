package service

import (
	"IM_Server/grpc/proto_service"
	"IM_Server/model"
	"IM_Server/util"
	"errors"
	"fmt"
)

func UserGroupEdit(r *proto_service.UserGroup) error {

	tx := model.DB.Begin()

	//todo 有错误的时候需不需要tx.Rollback() ？

	//ProjectId    int64   `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	//ProjectUid   int64   `protobuf:"varint,2,opt,name=project_uid,json=projectUid,proto3" json:"project_uid,omitempty"`
	//ProjectUsers []*User `protobuf:"bytes,3,rep,name=project_users,json=projectUsers,proto3" json:"project_users,omitempty"`
	//GroupId      int64   `protobuf:"varint,4,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`

	if r.GroupId <= 0 {
		return errors.New("用户组id不合法")
	}
	if len(r.ProjectUsers) <= 0 {
		return errors.New("用户ids不合法")
	}
	//todo 判断group_id是否存在

	//删除
	tx.Where(map[string]interface{}{"group_id": r.GroupId}).Delete(model.UserGroup{})
	//新增
	dbUids, _ := util.ArrayColumn(r.ProjectUsers, "ProjectUid")

	addIMUsers := []model.User{}
	tx.Where(map[string]interface{}{"project_id": r.ProjectId}).Where("project_uid in (?)", dbUids).Find(&addIMUsers)

	//todo 这里查询不出来
	//把这些人新增
	addUserGroup := []model.UserGroup{}
	for _, v := range addIMUsers {
		addUserGroup = append(addUserGroup, model.UserGroup{
			Uid:     int(v.ID),
			GroupId: int(r.GroupId),
		})
	}
	fmt.Println(111, addUserGroup)
	if len(addUserGroup) > 0 {
		if err := tx.Create(&addUserGroup).Error; err != nil {
			return errors.New("UserGroupEdit")
		}
	}

	tx.Commit()
	return nil

}
