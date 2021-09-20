package net

import (
	"github.com/Peanuttown/tzzGoUtil/process"
	"github.com/pkg/errors"
)

type VIRTUAL_DEV int

const (
	TUN VIRTUAL_DEV = iota + 1
	TAP
)


func SetIp(devName string, ip IpWithMask) error {
	_, errOut, err := process.ExeOutput("ip", "addr", "add", string(ip), "dev", devName)
	if err != nil {
		return errors.Wrap(err, errOut)
	}
	return nil
}

func DevUp(devName string) error {
	_, errOut, err := process.ExeOutput("ip", "link", "set", devName, "up")
	if err != nil {
		return errors.Wrap(err, errOut)
	}
	return nil
}
