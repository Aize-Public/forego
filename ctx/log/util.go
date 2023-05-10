package log

import (
	"fmt"
	"runtime"
)

func caller(above int) string {
	_, file, line, _ := runtime.Caller(above + 1)
	return fmt.Sprintf("%s:%d", file, line)
}
