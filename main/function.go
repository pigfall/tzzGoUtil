package main

import(
	"fmt"
	"github.com/Peanuttown/tzzGoUtil/debug"
	"github.com/Peanuttown/tzzGoUtil/log"
	linux "github.com/c9s/goprocinfo/linux"
)

func Caller() string{
	return debug.CallerName()
}

func main(){
	fmt.Println(debug.CallerName())
	fmt.Println(Caller())
    stat,err := linux.ReadStat(/proc/stat)
    if err != nil{
        log.Debug(err)
        return
    }
}
