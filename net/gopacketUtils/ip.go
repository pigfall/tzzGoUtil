package gopacketUtils

import (
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func ChangeDstIp(packet gopacket.Packet, ip net.IP) (gopacket.Packet, error) {
	// get all layers
	allLayers := packet.Layers()
	replaceLayer := make([]gopacket.SerializableLayer, 0, len(allLayers))
	var hasIpLayer = false
	for _, v := range allLayers {
		if v.LayerType() == layers.LayerTypeIPv4 {
			hasIpLayer = true
			ipLayer := v.(*layers.IPv4)
			ipLayer.DstIP = ip
		}
		if l, ok := v.(gopacket.SerializableLayer); ok {
			replaceLayer = append(replaceLayer, l)
		}
	}
	if !hasIpLayer {
		return nil, fmt.Errorf("no ip layer")
	}
	buffer := gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{}, replaceLayer...)
	if err != nil {
		return nil, err
	}
	p := gopacket.NewPacket(buffer.Bytes(), layers.LayerTypeIPv4, gopacket.Default)
	if p.ErrorLayer() != nil {
		return nil, err
	}
	return p, nil
}
