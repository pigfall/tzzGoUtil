package winsys

import (
	"encoding/binary"
	"net"
)

type MIB_IPFORWARDROW struct {
	DwForwardDest      DWORD
	DwForwardMask      DWORD
	DwForwardPolicy    DWORD
	DwForwardNextHop   DWORD
	DwForwardIfIndex   IF_INDEX
	ForwardType        MIB_IPFORWARD_TYPE
	ForwardProto       MIB_IPFORWARD_PROTO
	DwForwardAge       DWORD
	DwForwardNextHopAS DWORD
	DwForwardMetric1   DWORD
	DwForwardMetric2   DWORD
	DwForwardMetric3   DWORD
	DwForwardMetric4   DWORD
	DwForwardMetric5   DWORD
}

func (this *MIB_IPFORWARDROW) Dst() net.IPNet {
	ip := make([]byte, 4)
	binary.LittleEndian.PutUint32(ip, uint32(this.DwForwardDest))
	ipObj := net.IPv4(ip[0], ip[1], ip[2], ip[3])
	mask := make([]byte, 4)
	binary.LittleEndian.PutUint32(mask, uint32(this.DwForwardMask))
	return net.IPNet{
		IP:   ipObj,
		Mask: net.IPMask(mask),
	}
}

func (this *MIB_IPFORWARDROW) Gateway() net.IP {
	gateway := this.DwForwardNextHop
	if gateway == 0 {
		return nil
	}

	ip := make([]byte, 4)
	binary.LittleEndian.PutUint32(ip, uint32(gateway))
	return net.IPv4(ip[0], ip[1], ip[2], ip[3])
}

func (this *MIB_IPFORWARDROW) SrcDevIndex() uint32 {
	return uint32(this.DwForwardIfIndex)
}
