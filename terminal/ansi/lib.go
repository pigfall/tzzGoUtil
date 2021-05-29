package ansi

import (
	"fmt"
//	"github.com/Peanuttown/tzzGoUtil/ascii"
//	"strings"
)

const (
	ESC = "\x1b["
)

func CursorMoveLeft(moveStep uint) string {
	return fmt.Sprintf("%s%dD", ESC, moveStep)
}

