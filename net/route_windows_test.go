package net

import (
	"testing"

	"github.com/pigfall/tzzGoUtil/syscall/winsys"
)

func TestGetDefaultRoute(t *testing.T) {
	err := winsys.LoadIpHelperDLL()
	if err != nil {
		t.Fatal(err)
	}
	rule, err := GetDefaultRouteRule()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rule)
}
