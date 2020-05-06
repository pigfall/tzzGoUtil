package net

import (
	"testing"

	lys "github.com/Peanuttown/gopacket/layers"
	"github.com/Peanuttown/tzzGoUtil/log"
)

func TestListent(t *testing.T) {
	packets, closeF, err := Listen("eth0")
	if err != nil {
		t.Fatal(err)
	}
	var count uint
	for p := range packets {
		var layers = p.Layers()
		//t.Logf("layers : %d\n", len(layers))
		for i, l := range layers {
			if l.LayerType() == lys.LayerTypeICMPv4 {
				lys.ICMPv4
				t.Logf("cout:%d,%d layer is :%s,contentLength:%d, content:%s\n", count, i, l.LayerType().String(), len(l.LayerContents()), l.LayerPayload())
				count++
			}
		}
	}
	closeF()
	log.Debug("wait over")
}
