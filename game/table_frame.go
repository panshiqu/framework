package game

// TableFrame 桌子框架
type TableFrame struct {
	id     int
	status int
	users  []*UserItem
}

// TableID 桌子编号
func (t *TableFrame) TableID() int {
	return t.id
}

// TableStatus 桌子状态
func (t *TableFrame) TableStatus() int {
	return t.status
}

// UserCount 用户数量
func (t *TableFrame) UserCount() int {
	return 0
}

// SitDown 坐下
func (t *TableFrame) SitDown(userItem *UserItem) {

}
