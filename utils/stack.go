package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

// SafeCall .
func SafeCall(fn func(...any), args ...any) {
	defer func() {
		if err := recover(); err != nil {
			errStack(debug.Stack(), err)
		}
	}()
	fn(args...)
}

func errStack(stack []byte, err ...any) {
	data := stack
	if len(err) > 0 {
		data = []byte(fmt.Sprintf("err: %v\n%s", err[0], string(stack)))
	}
	if f, e := os.OpenFile(time.Now().Format("stack20060102150405.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); e != nil {
		log.Printf("OpenFile %v\n%s", e, string(data))
	} else {
		f.Write(data)
		f.Close()
	}
}

func Stack(err any, e ...*error) {
	if err != nil {
		if len(e) > 0 {
			*(e[0]) = errors.New("panic")
		}
		errStack(debug.Stack(), err)
	}
}

func StackAll() {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			errStack(buf[:n])
			return
		}
		buf = make([]byte, 2*len(buf))
	}
}
