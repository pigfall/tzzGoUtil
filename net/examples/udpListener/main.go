package main

import (
	"log"
	"os"
	"github.com/pigfall/tzzGoUtil/net"
)

func main() {
	ipToListen,err := net.ParseIPv4("0.0.0.0")
	if err != nil{
		panic(err)
	}
	const port  = 10101
	log.Printf("Listen at %v %v",ipToListen.String(),port)
	l,err :=net.UDPListen(ipToListen,port)
	if err != nil{
		log.Printf("Failed to listen at %v port %v\n",ipToListen.String(),port)
		os.Exit(1)
	}
	buffer := make([]byte,4*2014)
	readN,remoteAddr,err := l.ReadFromUDP(buffer)
	if err != nil{
		log.Println("ReadFromUDP　failed: ",err)
		os.Exit(1)
	}
	log.Printf("ReadFromUDP remote addr %v, content: %v\n",remoteAddr.String(),string(buffer[:readN]))
}
