package main

import (
	"runtime/debug"
	"time"
	"log"
	"fmt"
	"github.com/pigfall/tzzGoUtil/net/wintun"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket"
)

func main() {
	defer func(){
		e := recover()
		if e != nil {
			fmt.Println(e)
			debug.PrintStack()
			time.Sleep(time.Hour)
		}
		time.Sleep(time.Second*100)
	}()
	tunIfce,err := wintun.NewTun("tzzTest",1500)
	if err != nil{
		panic(err)
	}
	fmt.Println(tunIfce)
	err = tunIfce.SetIp("192.168.2.1/24")
	if err != nil{
		panic(err)
	}
	var buf = make([]byte,1024*4)
	fmt.Println("Reading")
	for{
		n,err :=tunIfce.Read(buf)
		if err != nil{
			panic(err)
		}
		//log.Println(string(buf[:n]))
		pac := gopacket.NewPacket(buf[:n],layers.LayerTypeEthernet,gopacket.Default)
		if ethLayer := pac.Layer(layers.LayerTypeEthernet);ethLayer!=nil{
			log.Println("ethernet packet ")
		}else{
			log.Println("not ethernet packet")
		}

	}
}
