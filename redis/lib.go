package redis

import(
	gredis "github.com/go-redis/redis/v8"
	"errors"
	"fmt"
	"context"
)

type Client struct{
	*gredis.Client
}

func NewClient(
	ctx context.Context,
	addr string,
	password string,
	db int,
	userName string,
)(*Client,error){
	options := gredis.Options{
		Addr:addr,
		Password:password,
		DB:db,
		Username:userName,
	}
	cli:= gredis.NewClient(
		&options,
	)
	err := cli.Ping(ctx).Err()
	if err != nil{
		return nil,err
	}
	return &Client{Client:cli},nil
}


// delete key ,if key not exist err is nil
func (this *Client) Del(ctx context.Context, keys ...string)error{
	return this.Client.Del(ctx,keys...).Err()
}

func (this *Client) SetKeyWithNoExpire(ctx context.Context,key string, value interface{})error{
	return this.Client.Set(ctx,key,value,0).Err()
}

func (this *Client) KeyExists(ctx context.Context,key string)(bool,error){
	existCount,err :=this.Client.Exists(ctx,key).Result()
	if err != nil{
		return false,err
	}
	if existCount >1 {
		return false,fmt.Errorf("unexpted existCount number %d",existCount)
	}

	return existCount == 1,nil
}

func ErrRedisReturnNil(err error)bool{
	return errors.Is(err,gredis.Nil)
}
