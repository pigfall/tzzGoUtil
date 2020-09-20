package net

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type EthernetFrame struct {
	frame gopacket.Packet
}

func DecodeICMPV4(toDecode []byte) (*layers.ICMPv4, error) {
	frame, err := DecodeEthernetFrame(toDecode)
	if err != nil {
		return nil, err
	}
	icmpLayerI := frame.frame.Layer(layers.LayerTypeICMPv4)
	if icmpLayerI == nil {
		return nil, fmt.Errorf("not a icmp")
	}
	return icmpLayerI.(*layers.ICMPv4), nil
}

func DecodeEthernetFrame(toDecode []byte) (*EthernetFrame, error) {
	pac := gopacket.NewPacket(toDecode, layers.LayerTypeEthernet, gopacket.Default)
	return &EthernetFrame{
		frame: pac,
	}, nil
}
