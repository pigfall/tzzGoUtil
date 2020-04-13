package debug

import(
	"testing"
    "strings"
)

func Caller() string{
	return CallerName()
}

func TestCallerName(t *testing.T){
	name := Caller()
	if !strings.Contains(name,"Caller"){
        t.Fatalf("test failed,expect name:'Caller',get %v",name)
		return
	}
}
