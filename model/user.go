package model

import (
	"golang.org/x/crypto/bcrypt"
)

// User 用户模型
type User struct {
	BaseModel
	Name       string
	ProjectId  int
	ProjectUid int
	Pwd        string
	Avatar     string `gorm:"size:1000"`
}

func (User) TableName() string {
	return "user"
}

const (
	// PassWordCost 密码加密难度
	PassWordCost = 12
	// Active 激活用户
	Active int = 1
	// Inactive 未激活用户
	Inactive int = 0
	// Suspend 被封禁用户
	Suspend int = 0
)

// GetUser 用ID获取用户
func GetUser(ID interface{}) (User, error) {
	var user User
	result := DB.First(&user, ID)
	return user, result.Error
}

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	user.Pwd = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(password))
	return err == nil
}
