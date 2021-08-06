package log

import(
	"fmt"
	kit_log "github.com/go-kit/kit/log"
)



type LogLevel uint

const LevelKey ="level"

const(
	LevelDebug LogLevel = iota
	LevelInfo 
	LevelWarn
	LevelError
)


// Helper is a logger helper.
type Helper struct {
	debug Logger_Log
	info  Logger_Log
	warn  Logger_Log
	err   Logger_Log
}


type emptyLog struct{}

func(this *emptyLog) Log(kv ...interface{})error{
	return nil
}

// NewHelper new a logger helper.
func NewHelper(name string, rawLogger Logger_Log,level LogLevel) *Helper {
	logger := kit_log.With(rawLogger, "module", name)
	debugLog := DebugLogger(logger)
	if LevelDebug < level {
		debugLog = &emptyLog{}
	}
	infoLog := InfoLogger(logger)
	if LevelInfo < level {
		infoLog= &emptyLog{}
	}
	warnLog :=WarnLogger(logger)
	if LevelWarn < level{
		warnLog= &emptyLog{}
	}
	errLog :=ErrorLogger(logger)
	if LevelError < level{
		errLog= &emptyLog{}
	}
	return &Helper{
		debug: debugLog,
		info:  infoLog,
		warn:  warnLog,
		err:   errLog,
	}
}

// Debug logs a message at debug level.
func (h *Helper) Debug(a ...interface{}) {
	h.debug.Log("msg", fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func (h *Helper) Debugf(format string, a ...interface{}) {
	h.debug.Log("msg", fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func (h *Helper) Debugw(kv ...interface{}) {
	h.debug.Log(kv...)
}

// Info logs a message at info level.
func (h *Helper) Info(a ...interface{}) {
	h.info.Log("msg", fmt.Sprint(a...))
}

// Infof logs a message at info level.
func (h *Helper) Infof(format string, a ...interface{}) {
	h.info.Log("msg", fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func (h *Helper) Infow(kv ...interface{}) {
	h.info.Log(kv...)
}

// Warn logs a message at warn level.
func (h *Helper) Warn(a ...interface{}) {
	h.warn.Log("msg", fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func (h *Helper) Warnf(format string, a ...interface{}) {
	h.warn.Log("msg", fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func (h *Helper) Warnw(kv ...interface{}) {
	h.warn.Log(kv...)
}

// Error logs a message at error level.
func (h *Helper) Error(a ...interface{}) {
	h.err.Log("msg", fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func (h *Helper) Errorf(format string, a ...interface{}) {
	h.err.Log("msg", fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func (h *Helper) Errorw(kv ...interface{}) {
	h.err.Log(kv...)
}

func DebugLogger(log Logger_Log) Logger_Log {
	return kit_log.With(log,LevelKey,"DEBUG")
}

func InfoLogger(log Logger_Log) Logger_Log {
	return kit_log.With(log,LevelKey,"INFO")
}

func WarnLogger(log Logger_Log) Logger_Log {
	return kit_log.With(log, LevelKey,"WARN")
}

func ErrorLogger(log Logger_Log) Logger_Log {
	return kit_log.With(log, LevelKey,"ERROR")
}

