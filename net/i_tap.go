package net

type ITap interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}

type ITun interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
}
