package utils

import (
	"log"
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
		log.Println("#####Exit", name, time.Since(start))
	}
}
