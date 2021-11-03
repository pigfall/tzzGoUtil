package net

import(
	"net"
)

type UDPSock struct{
	*net.UDPConn
}

func UDPListen(ipToListen net.IP,port int)(*UDPSock,error){
	udpSock,err := net.ListenUDP("udp",&net.UDPAddr{
		IP:ipToListen,
		Port:port,
	})
	if err !=nil{
		return nil,err
	}
	return &UDPSock{
		UDPConn:udpSock,
	},nil
}
