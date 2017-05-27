package define

import (
	"fmt"
)

// NewError 创建错误
func NewError(desc string) error {
	return fmt.Errorf(`{"errno":1,"errdesc":"%s"}`, desc)
}
