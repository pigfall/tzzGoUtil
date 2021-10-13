package net

import(
		"testing"
		"net"
)


func TestIpWithMaskBaseIpNet(t *testing.T){
	_,ipNet,err := net.ParseCIDR("192.168.127.1/17")
	if err != nil{
		t.Fatal(err)
	}
	t.Log(ipNet)
}

func TestIpWithMaskForEach(t *testing.T){
	ipNet,err := FromIpSlashMask("192.168.0.1/24")
	if err != nil{
		t.Fatal(err)
	}
	var addr = make([]string,0)
	ipNet.ForEachIpInThisCidr(
		func(ipNet *IpWithMask)(stop bool, err error){
			addr =append(addr,ipNet.FormatAsIpSlashMask())
			return false,nil
		},
	)
	if len(addr) != 254{
		t.Fatalf("not expected addr num %d, %v",len(addr),addr)

	}
}

func TestMaskFormatTo255(t *testing.T){
	ipMask,err := FromIpSlashMask("10.8.0.1/24")
	if err != nil{
		t.Fatal(err)
	}
	s := MaskFormatTo255(ipMask.Mask)
	if s != "255.255.255.0"{
		t.Fatal("failed")
	}
}
