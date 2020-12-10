package ssh

import (
	"testing"
)

func TestSession(t *testing.T) {
	clt, err := Dial("172.16.2.83:22", &DialCfg{User: "root", Passwd: "123456"})
	if err != nil {
		t.Fatal(err)
	}
	defer clt.Close()
	s, err := clt.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	err = s.Run("ls")
	if err != nil {
		t.Fatal(err)
	}
	err = s.Run("ls")
	if err != nil {
		t.Fatal(err)
	}
}
