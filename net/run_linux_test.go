package net

import(
		"testing"
)

func TestNewTun(t *testing.T){
	const tunName = "testTun"
	tun,err := NewTun(tunName)
	if err != nil{
		t.Fatal(err)
	}
	t.Log(tun)
}
