package net

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/pigfall/tzzGoUtil/sys/io"
	"golang.org/x/sys/unix"
)

type TunTapBase struct {
	DevFile *os.File
}

func (this *TunTapBase) Read(b []byte) (n int, err error) {
	return this.DevFile.Read(b)
}

func (this *TunTapBase) Write(b []byte) (n int, err error) {
	return this.DevFile.Write(b)
}

func (this *TunTapBase) Close(){
	this.DevFile.Close()
}

func newTunTapBase(file *os.File) *TunTapBase {
	return &TunTapBase{
		DevFile: file,
	}
}

func newTunTap(tpe VIRTUAL_DEV, devName string) (fd uintptr, err error) {
	var flag int16
	switch tpe {
	case TAP:
		flag = unix.IFF_TAP
	case TUN:
		flag = unix.IFF_TUN
	default:
		return 0, fmt.Errorf("unknown tun tap type %d", tpe)
	} // open /dev/net/tun
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return 0, err
	}

	var ifr [unix.IFNAMSIZ + 64]byte
	copy(ifr[:unix.IFNAMSIZ], []byte(devName))
	*(*int16)(unsafe.Pointer(&ifr[unix.IFNAMSIZ])) = flag | unix.IFF_NO_PI
	// ioctl
	err = io.IOCtl(
		file.Fd(),
		unix.TUNSETIFF,
		(uintptr)(unsafe.Pointer(&ifr[0])),
	)
	if err != nil {
		return 0, err
	}
	return file.Fd(), nil
}
