package winsys

import (
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
		t.Logf("dst %s", row.Dst().FormatAsIpSlashMask())
		t.Logf("gateway %v", row.Gateway())
		t.Logf("ifceIndex %v", row.SrcDevIndex())
	}
}
