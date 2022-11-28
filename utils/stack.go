package utils

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

// SafeCall .
func SafeCall(fn func(...any), args ...any) {
	defer func() {
		if err := recover(); err != nil {
			f, _ := os.OpenFile(time.Now().Format("stack2006-01-02T15:04:05.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			fmt.Fprintln(f, "err:", err)
			f.Write(debug.Stack())
			f.Close()
		}
	}()
	fn(args...)
}
