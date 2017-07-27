package game

import "github.com/panshiqu/framework/define"

// TableFrame 桌子框架
type TableFrame struct {
	id     int         // 编号
	status int         // 状态
	users  []*UserItem // 用户
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
func (t *TableFrame) UserCount() (cnt int) {
	for _, v := range t.users {
		if v != nil {
			cnt++
		}
	}

	return
}

// NilChairID 空椅子编号
func (t *TableFrame) NilChairID() int {
	for k, v := range t.users {
		if v == nil {
			return k
		}
	}

	return define.InvalidChair
}

// SitDown 坐下
func (t *TableFrame) SitDown(userItem *UserItem) {
	chair := t.NilChairID()
	t.users[chair] = userItem
	userItem.SetChairID(chair)
	userItem.SetTableFrame(t)
}

// StandUp 站起
func (t *TableFrame) StandUp(userItem *UserItem) {
	t.users[userItem.ChairID()] = nil
	userItem.SetChairID(define.InvalidChair)
	userItem.SetTableFrame(nil)
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
