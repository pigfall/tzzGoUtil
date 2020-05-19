package mongo

import (
	"testing"
)

func TestConnect(t *testing.T) {
	uri := MongoURI{}
	(&uri).Init("127.0.0.1", 27017)
	_, err := Connect(&uri)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMgInsert(t *testing.T) {

}
