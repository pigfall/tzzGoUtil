package io


import(
    "fmt"
	"golang.org/x/sys/unix"
)

/*
The fd argument is an open file descriptor for the device or file upon which the
control operation specified by request is to be performed. Device-specific header
files define constants that can be passed in the request argument

the third argument toioctl(), which we label argp, can be of 
any type. The value of the request argumentenables ioctl() to determin
e what type of value to expect in argPointer. Typically, argp is apointer to
 either an integer or a structure; in some cases, it is unused
*/
func IOCtl (fd uintptr,req int,argPointer uintptr)(error){
    _,_,errNo := unix.Syscall(
        unix.SYS_IOCTL,
        fd,
        uintptr(req),
        argPointer,
    )
    if errNo != 0{
        return fmt.Errorf("ioctl failed: %w",errNo)
    }
    return nil
}



