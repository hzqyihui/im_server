package service

import (
	"IM_Server/grpc/proto_service"
	"IM_Server/model"
	"errors"
	"fmt"
)

func GroupEdit(r *proto_service.Group) error {

	tx := model.DB.Begin()
	groupModel := model.Group{
		Name:      r.Name,
		Avatar:    *r.Avatar,
		CreateUid: int(r.Uid),
		Des:       r.Des,
		BaseModel: model.BaseModel{
			ID: uint(*r.GroupId),
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
