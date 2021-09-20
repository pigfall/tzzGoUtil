package bytes

import(
		"testing"
)


func TestHowManyBitIsOneInByte(t *testing.T){
	b := byte(1)
	res := HowManyBitIsOneInByte(b)
	if res != 1{
		t.Fatal("unexpted")
	}

	b = byte(2)
	res = HowManyBitIsOneInByte(b)
	if res != 1{
		t.Fatal("unexpted")
	}

	b = byte(3)
	res = HowManyBitIsOneInByte(b)
	if res != 2{
		t.Fatal("unexpted")
	}
	
	b= byte(255)
	res = HowManyBitIsOneInByte(b)
	if res != 8{
		t.Fatal("unexpted")
	}
}


