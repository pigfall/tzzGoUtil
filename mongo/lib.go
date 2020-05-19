package mongo

import (
	"context"
	"fmt"
	"time"

	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoURI struct {
	Addr string
	Port uint32
}

func (this *MongoURI) Init(addr string, port uint32) {
	this.Addr = addr
	this.Port = port
}

func (this MongoURI) ToString() string {
	return fmt.Sprintf("mongodb://%s:%d", this.Addr, this.Port)
}

type Client struct {
	c *mg.Client
}

func newClient(c *mg.Client) *Client {
	return &Client{c: c}
}

func Connect(uri *MongoURI) (c *Client, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var client *mg.Client
    client, err = mg.Connect(ctx, options.Client().ApplyURI(uri.ToString()))
	if err != nil {
		return
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return
	}
	c = newClient(client)
	return
}
