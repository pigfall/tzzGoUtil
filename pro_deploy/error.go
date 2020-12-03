package main

import (
	"github.com/Peanuttown/tzzGoUtil/output"
	"os"
)

func handleErr(err error) {
	if err != nil {
		output.Err(err)
		output.Err("❌ ❌  部署失败")
		os.Exit(1)
	}
}
