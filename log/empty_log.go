package log

type EmptyLogger struct{}

func NewEmptyLogger()*EmptyLogger{
	return &EmptyLogger{}
}



func(this *EmptyLogger)	Error(msg ...interface{}){

}
func(this *EmptyLogger)	Errorf(format string, msg ...interface{}){}
func(this *EmptyLogger)	Debug(msg ...interface{}){}
func(this *EmptyLogger)	Debugf(format string, args ...interface{}){}
func(this *EmptyLogger)	Info(msg ...interface{}){}
func(this *EmptyLogger)	Infof(format string, msg ...interface{}){}
