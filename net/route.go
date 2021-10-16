package net

import (
	"fmt"
	"net"
)

type RouteRule struct {
	TargetIpNet *IpWithMask
	DevName     string
	DevIndex    int
	Via         net.IP
}

func NewRouteRule(targetIp *IpWithMask, devName string, devIndex int, via net.IP) *RouteRule {
	if targetIp != nil {
		if targetIp.String() == "0.0.0.0/0" {
			targetIp = nil
		}

	}
	return &RouteRule{
		// ethier targetIp or targetIpNet
		TargetIpNet: targetIp,
		DevName:     devName,
		DevIndex:    devIndex,
		Via:         via,
	}
}
func GetDefaultRouteRule() (*RouteRule, error) {
	rules, err := RouteList()
	if err != nil {
		return nil, err
	}

	for _, rule := range rules {
		if rule.TargetIpNet == nil {
			return rule, nil
		}
	}
	return nil, fmt.Errorf("Not found default route")
}
