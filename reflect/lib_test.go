package reflect

import(
    "testing"
)

type testStruct struct{
    Bool bool
    Age int32
}

func TestSetValue(t *testing.T){
    s := testStruct{Bool:false}
    SetValue(&s,"Bool",true,)
    SetValue(&s,"Age",int32(32))
    t.Log(s.Bool)
    t.Log(s.Age)
}
