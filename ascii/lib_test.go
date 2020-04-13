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
