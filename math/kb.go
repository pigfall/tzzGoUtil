package math

import (
	"fmt"
	//stdmath "math"
)

type ByteUnit uint64

func ToByteUnit(numUnitByte uint64) ByteUnit {
	return ByteUnit(numUnitByte)
}

func (this ByteUnit) ToString() string {
	gb := float64(this) / (1024 * 1024 * 1024)
	if gb >= 1 {
		return fmt.Sprintf("%0.2f GB", gb)
	}
	mb := float64(this) / (1024 * 1024)
	if mb >= 1 {
		return fmt.Sprintf("%0.2f MB", mb)
	}
	kb := float64(this) / (1024)
	if kb >= 1 {
		return fmt.Sprintf("%0.2f KB", kb)
	}
	return fmt.Sprintf("%d Bytes", this)
}
