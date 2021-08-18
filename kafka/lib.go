package kafka

import(
	"fmt"
	//"github.com/confluentinc/confluent-kafka-go/kafka"
//	kf "github.com/confluentinc/confluent-kafka-go/kafka"
//	"context"
	kf "github.com/Shopify/sarama"
)




type ConsumerI interface{
	Close() error
	Messages() ConsumerMsg
}


type Consumer struct{
	c kf.Consumer
	partitionConsumer kf.PartitionConsumer
}

func NewConsumer(servers []string,offset int,topic string,partition int,opts ...Option)(ConsumerI,error){
	cfg := kf.NewConfig()
	for _,opt := range opts{
		opt(cfg)
	}
	fmt.Println(servers)
	fmt.Println(cfg.Net)
	c,err := kf.NewConsumer(servers,cfg)
	if err != nil{
		return nil,fmt.Errorf("Create consumer failed: %w",err)
	}
	partitionConsumer,err := c.ConsumePartition(topic,int32(0),int64(offset))
	if err != nil{
		c.Close()
		return nil, fmt.Errorf("Create partitionConsumer failed: %w",err)
	}
	return &Consumer{
		c:c,
		partitionConsumer:partitionConsumer,
	},nil
}


func (this *Consumer) Close() error{
	this.partitionConsumer.Close()
	return this.c.Close()
}

func (this *Consumer) Messages() ConsumerMsg {
	return this.partitionConsumer.Messages()
}

