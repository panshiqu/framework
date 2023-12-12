package define

import (
	"fmt"
)

const (
	// ErrnoSuccess 成功
	ErrnoSuccess int = 0

	// ErrnoFailure 失败
	ErrnoFailure int = 1
)

var (
	// ErrSuccess 成功
	ErrSuccess = &MyError{ErrField{Errno: ErrnoSuccess, Errdesc: "success"}}

	// ErrFailure 失败
	ErrFailure = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "failure"}}

	// ErrSignature 签名
	ErrSignature = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "signature"}}

	// ErrDisconnect 断开连接
	ErrDisconnect = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "disconnect"}}

	// ErrLengthLimit 长度限制
	ErrLengthLimit = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "length limit"}}

	// ErrMessageHead 非法消息头
	ErrMessageHead = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "message head"}}

	// ErrUnknownMainCmd 未知主命令
	ErrUnknownMainCmd = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "unknown main cmd"}}

	// ErrUnknownSubCmd 未知子命令
	ErrUnknownSubCmd = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "unknown sub cmd"}}

	// ErrRepeatRegisterService 重复注册服务
	ErrRepeatRegisterService = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "repeat register service"}}

	// ErrNotExistService 不存在该服务
	ErrNotExistService = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "not exist service"}}

	// ErrServiceAlreadyOpen 服务已经开启
	ErrServiceAlreadyOpen = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "service already open"}}

	// ErrServiceAlreadyShut 服务已经关闭
	ErrServiceAlreadyShut = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "service already shut"}}

	// ErrNotExistUser 不存在该用户
	ErrNotExistUser = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "not exist user"}}

	// ErrUserNotSit 用户没有坐下
	ErrUserNotSit = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "user not sit"}}

	// ErrTableStatus 桌子状态
	ErrTableStatus = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "table status"}}

	// ErrNotEnoughScore 分数不足
	ErrNotEnoughScore = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "not enough score"}}

	// ErrNotEnoughDiamond 钻石不足
	ErrNotEnoughDiamond = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "not enough diamond"}}

	// ErrNotYourTurn 没轮到你
	ErrNotYourTurn = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "not your turn"}}

	// ErrAlreadyPlaceStone 已经落子
	ErrAlreadyPlaceStone = &MyError{ErrField{Errno: ErrnoFailure, Errdesc: "already place stone"}}

	// ErrInAnotherGame 在其它游戏中
	ErrInAnotherGame = ErrField{Errno: ErrnoFailure, Errdesc: "in another game"}
)

// MyError 错误
type MyError struct {
	ErrField
}

// ErrField 错误字段
// 所有回复结构想嵌套MyError，但db需断言error
type ErrField struct {
	Errno   int    `json:",omitempty"` // 错误码
	Errdesc string `json:",omitempty"` // 错误描述
}

func (m *MyError) Error() string {
	return fmt.Sprintf(`{"Errno":%d,"Errdesc":"%s"}`, m.Errno, m.Errdesc)
}
