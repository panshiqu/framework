package define

import (
	"encoding/json"
	"fmt"
)

const (
	// ErrSuccess 成功
	ErrSuccess int = 0
)

// MyError 错误
type MyError struct {
	Errno   int    `json:",omitempty"` // 错误码
	Errdesc string `json:",omitempty"` // 错误描述
}

func (m *MyError) Error() string {
	return fmt.Sprintf(`{"Errno":%d,"Errdesc":"%s"}`, m.Errno, m.Errdesc)
}

// NewError 创建错误
func NewError(desc string) error {
	return fmt.Errorf(`{"Errno":1,"Errdesc":"%s"}`, desc)
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
