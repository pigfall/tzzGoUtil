package net

type TunIfce interface{
	Write(p []byte) (n int, err error)
	Read(p []byte) (n int, err error)
	Close() error
	Name() string
}

