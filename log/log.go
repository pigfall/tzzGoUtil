package log

import (
	"fmt"
	"io"
	golog "log"
	"os"
)

const callDepth = 3

const (
	DEBUG_PREFIX = "DEBUG"
	ERROR_PREFIX = "ERROR"
	INFO_PREFIX  = "INFO"
)

type Logger struct {
	l *golog.Logger
}

func NewLogger() *Logger {
	return &Logger{l: golog.New(os.Stdout, "", golog.Lshortfile|golog.LstdFlags)}
}

func (this *Logger) Infof(format string, msg ...interface{}) {
	this.printf(INFO_PREFIX, format, msg...)
}

func (this *Logger) Info(msg ...interface{}) {
	this.print(INFO_PREFIX, msg)
}

func (this *Logger) Debugf(format string, msg ...interface{}) {
	this.printf(INFO_PREFIX, format, msg...)
}

func (this *Logger) Debug(msg ...interface{}) {
	this.print(INFO_PREFIX, msg)
}

func (this *Logger) Error(msg ...interface{}) {
	this.print(ERROR_PREFIX, msg...)
}

func (this *Logger) Errorf(format string, msg ...interface{}) {
	this.printf(ERROR_PREFIX, format, msg...)

}
func (this *Logger) printf(prefix string, format string, msg ...interface{}) {
	this.l.Output(callDepth, fmt.Sprintf("%s: %s", prefix, fmt.Sprintf(format, msg...)))
}

func (this *Logger) print(prefix string, msg ...interface{}) {
	this.l.Output(callDepth, fmt.Sprintln(prefix, ": ", msg))
}

func (this *Logger) SetLongOutput() {
	this.l.SetFlags(golog.LstdFlags | golog.Llongfile)
}

func (this *Logger) SetOutput(writer io.Writer) {
	this.l.SetOutput(writer)
}

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
