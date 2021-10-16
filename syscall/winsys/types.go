package winsys

type DWORD uint32
type ULONG uint32
type NL_ROUTE_PROTOCOL int32
type IF_INDEX NET_IFINDEX
type NET_IFINDEX ULONG
type MIB_IPFORWARD_TYPE int32
type MIB_IPFORWARD_PROTO NL_ROUTE_PROTOCOL

const ANY_SIZE = 1

type MIB_IPFORWARDTABLE struct {
	DwNumEntries DWORD
	Table        []MIB_IPFORWARDROW
}

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

type PMIB_IPFORWARDTABLE *MIB_IPFORWARDTABLE
