package process

import (
	"bytes"
	"io"
	"os/exec"
	"sync"
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

func ExecOutput(normalWriter io.Writer, errWriter io.Writer, cmdName string, args ...string) error {
	cmd := exec.Command(cmdName, args...)
	pipeErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer pipeErr.Close()
	pipeNormal, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer pipeNormal.Close()
	err = cmd.Start()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	var f = func(rd io.Reader, wd io.Writer) {
		defer func() {
			wg.Done()
		}()
		io.Copy(wd, rd)
	}
	go f(pipeErr, errWriter)
	go f(pipeNormal, normalWriter)
	wg.Wait()
	return cmd.Wait()
}

func ExecCombinedOutput(writer io.Writer, cmdName string, args ...string) error {
	cmd := exec.Command(cmdName, args...)
	pipeErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer pipeErr.Close()
	pipeNormal, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer pipeNormal.Close()
	err = cmd.Start()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	for _, v := range []io.Reader{pipeErr, pipeNormal} {
		go func(rd io.Reader) {
			defer func() {
				wg.Done()
			}()
			io.Copy(writer, rd)
		}(v)
	}
	wg.Wait()
	return cmd.Wait()
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
