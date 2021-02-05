package os

import (
	"fmt"
	stdos "os"
	"path/filepath"
)

// 获取进程的执行路径
func GetExecutablePath() (string, error) {
	p, err := stdos.Executable()
	if err != nil {
		return "", fmt.Errorf("Get executable path error: %w", err)
	}
	ret, err := filepath.EvalSymlinks(p)
	if err != nil {
		return "", fmt.Errorf("Get executable path error: %w", err)
	}
	return ret, nil
}
