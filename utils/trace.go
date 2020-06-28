package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// TraceSwitch 开关
var TraceSwitch = true

// Trace trace
func Trace(name string, param ...interface{}) func() {
	if !TraceSwitch {
		return func() {}
	}

	start := time.Now()
	log.Println("####Enter", name, param)

	return func() {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		}

		file = strings.TrimSuffix(filepath.Base(file), ".go")

		log.Println("#####Exit", name, fmt.Sprintf("%s:%d", file, line), time.Since(start))
	}
}
