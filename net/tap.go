package net

import (
	"os"
)

type Tap struct {
	TunTapBase
}

func NewTap(devName string) (ITap, error) {
	fd, err := newTunTap(TAP, devName)
	if err != nil {
		return nil, err
	}
	return &Tap{
		TunTapBase: *newTunTapBase(os.NewFile(fd, "tap")),
	}, nil
}
