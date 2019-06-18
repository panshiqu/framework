package models

import (
	"strconv"
	"time"
)

type UserInformation struct {
	Id          uint32    `xorm:"autoincr pk notnull"`
	UserId      uint32    `xorm:"notnull unique"`
	UserAccount string    `xorm:"varchar(30) notnull unique"`
	UserName    string    `xorm:"notnull unique"`
	UserIcon    int       `xorm:"notnull"`
	UserGender  uint8     `xorm:"notnull"` //性别
	UserLevel   int       `xorm:"notnull"`
	BindPhone   string    `xorm:"char(11) notnull"`
	UserScore   int64    `xorm:"notnull"`
	UserDiamond int64    `xorm:"notnull"`
	IsRobot     uint8     `xorm:"notnull"`
	CreatedAt   time.Time `xorm:"created notnull"`
	UpdatedAt   time.Time `xorm:"updated notnull"`
}

func NewUserInfo() *UserInformation {
	return &UserInformation{}
}

// 添加用户信息
func (userInfo *UserInformation) AddUserInfoByUser(user *User)(*UserInformation,error) {
	info := UserInformation{
		UserId:      user.Id,
		UserDiamond: 0,
		UserScore:   0,
		UserLevel:   1,
		UserName:    user.Name,
		UserAccount: user.Account,
	}

	_, err := GameEngine.Insert(&info)
	if err != nil {
		DBLogError("db insert user info error",err)
		return nil,err
	}
	return &info, nil
}

// 通过账号获得用户信息
func (userInfo *UserInformation) GetInfoByAccount(account string) (*UserInformation, error) {
	retUserInfo := new(UserInformation)
	has, err := GameEngine.Where("user_account = ?", account).Get(retUserInfo)
	if err != nil {
		DBLogError("GetInfoByAccount error,account:"+account,err)
		return nil, err
	}
	if has {
		return retUserInfo, nil
	}else {
		return nil, nil
	}

}

func (userInfo *UserInformation) GetInfoByName(name string) (*UserInformation, error) {
	retUserInfo := new(UserInformation)
	has, err := GameEngine.Where("user_name = ?", name).Get(retUserInfo)
	if err != nil {
		DBLogError("GetInfoByName error,name:"+name,err)
		return nil, err
	}
	if has {
		return retUserInfo, nil
	}else {
		return nil, nil
	}
}

//  根据用户ID获得用户信息
func (userInfo *UserInformation)GetInfoByUserId(userId uint32)(*UserInformation,error)  {
	info := new(UserInformation)
	has, err := GameEngine.Where("user_id=?", userId).Get(info)
	if err != nil {
		DBLogError("GetInfoByName error,UserId:"+strconv.Itoa(int(userId)),err)
		return nil,err
	}
	if has {
		return info, nil
	}else {
		return nil, nil
	}
}


