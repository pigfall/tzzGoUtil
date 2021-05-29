package syscall

import (
	"fmt"
	libnl "github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"os"
	gosys "syscall"
)

const (
	RECEIVE_BUFFER_SIZE = 65536
)

func Syscall() {
	fd, _, e := gosys.RawSyscall(gosys.SYS_SOCKET, uintptr(gosys.AF_NETLINK), uintptr(gosys.SOCK_RAW), uintptr(gosys.NETLINK_ROUTE))
	if e != 0 {
		fmt.Println(e.Error())
		os.Exit(1)
	}
	s := &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
	}
	for _, g := range []int{unix.RTNLGRP_IPV4_ROUTE, unix.RTNLGRP_IPV6_ROUTE} {
		s.Groups |= (1 << (g - 1))
	}

	err := unix.Bind(int(fd), s)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	go func() {
		notification := make(chan libnl.RouteUpdate, 1)
		done := make(chan struct{})
		err := libnl.RouteSubscribe(notification, done)
		if err != nil {
			panic(err)
		}
		for {
			msg := <-notification
			fmt.Println(msg)
		}
	}()
	var rb [RECEIVE_BUFFER_SIZE]byte
	for {
		n, _, err := unix.Recvfrom(int(fd), rb[:], 0)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(rb[:n])
	}

}
