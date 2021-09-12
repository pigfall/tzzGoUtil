package water_wrap

import(
	"fmt"
	"testing"
"github.com/Peanuttown/tzzGoUtil/process"
)


func TestNewTun(t *testing.T){
	ifce,err := NewTun()
	if err != nil{
		t.Fatal(err)
	}
	t.Log("tun created")
	out,errOut,err := process.ExeOutput("ifconfig",ifce.Name(),"10.1.0.10","10.1.0.11","up")
	if err != nil{
		t.Fatal(fmt.Errorf("%w, %v, %v",err,out,errOut))
	}
	_,err = ifce.Write(send())
	if err != nil{
		t.Fatal(err)
	}
}
