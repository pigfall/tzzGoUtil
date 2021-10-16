package net

import (
	"fmt"
	"net"

	"github.com/pigfall/tzzGoUtil/process"
	"github.com/pigfall/tzzGoUtil/syscall/winsys"
)

func AddRoute(target net.IP, devIndex string, via net.IP) error {
	if !IsIpv4(target) {
		panic("TODO ,only supported ipv4")
	}
	return addRoute(target, net.IPv4Mask(255, 255, 255, 255), devIndex, via)
}

// route ADD destination_network MASK subnet_mask  gateway_ip metric_cost
func addRoute(dstIp net.IP, mask net.IPMask, devIndex string, gateway net.IP) error {
	var out, errOut string
	var err error
	if gateway == nil {
		out, errOut, err = process.ExeOutput("route", "add", dstIp.String(), "MASK", MaskFormatTo255(mask), "IF", devIndex)
	} else {
		out, errOut, err = process.ExeOutput("route", "add", dstIp.String(), "MASK", MaskFormatTo255(mask), "gateway", gateway.String(), "IF", devIndex)
	}
	if err != nil {
		return fmt.Errorf("%v, %v,%v", err, out, errOut)
	}
	return nil
}

func AddRouteIpNet(target *IpWithMask, devIndex string, via net.IP) error {
	return addRoute(target.Ip, target.Mask, devIndex, via)
}

func DelRoute(target net.IP) error {
	panic("TODO")
}

func RouteList() ([]*RouteRule, error) {
	rules, err := winsys.GetIpForwardTable()
	if err != nil {
		return nil, err
	}
	retRules := make([]*RouteRule, 0, len(rules))
	ifces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	ifceMap := make(map[int]net.Interface)
	for _, ifce := range ifces {
		ifceMap[ifce.Index] = ifce
	}
	for _, rule := range rules {
		var ifceIndex = int(rule.SrcDevIndex())
		dev, ok := ifceMap[ifceIndex]
		if !ok {
			return nil, fmt.Errorf("Not found net inteface by ifce index %d in %v", ifceIndex, ifceMap)
		}
		dstIpNet := rule.Dst()
		retRules = append(
			retRules,
			NewRouteRule(
				&IpWithMask{Ip: dstIpNet.IP, Mask: dstIpNet.Mask},
				dev.Name,
				int(ifceIndex),
				rule.Gateway(),
			),
		)
	}
	return retRules, nil
}
