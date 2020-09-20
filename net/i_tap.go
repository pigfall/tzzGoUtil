package net

type ITunTap interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close()
}

type ITap interface {
	ITunTap
}

type ITun interface {
	ITunTap
}
