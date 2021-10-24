package net

import(
		
	"net"
)

func IsIpv6(ip net.IP)bool{
	return ip.To4() == nil
}
