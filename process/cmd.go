package process

import (
	"bytes"
	"io"
	"os/exec"
)

// 执行 命令 并返回 标准输出和标准错误输出
func ExeOutput(name string, args ...string) (out string, errOut string, err error) {
	cmd := exec.Command(name, args...)
	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer
	err = cmd.Run()
	return outBuffer.String(), errBuffer.String(), err
}

func ExecWithErrOutput(writer io.Writer, cmdName string, args ...string) error {
	cmd := exec.Command(cmdName, args...)
	pipeErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer pipeErr.Close()
	err = cmd.Start()
	if err != nil {
		return err
	}
	io.Copy(writer, pipeErr)
	return cmd.Wait()
}
