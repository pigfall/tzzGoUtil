package redis

import(
	"context"
	"testing"
)


func TestRedisDel(t *testing.T){
	ctx := context.Background()
	cli,err := NewClient(ctx,"localhost:6379","",0)
	if err !=nil{
		t.Fatal(err)
	}
	const testKey = "testKey"
	cli.Del(ctx,testKey)
	// key not exist , exptec err is nil
	err = cli.Del(ctx,testKey)
	if err != nil{
		t.Fatal(err)
	}
	err = cli.SetKeyWithNoExpire(ctx,testKey,"testValue")
	if err != nil{
		t.Fatal(err)
	}
	err = cli.Del(ctx,testKey)
	if err != nil{
		t.Fatal()
	}
	existCount,err := cli.Exists(ctx,testKey).Result()
	if err !=nil{
		t.Fatal(err)
	}
	if existCount != 0 {
		t.Fatal("delted failed")
	}
}
