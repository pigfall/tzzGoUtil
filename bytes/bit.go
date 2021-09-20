package bytes


func HowManyBitIsOneInByte(bte byte)(int){
	var count int
	for i:=7;i>=0;i--{
		if ((bte >> i) & (1)) > 0{
			count++
		}
	}
	return count
}
