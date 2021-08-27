package reflect

import(
	"testing"
)


func TestFieldStringIsNilOrLenZero(t *testing.T){
	type Demo struct{
		F1 string
		F2 *string
	}
	d := &Demo{}
	fields := FieldStringIsNilOrLenZero(d)
	if len(fields) != 2 {
		t.Fatalf("Unexptected, %v",fields)
	}
	if fields[0] != "F1" || fields[1] != "F2"{
		t.Fatalf("Unexptected, %v",fields)
	}
}
