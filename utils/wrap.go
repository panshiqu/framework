package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// Wrap .
func Wrap(err error, info ...interface{}) error {
	if err == nil {
		return nil
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "Unknown"
		line = -1
	}

	file = strings.TrimSuffix(filepath.Base(file), ".go")

	if len(info) == 0 {
		return fmt.Errorf("%s:%d > %w", file, line, err)
	}

	return fmt.Errorf("%s:%d%v > %w", file, line, info, err)
}
