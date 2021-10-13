package main

import (
	"time"
	"fmt"
	"github.com/pigfall/tzzGoUtil/net/wintun"
)

func main() {
	tunIfce,err := wintun.NewTun("tzzTest",1500)
	if err != nil{
		panic(err)
	}
	fmt.Println(tunIfce)
	err = tunIfce.SetIp("192.168.2.1/24")
	if err != nil{
		fmt.Println(err)
	}
	time.Sleep(time.Hour)
}
