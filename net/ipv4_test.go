package net

import(
    "testing"
)

func TestParseIpv4(t *testing.T){
    ip := "1.1.1.1"
    v,err := ParseIPv4(ip)
    if err != nil{
        t.Fatal(err)
    }
    t.Logf("%v.%v.%v.%v",v[0],v[1],v[2],v[3])
}
