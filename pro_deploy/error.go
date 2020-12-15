package main

import (
	"github.com/Peanuttown/tzzGoUtil/output"
	"os"
	"runtime/debug"
)

func handleErr(err error) {
	if err != nil {
		output.Err(err)
		output.Err("\n > STACK <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< ")
		output.Err(string(debug.Stack()))
		output.Err("\n < STACK >>>>>>>>>>>>>>>>>>>>>>>>>>>>> ")
		output.Err("❌ ❌  部署失败")
		os.Exit(1)
	}
}
