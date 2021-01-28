package path

import (
	stdpath "path"
)

// 和 stdpath 中不一的是, 不会去 Clean
func DirPath(filepath string) string {
	dir, _ := stdpath.Split(filepath)
	return dir
}
