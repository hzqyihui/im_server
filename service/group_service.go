package service

import (
	"IM_Server/Util"
	"IM_Server/grpc/proto_service"
	"IM_Server/model"
	"errors"
	"fmt"
)

func GroupEdit(r *proto_service.Group) error {

	tx := model.DB.Begin()
	groupModel := model.Group{
		Name:   r.Name,
		Avatar: r.Avatar,
		//CreateUid: int(r.Uid),todo
		Des: r.Des,
		BaseModel: model.BaseModel{
			ID: uint(r.GroupId),
		},
	}

	fmt.Println(groupModel)

	if groupModel.ID > 0 {
		//修改
		fmt.Println("编辑", groupModel)
		dbCount := 0
		dbItem := &model.Group{}
		tx.Where(map[string]interface{}{"id": groupModel.ID}).Find(dbItem).Count(&dbCount)
		dbItem.Name = groupModel.Name
		dbItem.Avatar = groupModel.Avatar
		if err := tx.Save(&dbItem).Error; err != nil {
			return errors.New("GroupEdit更新失败")
		}
	} else {
		//新增
		fmt.Println("新增", groupModel)
		if err := tx.Create(&groupModel).Error; err != nil {
			return errors.New("GroupEdit新增失败")
		}
	}

	tx.Commit()
	return nil

}

func GroupDel(r *proto_service.Group) error {

	tx := model.DB.Begin()
	if err := tx.Where(map[string]interface{}{"id": r.GroupId}).Delete(model.Group{}); err != nil {
		return errors.New("GroupDel失败")
	}

	tx.Commit()
	return nil

}

//查询这个人有那些组
func QueryGroupList(req *proto_service.QueryCommonReq) []*proto_service.IMGroup {

	tx := model.DB.Begin()
	uid, _ := model.GetUidByProjectInfo(int(req.ProjectUid), int(req.ProjectId))
	groupUserList := []model.UserGroup{}
	tx.Where(model.UserGroup{Uid: uid}).Find(&groupUserList)
	groupIds, _ := util.ArrayColumn(groupUserList, "Group")

	groupList := []model.Group{}
	tx.Where("id in (?)", groupIds).Find(&groupList)

	IMGroups := []*proto_service.IMGroup{}
	for _, v := range groupList {
		IMGroups = append(IMGroups, &proto_service.IMGroup{
			Id:     int64(v.ID),
			Name:   v.Name,
			Avatar: v.Avatar,
			Des:    v.Des,
		})
	}
	tx.Commit()
	return IMGroups

}
