package bean

import (
	"time"
)

//用户信息
type User struct {
	Uuid       string
	TimeCreate time.Time //创建时间
	TimeUpdate time.Time //更新时间
	Nickname   string    //昵称
	Salt       string    `json:"omitempty"` //密码加盐
	Password   string    `json:"omitempty"` //密码
	Email      string
	Phone      string //手机号
	Role       string //角色 ENUM('admin','sales')

	Pid string `xorm:"-"`
}

func (this *User) TableName() string {
	return "user"
}
