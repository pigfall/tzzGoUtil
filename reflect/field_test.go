package reflect

import(
		"testing"
)

func TestGetEmptyField(t *testing.T){
	type T struct{
		F1 string
		F2 string
	}
	tt := T{}
	t.Log(GetStructEmptyField(tt))
}
