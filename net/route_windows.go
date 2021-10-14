package net

import(
		"net"
)

func AddRoute(target net.IP,devName string,via net.IP)error{
	 panic("TODO")
}

func AddRouteIpNet(target *IpWithMask,devName string,via net.IP)error{
	 panic("TODO")
}

func DelRoute(target net.IP)error{
	 panic("TODO")
}


func GetDefaultRouteRule()(*RouteRule,error){
	 panic("TODO")
}

func RouteList()([]*RouteRule,error){
	panic("TODO")
}

