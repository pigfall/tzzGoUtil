package debug

import(
	"testing"
)

func Caller() string{
	return CallerName()
}

func TestCallerName(t *testing.T){
	name := Caller()
	if name!="Caller"{
		t.Error("test failed")
		t.Log(name)
		return
	}
}
