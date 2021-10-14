package main

import (
	gonet "net"
	"os"

	"github.com/pigfall/tzzGoUtil/log"
	"github.com/pigfall/tzzGoUtil/net"
	"github.com/pigfall/tzzGoUtil/net/gopacketUtils"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func main() {
	var devName = "testDev"
	tun, err := net.NewTun(devName)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}


	ipWithMask,err := net.FromIpSlashMask("10.8.0.1/8")
	if err != nil{
		log.Error(err)
		os.Exit(1)
	}
	err = net.SetIp(devName, *ipWithMask)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	err = net.DevUp(devName)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	var b = make([]byte, 1024)
	for {
		n, err := tun.Read(b)
		if err != nil {
			log.Error(err)
			return
		}
		replyV2(b[:n], tun)
	}
}

func replyV2(b []byte, tun net.ITun) {
	packet := gopacket.NewPacket(b, layers.LayerTypeIPv4, gopacket.Default)
	var rep = make([]gopacket.SerializableLayer, 0)
	if iplayer, ok := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4); ok {
		log.Debug("src ", iplayer.SrcIP.String())
		log.Debug("dst ", iplayer.DstIP.String())
		//var dst = iplayer.DstIP
		//var src = iplayer.SrcIP
		//iplayer.DstIP = src
		//iplayer.SrcIP = ds
		if iplayer.DstIP.String() =="10.8.0.2"{
			log.Debug("Change dst ip")
			iplayer.DstIP = gonet.ParseIP("8.8.8.8")
		}
		rep = append(rep, iplayer)
	} else {
		return
	}
	allLayers := packet.Layers()
	var startAppend bool
	for _,layer := range  allLayers {
		if startAppend{
			if sLayer,ok := layer.(gopacket.SerializableLayer);ok{
				log.Debug("Is SerializableLayer")
				rep = append(rep,sLayer)
			}
		}else{
			if layer.LayerType() == layers.LayerTypeIPv4{
				startAppend = true
			}
		}
	}
	buffer := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{}, rep...)
	if err != nil {
		log.Error(err)
		return
	}
	n, err := tun.Write(buffer.Bytes())
	if err != nil {
		log.Error(err)
	}
	log.Debugf("write %d data\n", n)
}

func reply(b []byte, tun net.ITun) {
	packet := gopacket.NewPacket(b, layers.LayerTypeIPv4, gopacket.Default)
	var rep = make([]gopacket.SerializableLayer, 0)
	if iplayer, ok := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4); ok {
		log.Debug("src ", iplayer.SrcIP.String())
		log.Debug("dst ", iplayer.DstIP.String())
		var dst = iplayer.DstIP
		var src = iplayer.SrcIP
		iplayer.DstIP = src
		iplayer.SrcIP = dst
		rep = append(rep, iplayer)
	} else {
		return
	}
	var payload gopacket.Payload
	if icmpLayer, ok := packet.Layer(layers.LayerTypeICMPv4).(*layers.ICMPv4); ok {
		icmpLayer.TypeCode = 0
		log.Debug("icmp payload length :", len(icmpLayer.Payload))
		rep = append(rep, icmpLayer)
		payload = gopacket.Payload(icmpLayer.Payload)
		rep = append(rep, payload)
	} else {
		return
	}
	buffer := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{}, rep...)
	if err != nil {
		log.Error(err)
		return
	}
	n, err := tun.Write(buffer.Bytes())
	if err != nil {
		log.Error(err)
	}
	log.Debugf("write %d data\n", n)
}

func handleData(b []byte, tun net.ITun) {
	packet := gopacket.NewPacket(b, layers.LayerTypeIPv4, gopacket.Default)
	if packet.ErrorLayer() != nil {
		log.Error(packet.ErrorLayer().Error())
		return
	}
	//for _, layer := range packet.Layers() {
	//	log.Debug(layer.LayerType().String())
	p, err := gopacketUtils.ChangeDstIp(packet, gonet.IPv4(127, 0, 0, 1))
	if err != nil {
		log.Error(err)
		return
	}

	n, err := tun.Write(p.Data())
	if err != nil {
		log.Error(err)
	}
	log.Debug("转发数据 %s ", n)
}
