package epoll_server

import (
	"IM_Server/model"
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
