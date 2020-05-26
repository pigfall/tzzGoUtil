package net

import (
	"sync"

	"github.com/Peanuttown/gopacket"
	"github.com/Peanuttown/gopacket/layers"
	"github.com/Peanuttown/gopacket/pcap"
	"github.com/Peanuttown/tzzGoUtil/log"
)

type Packet struct {
	gopacket.Packet
}

func Listen(netInterface string) (packetsChRet <-chan Packet, closeF func(), err error) {
	//var devices []pcap.Interface
	//devices, err = pcap.FindAllDevs()
	//if err != nil {
	//	return
	//}
	var handler *pcap.Handle
	handler, err = pcap.OpenLive(netInterface, 65536, true, pcap.BlockForever)
	if err != nil {
		return
	}
	closeSig := make(chan struct{})
	src := gopacket.NewPacketSource(handler, layers.LayerTypeEthernet)
	innerpacketCh := src.Packets()
	packetsCh := make(chan Packet, 100)
	wg := sync.WaitGroup{}
	wg.Add(1)
	closeF = func() {
		close(closeSig)
		log.Debug("listen wait")
		wg.Wait()
	}
	go func() {
		defer func() {
			handler.Close()
			log.Debug("listen done")
			wg.Done()
		}()
		for {
			select {
			case <-closeSig:
				close(packetsCh)
				return
			case packet := <-innerpacketCh:
				p := Packet{}
				p.Packet = packet
				packetsCh <- p
			}
		}
	}()
	packetsChRet = packetsCh
	return
}
