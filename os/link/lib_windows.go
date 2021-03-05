package link

import (
	"fmt"
	"github.com/Peanuttown/tzzGoUtil/sys/windows"
	"syscall"
)

func Symbolic(from, to string) error {
	syscall.LoadLibrary(windows.LIB_NAME_USER32)
	fmt.Println("windows_version")
	return nil

}
