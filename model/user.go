package model

import (
	"IM_Server/cache"
	"strconv"
)

// User 用户模型
type User struct {
	BaseModel
	Name       string
	ProjectId  int
	ProjectUid int
	Pwd        string
	Avatar     string `gorm:"size:1000"`
	IsOnline   int
}

func (User) TableName() string {
	return "user"
}

func GetUidByProjectInfo(projectUid int, projectId int) (int, error) {
	hashKey := "user_project_relation"
	relationField := strconv.Itoa(projectId) + "_" + strconv.Itoa(projectUid)

	if res, err := cache.RedisClient.HGet(hashKey, relationField).Result(); err != nil {
		return strconv.Atoi(res)
	} else {
		//查询
		userInfo := User{}
		DB.Where(map[string]interface{}{"project_id": projectId, "project_uid": projectUid}).Find(&userInfo)
		if userInfo.ID > 0 {
			cache.RedisClient.HSet(hashKey, relationField, userInfo.ID)
			//cache.RedisClient.Expire(hash_key,60*time.Second) 设置过期时间
			return int(userInfo.ID), nil
		}
	}
	return 0, nil
}
