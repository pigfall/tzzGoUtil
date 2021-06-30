package log

import (
	"io"
)

type LoggerI interface {
	Error(msg ...interface{})
	Errorf(format string, msg ...interface{})
	Debug(msg ...interface{})
	Debugf(format string, args ...interface{})
	Info(msg ...interface{})
	Infof(format string, msg ...interface{})
	SetLongOutput()
	SetOutput(writer io.Writer)
}

type LoggerLite interface{
	Error(msg ...interface{})
	Errorf(format string, msg ...interface{})
	Debug(msg ...interface{})
	Debugf(format string, args ...interface{})
	Info(msg ...interface{})
	Infof(format string, msg ...interface{})
}
