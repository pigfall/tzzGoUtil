package net

import (
	"fmt"
	"net"

	"github.com/pigfall/tzzGoUtil/process"
	nl "github.com/vishvananda/netlink"
)

func AddRoute(target net.IP, devName string, via net.IP) error {
	return addRoute(target.String(), devName, via)
}

func AddRouteIpNet(target *IpWithMask, devName string, via net.IP) error {
	return addRoute(target.FormatAsIpSlashMask(), devName, via)
}

func DelRoute(target net.IP) error {
	_, errOut, err := process.ExeOutput("ip", "route", "del", target.String())
	if err != nil {
		return fmt.Errorf("Del route to %s failed %v, %v", target.String(), err, errOut)
	}
	return nil
}

func addRoute(target string, devName string, via net.IP) error {
	var errOut string
	var err error
	if via != nil {
		_, errOut, err = process.ExeOutput("ip", "route", "add", target, "via", via.String(), "dev", devName)
	} else {
		_, errOut, err = process.ExeOutput("ip", "route", "add", target, "dev", devName)
	}

	if err != nil {
		return fmt.Errorf("%s, %w", errOut, err)
	}

	return err
}

func RouteList() ([]*RouteRule, error) {
	routeRules, err := nl.RouteList(nil, nl.FAMILY_ALL)
	if err != nil {
		return nil, err
	}
	ifces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	ifceMap := make(map[int]net.Interface)
	for _, ifce := range ifces {
		ifceMap[ifce.Index] = ifce
	}

	rules := make([]*RouteRule, 0, len(routeRules))
	for _, rule := range routeRules {
		var ifceIndex = rule.LinkIndex
		dev, ok := ifceMap[ifceIndex]
		if !ok {
			return nil, fmt.Errorf("Not found net inteface by ifce index %d in %v", ifceIndex, ifceMap)
		}
		rules = append(rules, NewRouteRule(
			IpWithMaskFromIpNet(rule.Dst), dev.Name, rule.LinkIndex, rule.Gw,
		))
	}

	return rules, nil
}
