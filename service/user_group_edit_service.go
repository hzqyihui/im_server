package service

import (
	"IM_Server/Util"
	"IM_Server/grpc/proto_service"
	"IM_Server/model"
	"errors"
)

func UserGroupEdit(r *proto_service.UserGroup) error {

	tx := model.DB.Begin()
	//ProjectId int64             `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	//Uid       int64             `protobuf:"varint,2,opt,name=uid,proto3" json:"uid,omitempty"`
	//Uids      []int64           `protobuf:"varint,3,rep,packed,name=uids,proto3" json:"uids,omitempty"`
	//GroupId   int64             `protobuf:"varint,4,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	//Optional  GroupOptionalType `protobuf:"varint,5,opt,name=optional,proto3,enum=proto_service.GroupOptionalType" json:"optional,omitempty"` //del / add

	if r.GroupId <= 0 {
		return errors.New("用户组id不合法")
	}
	if len(r.Uids) <= 0 {
		return errors.New("用户ids不合法")
	}
	relationList := []model.UserGroup{}
	tx.Where(map[string]interface{}{"group_id": r.GroupId}).Find(&relationList)
	optionalUids := r.Uids
	dbUids, _ := Util.ArrayColumn(relationList, "Uid")

	//先简化下 ，必须要有group_id
	if r.Optional == proto_service.GroupOptionalType_DEL {
		//把这些人剔除
		deleteIds := []int{}

		for _, v := range relationList {
			if in, _ := Util.InArray(v.Uid, optionalUids); in {
				deleteIds = append(deleteIds, int(v.ID))
			}
		}

		if err := tx.Where(deleteIds).Delete(model.UserGroup{}).Error; err != nil {
			return errors.New("UserGroupEdit删除错误：" + err.Error())
		}
	} else {
		//把这些人新增
		addUIds := []int{}
		for _, v := range optionalUids {
			if in, _ := Util.InArray(v, dbUids); in == false {
				addUIds = append(addUIds, int(v))
			}
		}
		addUserGroup := []model.UserGroup{}
		for _, v := range addUIds {
			addUserGroup = append(addUserGroup, model.UserGroup{
				Uid:     int(v),
				GroupId: int(r.GroupId),
			})
			insertModel := model.UserGroup{
				Uid:     int(v),
				GroupId: int(r.GroupId),
			}
			if err := tx.Create(&insertModel).Error; err != nil {
				return errors.New("UserGroupEdit")
			}
		}
	}
	tx.Commit()
	return nil

}
