package time

import(
		"testing"
)


func TestFormat(t *testing.T){
	tim,err := ParseFromYYYY_MM_DD_HH_MM_SS(format_yyyy_mm_dd_hh_mm_ss)
	if err != nil{
		t.Fatal(err)
	}
	timStr:= FormatToHH_MM_SS(tim)
	if  timStr!= format_hh_mm_ss{
		t.Fatalf("time format error, expect result %s, but get %s",format_hh_mm_ss,timStr)
	}
}
