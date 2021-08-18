package kafka

import(
	kf "github.com/Shopify/sarama"
)


type ConsumerMsg <-chan *kf.ConsumerMessage


