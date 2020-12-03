package terminal

import (
	"github.com/buger/goterm"
)

type Box struct {
	box *goterm.Box
}

func NewBox(widthPct int, heightPct int) *Box {
	return &Box{
		box: goterm.NewBox(widthPct|goterm.PCT, heightPct|goterm.PCT, 0),
	}

}

func (b *Box) Write(data []byte) (int, error) {
	return b.box.Write(data)
}
