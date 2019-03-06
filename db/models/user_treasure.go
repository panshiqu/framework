package models

import (
	"strconv"
	"time"
)

type UserTreasure struct {
	Id          uint32    `xorm:"autoincr pk notnull"`
	UserId      uint32    `xorm:"notnull"`
	UserScore   int64     `xorm:"notnull"`
	UserDiamond int64     `xorm:"notnull"`
	CreatedAt   time.Time `xorm:"created notnull"`
	UpdatedAt   time.Time `xorm:"updated notnull"`
}

func NewUserTreasure() *UserTreasure {
	return &UserTreasure{}
}

// 初始化用户财富数据
func (ut *UserTreasure) InitUserTreasure(userId uint32) (*UserTreasure, error) {
	newUserTreasure := UserTreasure{
		UserId:      userId,
		UserScore:   0,
		UserDiamond: 0,
	}
	_, err := GameEngine.Insert(&newUserTreasure)
	if err != nil {
		DBLogError("db insert user treasure error", err)
		return nil, err
	}
	return &newUserTreasure, nil
}

// 通过用户ID获得用户财富信息
func (ut *UserTreasure) GetTreasureByUserId(userId uint32) (*UserTreasure, error) {
	treasure := new(UserTreasure)
	has, err := GameEngine.Where("user_id=?", userId).Get(treasure)
	if err != nil {
		DBLogError("GetTreasureByUserId error,UserId:"+strconv.Itoa(int(userId)), err)
		return nil, err
	}
	if has {
		return treasure,nil
	}else {
		return nil,nil
	}
}

// 更新分数钻石
func (ut *UserTreasure) UpdateUserTreasure(userId uint32, changeScore int64, changeDiamond int64) error {

	if _, err := GameEngine.Exec("UPDATE game_user_treasure SET user_score = user_score + ?, user_diamond = user_diamond + ? WHERE user_id = ?", changeScore, changeDiamond, userId); err != nil {
		DBLogError("UpdateUserTreasure error,UserId:"+strconv.Itoa(int(userId)), err)
		return err
	}
	return nil
}
