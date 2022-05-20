package model

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
