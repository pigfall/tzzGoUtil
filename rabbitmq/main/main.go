package main

import (
	"time"
	"github.com/streadway/amqp"
	rb "github.com/Peanuttown/tzzGoUtil/rabbitmq"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// conn, err := amqp.Dial("amqp://MjphbXFwLWNuLTJyNDI4YmlpYzAwODpMVEFJNXQ3WHZkYWoxUFA4aE1RVFRqemI=:MjlBODdERDIzMURBODAwNUEyMzQyQkQzOThGMEJGOUE5M0Y1QTdBRDoxNjI0NDI3MjE3MTQ1@amqp-cn-2r428biic008.mq-amqp.cn-shanghai-867405-a.aliyuncs.com:5672/sweetlucky")
	cfg := rb.NewConnCfg(
		true,
		"MjphbXFwLWNuLTJyNDI4YmlpYzAwODpMVEFJNXQ3WHZkYWoxUFA4aE1RVFRqemI=",
		"MjlBODdERDIzMURBODAwNUEyMzQyQkQzOThGMEJGOUE5M0Y1QTdBRDoxNjI0NDI3MjE3MTQ1",
		"amqp-cn-2r428biic008.mq-amqp.cn-shanghai-867405-a.aliyuncs.com",
		"5672",
	)
	cli,err := rb.NewClient(cfg,rb.NewClientOptionVHost("sweetlucky"))
	if err != nil{
		panic(err)
	}
	failOnError(err, "Failed to connect to RabbitMQ")
	defer cli.Close()
	ch, err :=cli.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	err = ch.ExchangeDeclare(
		"testExchange",
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil{
		panic(err)
	}
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	err = ch.QueueBind(q.Name,"","testExchange",false,nil)
	if err != nil{
		panic(err)
	}

	body := "Hello World!"
	var confirms = make(chan amqp.Confirmation,100)
	go func(){
		err := ch.Confirm(false)
		if err != nil{
			panic(err)
		}
		confirms= ch.NotifyPublish(confirms)
		if err != nil{
			panic(err)
		}
		for cfm := range confirms{
			log.Println(cfm)
		}
	}()
	for{
		err = ch.Publish(
			"testExchange",     // exchange
			//q.Name, // routing key
			"",
			false,  // mandatory
			false,  // immediate
			amqp.Publishing {
				ContentType: "text/plain",
				Body:        []byte(body),
			})
			if err != nil{
				log.Println("publish err ",err)
			}else{
				log.Println("publis success")
			}

		time.Sleep(time.Second*2)
	}
}
