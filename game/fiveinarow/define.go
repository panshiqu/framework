package fiveinarow

const (
	// Timeout 超时
	Timeout = 15

	// LineNumber 线数
	LineNumber = 15
)

const (
	// GamePlaceStone 落子
	GamePlaceStone = 1

	// GameBroadcastPlaceStone 广播落子
	GameBroadcastPlaceStone = 101

	// GameNotifyScene 通知场景
	GameNotifyScene = 201

	// GameBroadcastStart 广播开始
	GameBroadcastStart = 202

	// GameBroadcastConclude 广播结束
	GameBroadcastConclude = 203
)

// NotifyScene 通知场景
type NotifyScene struct {
	Timeout    int `json:",omitempty"` // 超时
	LineNumber int `json:",omitempty"` // 线数
}

// BroadcastStart 广播开始
type BroadcastStart struct {
	ChairID int `json:",omitempty"` // 玩家
}

// PlaceStone 落子
type PlaceStone struct {
	PositionX int `json:",omitempty"` // 位置
	PositionY int `json:",omitempty"` // 位置
}

// BroadcastPlaceStone 广播落子
type BroadcastPlaceStone struct {
	ChairID   int `json:",omitempty"` // 玩家
	PositionX int `json:",omitempty"` // 位置
	PositionY int `json:",omitempty"` // 位置
}

// BroadcastConclude 广播结束
type BroadcastConclude struct {
	ChairID int `json:",omitempty"` // 赢家
}
