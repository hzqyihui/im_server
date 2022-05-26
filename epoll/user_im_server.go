package epoll_server

import (
	"IM_Server/Util"
	"IM_Server/model"
	"errors"
	"fmt"
)

/**
用户上线：
1.修改online字段
*/
func UserIMOnline(request IMMessage) model.User {

	tx := model.DB.Begin()
	userInfo := model.User{}
	tx.Where(map[string]interface{}{"project_id": request.ProjectId, "project_uid": request.ProjectUid}).Find(&userInfo)

	//上线
	userInfo.IsOnline = 1

	if err := tx.Save(&userInfo).Error; err != nil {
		fmt.Println("上线更新失败")
	}

	return userInfo
}

func UserIMSaveMessage(request IMMessage, isSend bool) error {

	tx := model.DB.Begin()
	sendUid, _ := model.GetUidByProjectInfo(request.ProjectUid, request.ProjectId)
	sendToUid, _ := model.GetUidByProjectInfo(request.ProjectUid, request.ProjectId)

	messageModel := model.Message{
		Uid:         sendUid,
		GroupId:     request.GroupId,
		ToUid:       sendToUid,
		Message:     request.Data,
		MessageType: request.Type,
		Time:        request.Time,
		IsSend:      Util.BoolToInt(isSend),
	}

	if err := tx.Create(&messageModel).Error; err != nil {
		return errors.New("UserIMSaveMessageError")
	}
	tx.Commit()
	return nil
}
