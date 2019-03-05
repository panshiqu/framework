package db

import "time"

//用户信息
type User struct {
	Id uint32 `xorm:"autoincr pk notnull"`
	Account string `xorm:"varchar(30) notnull unique"`		// 账户
	Password string `xorm:"char(32) notnull"`	// 密码（未使用）
	Name string `xorm:"varchar(30) notnull"`	// 名称
	Icon int `xorm:"notnull"`	// 图标
	Gender uint8  `xorm:"notnull"` // 性别
	Ip	string `xorm:"varchar(15) notnull"` //Ip地址
	CreatedAt time.Time `xorm:"created notnull"`
	UpdatedAt time.Time `xorm:"updated notnull"`
}

type UserInformation struct {
	Id uint32
	UserId uint32 `xorm:"notnull unique"`
	UserName string `xorm:"notnull unique"`
	UserIcon int `xorm:"notnull"`
	UserGender uint8 `xorm:"notnull"` //性别
	UserLevel int `xorm:"notnull"`
	BindPhone string `xorm:"char(11) notnull"`
	UserScore uint32 `xorm:"notnull"`
	UserDiamond uint32 `xorm:"notnull"`
	IsRobot uint8 `xorm:"notnull"`
	CreatedAt time.Time `xorm:"created notnull"`
	UpdatedAt time.Time `xorm:"updated notnull"`
}


type UserTreasure struct {
	Id uint32 `xorm:"autoincr pk notnull"`
	UserId uint32 `xorm:"notnull"`
	UserScore uint32 `xorm:"notnull"`
	UserDiamond uint32 `xorm:"notnull"`
	CreatedAt time.Time `xorm:"created notnull"`
	UpdatedAt time.Time `xorm:"updated notnull"`
}
