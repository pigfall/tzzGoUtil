package main

import (
	"fmt"
	"github.com/Peanuttown/tzzGoUtil/debug"
)

func Caller() string {
	return debug.CallerName()
}

func main() {
	fmt.Println(debug.CallerName())
	fmt.Println(Caller())
}
