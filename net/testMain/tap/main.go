package main

import (
	"bufio"
	"net"
	"os"

	"github.com/Peanuttown/tzzGoUtil/log"
	tzNet "github.com/Peanuttown/tzzGoUtil/net"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func main() {
	var devName = "testDev"
	tap, err := tzNet.NewTap(devName)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	err = tzNet.SetIp(devName, "172.16.2.151/16")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	err = tzNet.DevUp(devName)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	var b = make([]byte, 1024)
	for {
		n, err := tap.Read(b)
		if err != nil {
			log.Error(err)
			break
		}
		handleData(b[:n], tap)
	}
}

func handleData(b []byte, tap tzNet.ITap) {
	// parse packet
	var err error
	packet := gopacket.NewPacket(b, layers.LayerTypeEthernet, gopacket.Default)
	if packet.ErrorLayer() != nil {
		log.Error(packet.ErrorLayer())
		return
	}

	ethernet := extractEthernet(packet)
	//ipframe :=extractIpFrame(packet)

	log.Debug("old dst mac", ethernet.DstMAC.String())
	log.Debug("old src mac", ethernet.SrcMAC.String())
	ethernet.DstMAC, err = net.ParseMAC("00:00:00:00:00:00")
	if err != nil {
		log.Error(err)
		return
	}
	buffer := gopacket.NewSerializeBuffer()
	option := gopacket.SerializeOptions{}
	err = ethernet.SerializeTo(buffer, option)
	if err != nil {
		log.Error(err)
		return
	}
	var toSvrEthernet layers.Ethernet
	err = (toSvrEthernet).DecodeFromBytes(buffer.Bytes(), gopacket.NilDecodeFeedback)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("new dst mac", toSvrEthernet.DstMAC.String())
	log.Debug("new src mac", toSvrEthernet.SrcMAC.String())

	var sBuffer = gopacket.NewSerializeBuffer()
	err = gopacket.SerializeLayers(sBuffer, option, &toSvrEthernet)
	if err != nil {
		log.Error(err)
		return
	}
	_, err = tap.Write(sBuffer.Bytes())
	if err != nil {
		log.Error(err)
	}
}

func extractEthernet(packet gopacket.Packet) *layers.Ethernet {
	l := packet.Layer(layers.LayerTypeEthernet)
	return l.(*layers.Ethernet)
}

func extractIpFrame(packet gopacket.Packet) *layers.IPv4 {
	l := packet.Layer(layers.LayerTypeEthernet)
	return l.(*layers.IPv4)
}

func handleTapData(tap tzNet.ITap) {
	var b = make([]byte, 1024)
	for {
		n, err := tap.Read(b)
		if err != nil {
			log.Error("读取 tap 中的数据失败: %w", err)
			return
		}
		pac := gopacket.NewPacket(b[:n], layers.LayerTypeEthernet, gopacket.Default)
		if errLayer := pac.ErrorLayer(); errLayer != nil {
			log.Error(errLayer)
			continue
		}
		// change mac address
		ethetnetFrame, err := changeMAC(pac)
		if err != nil {
			log.Error(err)
		}
		log.Debug(ethetnetFrame)
	}
}

func changeMAC(packet gopacket.Packet) (gopacket.Packet, error) {
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	ethernet := ethernetLayer.(*layers.Ethernet)
	loopMac, err := net.ParseMAC("00:00:00:00:00:00")
	if err != nil {
		return nil, err
	}
	ethernet.DstMAC = loopMac
	serializeBuffer := gopacket.NewSerializeBuffer()
	option := gopacket.SerializeOptions{}
	err = ethernet.SerializeTo(serializeBuffer, option)
	if err != nil {
		return nil, err
	}
	var newEthernet = &layers.Ethernet{}
	err = newEthernet.DecodeFromBytes(serializeBuffer.Bytes(), gopacket.NilDecodeFeedback)
	if err != nil {
		return nil, err
	}
	err = gopacket.SerializeLayers(serializeBuffer, option, newEthernet)
	if err != nil {
		return nil, err
	}
	return gopacket.NewPacket(serializeBuffer.Bytes(), layers.LayerTypeEthernet, gopacket.Default), nil
}

func servce(l net.Listener, finish chan bool) {
	defer func() {
		finish <- true
	}()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error(err)
			return
		}
		// print remote address
		log.Debugf("与 %s 建立起连接\n", conn.RemoteAddr().String())
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	rd := bufio.NewReader(conn)
	var buf = make([]byte, 1024)
	for {
		n, err := rd.Read(buf)
		if err != nil {
			log.Error(err)
			return
		}
		log.Debug("Read data : %s", string(buf[:n]))
	}

}
