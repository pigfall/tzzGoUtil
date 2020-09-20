package log

import (
	"fmt"
	golog "log"
	"os"
)

const callDepth = 3

const (
	DEBUG_PREFIX = "DEBUG"
	ERROR_PREFIX = "ERROR"
)

var i *golog.Logger

func init() {
	i = golog.New(os.Stdout, "tzzLog", golog.Lshortfile|golog.LstdFlags)
}

func print(prefix string, msg ...interface{}) {
	i.Output(callDepth, fmt.Sprintln(prefix, ":", msg))
}

func printf(prefix string, format string, args ...interface{}) {
	i.Output(callDepth, fmt.Sprintf("%s:%s", prefix, (fmt.Sprintf(format, args...))))
}

func Debug(msg ...interface{}) {
	print(DEBUG_PREFIX, msg)
}

func Debugf(format string, args ...interface{}) {
	printf(DEBUG_PREFIX, format, args...)
}

func Error(msg ...interface{}) {
	print(ERROR_PREFIX, msg)
}

func Errorf(format string, args ...interface{}) {
	printf(ERROR_PREFIX, format, args...)
}
