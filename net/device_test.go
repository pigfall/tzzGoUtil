package net

import(
		"testing"
)



func TestListDevices(t *testing.T){
	devs,err := ListDevices()
	if err != nil{
		t.Fatal(err)
	}
	for _,dev := range devs{
		addrs,err := dev.Addrs()
		if err != nil{
			t.Fatal(err)
		}
		for _,addr := range addrs{
			t.Log(addr.String())
		}
	}
}
