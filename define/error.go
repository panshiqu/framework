package define

import (
	"encoding/json"
	"fmt"
)

const (
	// ErrSuccess 成功
	ErrSuccess int = 0
)

// Error 错误
type Error struct {
	Errno   int    `json:",omitempty"` // 错误码
	Errdesc string `json:",omitempty"` // 错误描述
}

func (e *Error) Error() string {
	return fmt.Sprintf(`{"Errno":%d,"Errdesc":"%s"}`, e.Errno, e.Errdesc)
}

// NewError 创建错误
func NewError(desc string) error {
	return fmt.Errorf(`{"Errno":1,"Errdesc":"%s"}`, desc)
}

// CheckError 检查错误
func CheckError(data []byte) error {
	eno := &Error{}

	if err := json.Unmarshal(data, eno); err != nil {
		return err
	}

	if eno.Errno != ErrSuccess {
		return eno
	}

	return nil
}
