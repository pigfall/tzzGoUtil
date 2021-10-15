package net

import(
	"fmt"
		"net"
		"github.com/pigfall/tzzGoUtil/process"
)

func AddRoute(target net.IP,devIndex string,via net.IP)error{
	if !IsIpv4(target){
		panic("TODO ,only supported ipv4")
	}
	return addRoute(target,net.IPv4Mask(255,255,255,255),devIndex,via)
}

// route ADD destination_network MASK subnet_mask  gateway_ip metric_cost
func addRoute(dstIp net.IP,mask net.IPMask,devIndex string,gateway net.IP)error{
	var out,errOut string
	var err error
	if gateway == nil{
		out,errOut,err = process.ExeOutput("route","add",dstIp.String(),"MASK",MaskFormatTo255(mask),"IF",devIndex)
	}else{
		out,errOut,err =process.ExeOutput("route","add",dstIp.String(),"MASK",MaskFormatTo255(mask),"gateway",gateway.String(),"IF",devIndex)
	}
	if err != nil{
		return fmt.Errorf("%v, %v,%v",err,out,errOut)
	}
	return nil
}

func AddRouteIpNet(target *IpWithMask,devIndex string,via net.IP)error{
	return addRoute(target.Ip,target.Mask,devIndex,via)
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

