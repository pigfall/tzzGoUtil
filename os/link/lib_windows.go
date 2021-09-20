package link

import (
	"fmt"
	"github.com/pigfall/tzzGoUtil/sys/windows"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const proc_name = "CreateSymbolicLinkA"

var symbolic_proc_addr uintptr
var load_lock sync.Mutex

func Symbolic(from, to string) error {
	load_lock.Lock()
	defer load_lock.Unlock()
	if symbolic_proc_addr == 0 {
		user32_lib_handle_addr, err := syscall.LoadLibrary(windows.LIB_NAME_USER32)
		if err != nil {
			return fmt.Errorf("Load library %s failed", windows.LIB_NAME_USER32, err)
		}
		defer syscall.FreeLibrary(user32_lib_handle_addr)
		symbolic_proc_addr, err = syscall.GetProcAddress(user32_lib_handle_addr, proc_name)
		if err != nil {
			symbolic_proc_addr = 0
			return fmt.Errorf("GetProcAddress failed: %w", err)
		}
	}
	ret,_,callErr := syscall.Syscall(
		symbolic_proc_addr,
		uintptr(syscall.StringToUTF16(to)),
		uintptr(syscall.StringToUTF16(from)),
		0,
	)
	if callErr != 0 {
		return fmt.Errorf(callErr.Error())
	}
	if ret 
}
