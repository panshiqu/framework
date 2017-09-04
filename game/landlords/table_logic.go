package landlords

import (
	"github.com/panshiqu/framework/define"
)

// TableLogic 桌子逻辑
type TableLogic struct {
	tableFrame define.ITableFrame
}

// OnInit 初始化
func (t *TableLogic) OnInit() error {
	return nil
}

// OnGameStart 游戏开始
func (t *TableLogic) OnGameStart() error {
	return nil
}

// OnGameConclude 游戏结束
func (t *TableLogic) OnGameConclude() error {
	return nil
}

// OnUserSitDown 用户坐下
func (t *TableLogic) OnUserSitDown(userItem define.IUserItem) error {
	return nil
}

// OnUserStandUp 用户站起
func (t *TableLogic) OnUserStandUp(userItem define.IUserItem) error {
	return nil
}

// OnUserReconnect 用户重连
func (t *TableLogic) OnUserReconnect(userItem define.IUserItem) error {
	return nil
}

// OnMessage 收到消息
func (t *TableLogic) OnMessage(scmd uint16, data []byte, userItem define.IUserItem) error {
	return nil
}

// OnTimer 定时器
func (t *TableLogic) OnTimer(id int, parameter interface{}) error {
	return nil
}

// NewTableLogic 新建桌子逻辑
func NewTableLogic(v define.ITableFrame) define.ITableLogic {
	return &TableLogic{
		tableFrame: v,
	}
}
