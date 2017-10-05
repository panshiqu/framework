package fiveinarow

import (
	"encoding/json"
	"log"
	"math/rand"

	"github.com/panshiqu/framework/define"
)

// TableLogic 桌子逻辑
type TableLogic struct {
	tableFrame define.ITableFrame

	currentChair int     // 当前椅子
	checkerBoard [][]int // 五子棋盘
}

// OnInit 初始化
func (t *TableLogic) OnInit() error {
	log.Println("TableLogic OnInit")

	// 默认椅子
	t.currentChair = define.InvalidChair

	// 初始化五子棋盘
	t.checkerBoard = make([][]int, LineNumber)
	for i := 0; i < LineNumber; i++ {
		t.checkerBoard[i] = make([]int, LineNumber)
	}

	return nil
}

// OnGameStart 游戏开始
func (t *TableLogic) OnGameStart() error {
	log.Println("TableLogic OnGameStart")

	// 随机玩家
	t.currentChair = rand.Intn(define.CG.UserPerTable)

	// 广播开始
	t.tableFrame.SendTableJSONMessage(define.GameTable, GameBroadcastStart, &BroadcastStart{
		ChairID: t.currentChair,
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

	// 通知场景
	userItem.SendJSONMessage(define.GameTable, GameNotifyScene, &NotifyScene{
		Timeout:    Timeout,
		LineNumber: LineNumber,
	})

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

	switch scmd {
	case GamePlaceStone:
		placeStone := &PlaceStone{}

		if err := json.Unmarshal(data, placeStone); err != nil {
			return err
		}

		// 没轮到你
		if userItem.ChairID() != t.currentChair {
			return define.ErrNotYourTurn
		}

		// 已经落子
		if t.checkerBoard[placeStone.PositionX][placeStone.PositionY] != 0 {
			return define.ErrAlreadyPlaceStone
		}

		// 标记落子
		t.checkerBoard[placeStone.PositionX][placeStone.PositionY] = t.currentChair + 1

		// 广播落子
		t.tableFrame.SendTableJSONMessage(define.GameTable, GameBroadcastPlaceStone, &BroadcastPlaceStone{
			ChairID:   t.currentChair,
			PositionX: placeStone.PositionX,
			PositionY: placeStone.PositionY,
		})

		// 轮转玩家
		t.currentChair = (t.currentChair + 1) % define.CG.UserPerTable

	default:
		return define.ErrUnknownSubCmd
	}

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
