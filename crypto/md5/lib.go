package md5

import (
	stdmd5 "crypto/md5"
	"fmt"
)

func Hash(toBeHashed []byte) string {
	return fmt.Sprintf("%x", stdmd5.Sum(toBeHashed))
}
