package log

import(
    golog "log"
    "os"
    "fmt"
)

const callDepth =2

var i *golog.Logger

func init(){
    i = golog.New(os.Stdout, "tzzLog",golog.Lshortfile|golog.LstdFlags)
}

func Debug(msg ...interface{}){
    i.Output(callDepth,fmt.Sprintln(msg...))
}

func Debugf(format string,args ...interface{}){
    i.Output(callDepth,fmt.Sprintf(format,args...))
}
