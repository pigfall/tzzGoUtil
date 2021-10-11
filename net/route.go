package net

import(
		"net"
)



type RouteRule struct {
	TargetIpNet *IpWithMask
	DevName string 
	DevIndex int
	Via net.IP
}

func NewRouteRule(targetIp *IpWithMask,devName string,devIndex int,via net.IP)*RouteRule{
	return &RouteRule {
		// ethier targetIp or targetIpNet
		TargetIpNet :targetIp,
		DevName :devName,
		DevIndex :devIndex,
		Via :via,
	}
}
