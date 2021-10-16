package winsys

import (
	"encoding/binary"
	"reflect"
	"unsafe"

	gosyscall "syscall"

	"github.com/pigfall/tzzGoUtil/syscall"
	"golang.org/x/sys/windows"
)

var ipHelperDLL *syscall.DLL

func LoadIpHelperDLL() error {
	var err error
	ipHelperDLL, err = syscall.LoadDLL("iphlpapi.dll")
	if err != nil {
		return err
	}
	return nil
}

//func GetIpForwardTable(pIpForwardTable PMIB_IPFORWARDTABLE, pdwSize *uint32, bOrder bool) (DWORD, error) {
func GetIpForwardTable() ([]*MIB_IPFORWARDROW, error) {
	procdure, err := ipHelperDLL.FindProc("GetIpForwardTable")
	if err != nil {
		return nil, err
	}
	var bufSize uint
	var buf = make([]byte, 1)
	for {
		r1, _, _ := gosyscall.Syscall(
			procdure.Addr(),
			3,
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&bufSize)),
			getUintptrFromBool(true),
		)
		if r1 == windows.NO_ERROR {
			// { parse entry num
			entryNum := binary.LittleEndian.Uint32(buf[:4])
			// {{ parse rows
			return parseRouteTableRows(int(entryNum), buf[4:])
			// }}
			// }
		} else if r1 == uintptr(windows.ERROR_INSUFFICIENT_BUFFER) {
			buf = make([]byte, bufSize)
		} else {
			panic(r1)
		}
	}
}

func parseRouteTableRows(entryNum int, buf []byte) ([]*MIB_IPFORWARDROW, error) {
	row := MIB_IPFORWARDROW{}
	rowSize := reflect.TypeOf(row).Size()
	elems := make([]*MIB_IPFORWARDROW, 0, entryNum)
	for i := 0; i < entryNum; i++ {
		row, err := parseRouteTableRow(buf[(int(i))*int(rowSize):])
		if err != nil {
			return nil, err
		}
		elems = append(elems, row)
	}
	return elems, nil
}

func parseRouteTableRow(buf []byte) (*MIB_IPFORWARDROW, error) {
	row := &MIB_IPFORWARDROW{}
	rv := reflect.ValueOf(row)
	rvElem := rv.Elem()
	var index = 0
	for i := 0; i < rvElem.Type().NumField(); i++ {
		v := binary.LittleEndian.Uint32(buf[index : index+4])
		// TODO
		if rvElem.Field(i).Type().Kind() == reflect.Uint32 {
			rvElem.Field(i).SetUint(uint64(v))
		} else {
			rvElem.Field(i).SetInt(int64(v))
		}
		index += 4
	}
	return row, nil
}

func getUintptrFromBool(b bool) uintptr {
	if b {
		return 1
	} else {
		return 0
	}
}
