package errorx

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func caller(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	file = filepath.Base(file)
	return fmt.Sprintf("%s:%d", file, line)
}
