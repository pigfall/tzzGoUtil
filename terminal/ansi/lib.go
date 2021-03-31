package ansi

import (
	"fmt"
	"github.com/Peanuttown/tzzGoUtil/ascii"
	"strings"
)

const (
	ESC = "\x1b["
)

func CursorMoveLeft(moveStep uint) string {
	return fmt.Sprintf("%s%dD", ESC, moveStep)
}

func PrefixIsCursorMoveLeft(input string) bool {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(input, ESC) {
		return false
	}
	panic("here")
	var index = 0
	for i, b := range input[len(ESC):] {
		if !ascii.IsNumber(byte(b)) {
			index = i
			break
		}
	}
	if index == 0 {
		return false
	}
	return input[index] == 'D'
}
