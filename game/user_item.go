package game

import (
	"net"

	"github.com/panshiqu/framework/define"
)

// UserItem 用户
type UserItem struct {
	id      int      // 编号
	name    string   // 名称
	icon    int      // 图标
	level   int      // 等级
	gender  int      // 性别
	phone   string   // 手机
	score   int64    // 分数
	diamond int64    // 钻石
	robot   bool     // 机器人
	conn    net.Conn // 网络连接

	status     int         // 状态
	chairID    int         // 椅子编号
	tableFrame *TableFrame // 桌子框架
}

// UserID 用户编号
func (u *UserItem) UserID() int {
	return u.id
}

// UserName 用户名称
func (u *UserItem) UserName() string {
	return u.name
}

// UserIcon 用户图标
func (u *UserItem) UserIcon() int {
	return u.icon
}

// UserLevel 用户等级
func (u *UserItem) UserLevel() int {
	return u.level
}

// UserGender 用户性别
func (u *UserItem) UserGender() int {
	return u.gender
}

// BindPhone 用户手机
func (u *UserItem) BindPhone() string {
	return u.phone
}

// UserScore 用户分数
func (u *UserItem) UserScore() int64 {
	return u.score
}

// UserDiamond 用户钻石
func (u *UserItem) UserDiamond() int64 {
	return u.diamond
}

// IsRobot 是否机器人
func (u *UserItem) IsRobot() bool {
	return u.robot
}

// UserStatus 用户状态
func (u *UserItem) UserStatus() int {
	return u.status
}

// SetUserStatus 设置用户状态
func (u *UserItem) SetUserStatus(v int) {
	u.status = v
}

// ChairID 椅子编号
func (u *UserItem) ChairID() int {
	return u.chairID
}

// SetChairID 设置椅子编号
func (u *UserItem) SetChairID(v int) {
	u.chairID = v
}

// TableID 桌子编号
func (u *UserItem) TableID() int {
	if u.tableFrame != nil {
		return u.tableFrame.TableID()
	}

	return define.InvalidTable
}

// TableFrame 桌子框架
func (u *UserItem) TableFrame() *TableFrame {
	return u.tableFrame
}

// SetTableFrame 设置桌子框架
func (u *UserItem) SetTableFrame(v *TableFrame) {
	u.tableFrame = v
}
