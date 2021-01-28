package filepath

import (
	std "path/filepath"
)

// 和 stdpath 中不一的是, 不会去 Clean
func DirPath(filepath string) string {
	dir, _ := std.Split(filepath)
	return dir
}

func Join(elems ...string) string {
	return std.Join(elems...)
}
