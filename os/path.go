package os

import (
	stdos "os"
	"path/filepath"
)

// 获取进程的执行路径
func GetExecutablePath() (string, error) {
	p, err := stdos.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(p)
}
