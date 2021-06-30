package redis

import(
	"context"
	"testing"
)


func testCli()*Client{
	ctx := context.Background()
	cli,err := NewClient(ctx,"localhost:6379","",0,"")
	if err !=nil{
		panic(err)
	}
	return cli
}

func TestRedisDel(t *testing.T){
	ctx := context.Background()
	cli,err := NewClient(ctx,"localhost:6379","",0,"")
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

func TestRedisLua(t *testing.T){
	ctx := context.Background()
	cli,err := NewClient(ctx,"localhost:6379","",0,"")
	if err !=nil{
		t.Fatal(err)
	}
	res,err := cli.Eval(
		ctx,
		`
		redis.call("set","testKey","value")
		local res =redis.call("incr","testKey",1)
		return res+1
		`,
		nil,
	).Result()
	t.Log(err)
	t.Log(res)
}

func TestRedisLuaNilSet(t *testing.T){
	cli := testCli()
	ctx := context.Background()
	t.Log(
			
cli.Eval(
		ctx,
		`
		redis.call("del","noExistKey")
		local value = redis.call("get","noExistKey")
		if not value  then
			return -1
		end
		return 1
		`,
		nil,
	).Result(),
	)
}
