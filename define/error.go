package define

import (
	"encoding/json"
	"fmt"
)

const (
	// ErrSuccess 成功
	ErrSuccess int = 0

	// ErrFailure 失败
	ErrFailure int = 1
)

var (
	// ErrDisconnect 断开连接
	ErrDisconnect = &MyError{Errno: ErrFailure, Errdesc: "disconnect"}

	// ErrUnknownMainCmd 未知主命令
	ErrUnknownMainCmd = &MyError{Errno: ErrFailure, Errdesc: "unknown main cmd"}

	// ErrUnknownSubCmd 未知子命令
	ErrUnknownSubCmd = &MyError{Errno: ErrFailure, Errdesc: "unknown sub cmd"}

	// ErrRepeatRegisterService 重复注册服务
	ErrRepeatRegisterService = &MyError{Errno: ErrFailure, Errdesc: "repeat register service"}

	// ErrNotExistService 不存在该服务
	ErrNotExistService = &MyError{Errno: ErrFailure, Errdesc: "not exist service"}

	// ErrServiceAlreadyOpen 服务已经开启
	ErrServiceAlreadyOpen = &MyError{Errno: ErrFailure, Errdesc: "service already open"}

	// ErrServiceAlreadyShut 服务已经关闭
	ErrServiceAlreadyShut = &MyError{Errno: ErrFailure, Errdesc: "service already shut"}
)

// MyError 错误
type MyError struct {
	Errno   int    `json:",omitempty"` // 错误码
	Errdesc string `json:",omitempty"` // 错误描述
}

func (m *MyError) Error() string {
	return fmt.Sprintf(`{"Errno":%d,"Errdesc":"%s"}`, m.Errno, m.Errdesc)
}

// CheckError 检查错误
func CheckError(data []byte) error {
	me := &MyError{}

	if err := json.Unmarshal(data, me); err != nil {
		return err
	}

	if me.Errno != ErrSuccess {
		return me
	}

	return nil
}
