package main

import (
	"fmt"
	"github.com/pigfall/tzzGoUtil/debug"
)

func Caller() string {
	return debug.CallerName()
}

func main() {
	fmt.Println(debug.CallerName())
	fmt.Println(Caller())
}
