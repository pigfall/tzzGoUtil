package ascii

import(
    "testing"
)

func RuneTest(b byte){

}

func TestIsAlha(t *testing.T){
    var s = "abc数"
    for i:=0;i<len(s);i++{
        if i<=2{
            if !IsAlpha(s[i]){
                t.Fatalf("expect isAlpha:%v\n",s[i])
            }
        }else{
            if IsAlpha(s[i]){
                t.Fatalf("expect no alpha:%v\n",s[i])
            }
        }
    }
}

func TestToLower(t *testing.T){
    var upper = byte('A')
    if ToLower(upper)!='a'{
        t.Fatal("to lower failed")
    }

}
func TestToUpper(t *testing.T){
    var lower = byte('a')
    if ToUpper(lower)!='A'{
        t.Fatal("to upper failed")
    }

}
