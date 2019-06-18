package landlords

import (
	"log"
	"time"

	"../../define"
)

// TableLogic 桌子逻辑
type TableLogic struct {
	tableFrame define.ITableFrame
}

// OnInit 初始化
func (t *TableLogic) OnInit() error {
	log.Println("TableLogic OnInit")
	return nil
}

// OnGameStart 游戏开始
func (t *TableLogic) OnGameStart() error {
	log.Println("TableLogic OnGameStart")
	time.AfterFunc(time.Minute, func() {
		if err := t.OnGameConclude(); err != nil {
			log.Println("TableLogic OnGameConclude", err)
		}
	})
	return nil
}

// OnGameConclude 游戏结束
func (t *TableLogic) OnGameConclude() error {
	log.Println("TableLogic OnGameConclude")
	t.tableFrame.ConcludeGame()
	return nil
}

// OnUserSitDown 用户坐下
func (t *TableLogic) OnUserSitDown(userItem define.IUserItem) error {
	log.Println("TableLogic OnUserSitDown", userItem.UserID())
	return nil
}

// OnUserStandUp 用户站起
func (t *TableLogic) OnUserStandUp(userItem define.IUserItem) error {
	log.Println("TableLogic OnUserStandUp", userItem.UserID())
	return nil
}

// OnUserReconnect 用户重连
func (t *TableLogic) OnUserReconnect(userItem define.IUserItem) error {
	log.Println("TableLogic OnUserReconnect", userItem.UserID())
	return nil
}

// OnMessage 收到消息
func (t *TableLogic) OnMessage(scmd uint16, data []byte, userItem define.IUserItem) error {
	log.Println("TableLogic OnMessage", userItem.UserID(), scmd)
	return nil
}

// OnTimer 定时器
func (t *TableLogic) OnTimer(id int, parameter interface{}) error {
	return nil
}

// NewTableLogic 新建桌子逻辑
func NewTableLogic(v define.ITableFrame) define.ITableLogic {
	t := &TableLogic{
		tableFrame: v,
	}

	if err := t.OnInit(); err != nil {
		log.Println("TableLogic OnInit", err)
		return nil
	}

	return t
}
