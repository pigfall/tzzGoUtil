package net

import (
	"fmt"
	"log"
	"net"
	"strconv"

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
	log.Println("devIndex ", devIndex)
	devIndexInt, err := strconv.ParseInt(devIndex, 10, 64)
	if err != nil {
		return err
	}
	netIfce, err := FindIfceByIndex(int(devIndexInt))
	if err != nil {
		return err
	}
	addrs, err := netIfce.Addrs()
	if err != nil {
		return err
	}
	var srcIpObj net.IP
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return err
		}
		if IsIpv4(ip) {
			srcIpObj = ip
			break
		}
	}

	if srcIpObj == nil {
		panic(fmt.Errorf("not found ipv4 on interfce %v", devIndex))
	}

	_ = srcIpObj.String()

	fmt.Println("addRoute")
	var out, errOut string
	if gateway == nil {
		cmds := []string{"route", "add", dstIp.String(), "MASK", MaskFormatTo255(mask), "0.0.0.0", "if", devIndex}
		log.Println(cmds)
		out, errOut, err = process.ExeOutput(cmds[0], cmds[1:]...)
	} else {
		cmds := []string{"route", "add", dstIp.String(), "MASK", MaskFormatTo255(mask), gateway.String(), "if", devIndex}
		log.Println(cmds)
		out, errOut, err = process.ExeOutput(cmds[0], cmds[1:]...)
	}
	if err != nil {
		return fmt.Errorf("%v, %v,%v", err, out, errOut)
	}
	return nil
}

func AddRouteIpNet(target *IpWithMask, devIndex string, via net.IP) error {
	log.Println("devIndex ", devIndex)
	return addRoute(target.Ip, target.Mask, devIndex, via)
}

func DelRoute(target net.IP) error {
	return nil
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
