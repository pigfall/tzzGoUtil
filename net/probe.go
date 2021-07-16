package net

import(
		stdnet "net"
)



type ProbeNetWork string

const(
		TCP ProbeNetWork= "tcp"
		UDP ProbeNetWork= "udp"
)

func ProbeAddr(network ProbeNetWork,addr string)(error){
	l,err := stdnet.Dial(string(network),addr)
	if err != nil{
		return  err
	}
	l.Close()
	return nil
}
