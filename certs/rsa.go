package certs

import (
	"crypto/rand"
	"crypto/rsa"
)

/*
RSAGenPrivateKey: 生成 rsa private key
*/

type PrivateKeyBitSize int

const (
	PrivateKeyBitSize_1024 PrivateKeyBitSize = 1024
	PrivateKeyBitSize_2048 PrivateKeyBitSize = 2048
)

func RSAGenPrivateKey(bitSize PrivateKeyBitSize) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, int(bitSize))
}
