package io

import(
    "testing"
    "unsafe"
    "os"
	"golang.org/x/sys/unix"
    "github.com/Peanuttown/tzzGoUtil/process"
)

func TestIOCtl(t *testing.T){
    // create a tap device
    var devName="testDev"
    var ifReq [unix.IFNAMSIZ+64]byte
    copy(ifReq[:],[]byte(devName))
    *(*uint16)(unsafe.Pointer(&ifReq[unix.IFNAMSIZ]))=unix.IFF_TAP
    file,err :=os.OpenFile("/dev/net/tun",os.O_RDWR,0)
    if err != nil{
        t.Fatal(err)
    }
    err = IOCtl(
        file.Fd(),
        unix.TUNSETIFF,
        uintptr(unsafe.Pointer(&ifReq[0])),
    )
    if err != nil{
        t.Fatal(err)
    }
    // check device
    _,_,err = process.ExeOutput("ip","link","show",devName)
    if err != nil{
        t.Fatal(err)
    }


}
