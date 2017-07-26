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

// StandUp 站起
func (t *TableFrame) StandUp(userItem *UserItem) {

}

// OnTimer 定时器
func (t *TableFrame) OnTimer(id int, parameter interface{}) {

}

// OnMessage 收到消息
func (t *TableFrame) OnMessage(mcmd uint16, scmd uint16, data []byte) {

}

// SendTableMessage 发送桌子消息
func (t *TableFrame) SendTableMessage(mcmd uint16, scmd uint16, data []byte) {

}

// SendChairMessage 发送椅子消息
func (t *TableFrame) SendChairMessage(mcmd uint16, scmd uint16, data []byte) {

}
