package models

import (
	"../../define"
	"time"
)

//用户信息
type User struct {
	Id        uint32    `xorm:"autoincr pk notnull"`
	Account   string    `xorm:"varchar(30) notnull unique"` // 账户
	Password  string    `xorm:"char(32) notnull"`           // 密码（未使用）
	Name      string    `xorm:"varchar(30) notnull"`        // 名称
	Icon      int       `xorm:"notnull"`                    // 图标
	Gender    uint8     `xorm:"notnull"`                    // 性别
	Ip        string    `xorm:"varchar(15) notnull"`        //Ip地址
	CreatedAt time.Time `xorm:"created notnull"`
	UpdatedAt time.Time `xorm:"updated notnull"`
}

func NewUser() *User  {
	return &User{}
}

// 添加新用户
func (u *User) RegisterUser(info *define.FastRegister)(*User,error) {
	newUser := User{
		Account:  info.Account,
		Password: info.Password,
		Name:     info.Name,
		Icon:     info.Icon,
		Ip:       info.IP,
		Gender:   uint8(info.Gender),
	}
	_, err := GameEngine.Insert(&newUser)

	if err != nil {
		DBLogError("db insert user error",err)
		return nil,err
	}
	return &newUser, nil
}