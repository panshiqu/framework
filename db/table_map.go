package db

//用户信息
type User struct {
	Id uint32 `xorm:"autoincr pk notnull"`
	Account string `xorm:"notnull"`		// 账户
	Password string `xorm:"notnull"`	// 密码（未使用）
	Name string `xorm:"notnull"`	// 名称
	Icon int16 `xorm:"notnull"`	// 图标
	Gender uint8  `xorm:"notnull"` // 性别
	Ip	string `xorm:"notnull"` //Ip地址
}


type UserTreasure struct {
	Id uint32 `xorm:"autoincr pk notnull"`
	UserId uint32 `xorm:"notnull"`
	UserScore uint32 `xorm:"notnull"`
	UserDiamond uint32 `xorm:"notnull"`
}
