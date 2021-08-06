package log

import(
	"os"
		"testing"
)



func TestJsonLogger(t *testing.T){
	rawLogger := NewJsonLogger(os.Stdout)
	logger := NewHelper("test",rawLogger,LevelDebug)
	logger.Info("key","output msg")
}
