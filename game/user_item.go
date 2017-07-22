package game

// UserItem 用户
type UserItem struct {
	id      int    // 编号
	name    string // 名称
	icon    int    // 图标
	level   int    // 等级
	gender  int    // 性别
	score   int64  // 分数
	diamond int64  // 钻石
	robot   bool   // 机器人
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
