package debug

import (
	"runtime"
)

func CallerName() (name string) {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	return frame.Function
}
