package net

import "os"

type Tun struct {
	TunTapBase
}


func NewTun(devName string) (ITun, error) {
	fd, err := newTunTap(TUN, devName)
	if err != nil {
		return nil, err
	}
	return &Tun{
		TunTapBase: *newTunTapBase(os.NewFile(fd, "tun")),
	}, nil
}
