package net

import(
)

type TunIfce interface{
	Write(p []byte) (n int, err error)
	Read(p []byte) (n int, err error)
	Close() error
	Name() string
	SetIp(ip ...string)(err error)
}

