package output

import (
	"fmt"
	"os"
)

func Err(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}
func Errf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}
