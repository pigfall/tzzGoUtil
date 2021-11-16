package log

import(
		"log"
		"fmt"
		"os"
)

var globalLogger = &GoLogger{calldepth:go_logger_calldepth+1,Logger:log.New(os.Stdout,"",log.LstdFlags | log.Llongfile)}

type GoLogger struct{
	*log.Logger
	calldepth  int
}

func Debug(output string){
	globalLogger.Debug(output)
}

func Info(output string){
	globalLogger.Info(output)
}

func Error(output string){
	globalLogger.Error(output)
}

const (
		go_logger_calldepth = 3
)

func (this *GoLogger) Debug(output string){
	this.Logger.Output(this.calldepth,fmt.Sprintf(" DEBUG: %s",output))
}

func (this *GoLogger) Info(output string){
	this.Logger.Output(this.calldepth,fmt.Sprintf(" INFO: %s",output))
}

func (this *GoLogger) Error(output string){
	this.Logger.Output(this.calldepth,fmt.Sprintf(" ERROR: %s",output))
}
