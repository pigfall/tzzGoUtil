package winsys

import (
	"encoding/binary"
	"net"
	"testing"
)

func TestGetIpForwardTable(t *testing.T) {
	err := LoadIpHelperDLL()
	if err != nil {
		t.Fatal(err)
	}

	//var table MIB_IPFORWARDTABLE
	rows, err := GetIpForwardTable()
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		dst := make([]byte, 4)
		binary.LittleEndian.PutUint32(dst, uint32(row.DwForwardDest))
		t.Logf("dst %s", net.IPv4(dst[0], dst[1], dst[2], dst[3]).String())
	}
}
