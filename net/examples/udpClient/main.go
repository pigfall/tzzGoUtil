package main

import (
	"github.com/pigfall/tzzGoUtil/net"
)

func main() {
	remoteIp,err := net.ParseIPv4("107.155.15.21")
	if err != nil{
		panic(err)
	}
	sock,err := net.UDPDial(remoteIp,10101)
	if err != nil{
		panic(err)
	}
	_,err = sock.Write([]byte("test"))
	if err != nil{
		panic(err)
	}
}
