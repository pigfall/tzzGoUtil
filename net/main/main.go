package main

import (
	"bufio"
	"errors"
	"fmt"
	//"github.com/pigfall/tzzGoUtil/terminal/ansi"
	"io"
	"os/exec"
)

func main() {
	cmd := exec.Command("you-get.exe", "-f", "https://www.bilibili.com/video/BV1Tt411776s?from=search&seid=9333383915621586500")
	rd, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	defer rd.Close()
	bufReader := bufio.NewReaderSize(rd, 2048)
	for {
		line, err := bufReader.ReadBytes('%')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			panic(err)
		}
		fmt.Println(string(line[len(line)-5:]))
	}
}
