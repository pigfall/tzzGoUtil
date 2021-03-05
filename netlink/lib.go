package netlink

import (
	libnl "github.com/vishvananda/netlink"
)

func handleErr()

func Test() {
	notification := make(chan libnl.AddrUpdate, 1)
	done := make(chan struct{})
	err := libnl.RouteSubscribe(notification, done)
	libnl.NeighSubscribe
	libnl.FilterList()
	if err != nil {

	}

}
