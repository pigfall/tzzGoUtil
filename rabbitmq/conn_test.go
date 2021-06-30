package rb

import(
	"time"
	"log"
	"testing"
	"fmt"

		"github.com/streadway/amqp"
)

func TestNewClient(t *testing.T){
	cfg := NewConnCfg(
		true,
		"MjphbXFwLWNuLTJyNDI4YmlpYzAwODpMVEFJNXQ3WHZkYWoxUFA4aE1RVFRqemI=",
		"MjlBODdERDIzMURBODAwNUEyMzQyQkQzOThGMEJGOUE5M0Y1QTdBRDoxNjI0NDI3MjE3MTQ1",
		"amqp-cn-2r428biic008.mq-amqp.cn-shanghai-867405-a.aliyuncs.com",
		"5672",
	)
	t.Log(cfg.ToUrl())
//NewClientOptionVHost("sweetlucky")
	_,err := NewClient(cfg,NewClientOptionVHost("sweetlucky"))
	if err != nil{
		t.Fatal(err)
	}
}

func TestPublishConfirm(t *testing.T){
	log.SetFlags(log.LstdFlags|log.Lshortfile)
	cfg := NewConnCfg(
		true,
		"MjphbXFwLWNuLTJyNDI4YmlpYzAwODpMVEFJNXQ3WHZkYWoxUFA4aE1RVFRqemI=",
		"MjlBODdERDIzMURBODAwNUEyMzQyQkQzOThGMEJGOUE5M0Y1QTdBRDoxNjI0NDI3MjE3MTQ1",
		"amqp-cn-2r428biic008.mq-amqp.cn-shanghai-867405-a.aliyuncs.com",
		"5672",
	)
	t.Log(cfg.ToUrl())
//NewClientOptionVHost("sweetlucky")
	go func(){
		cli,err := NewClient(cfg,NewClientOptionVHost("sweetlucky"))
		if err != nil{
			panic(err)
		}
		ch,err := cli.Channel()
		if err != nil{
			panic(err)
		}
		for{
			time.Sleep(time.Second*3)
			err := ch.Publish(
				"testExchange",
				"",
				false,
				false,
				amqp.Publishing{
					Headers:map[string]interface{}{
						"test":"value",
					},
				},
			)
			if err != nil{
				log.Println("publish msg err ",err.Error())
				log.Println("connection is close ???",cli.IsClosed())
			}else{
				log.Println("publish msg success")
			}
		}
	}()
	go func(){
		cli,err := NewClient(cfg,NewClientOptionVHost("sweetlucky"))
		if err != nil{
			panic(err)
		}
		ch,err := cli.Channel()
		if err != nil{
			panic(err)
		}
		_,err = ch.QueueDeclare("testQueue",true,false,false,false,nil)
		if err != nil{
			panic(err)
		}
		err = ch.QueueBind("testQueue","","testExchange",false,nil)
		if err != nil{
			panic(err)
		}
		consumeCh,err := ch.Consume("testQueue","testConsumer",false,false,false,false,nil)
		if err != nil{
			panic(err)
		}
		log.Println("ready to consume msg")
		for msg := range consumeCh {
			log.Println("consume msg ",msg)
		}
	}()
		cli,err := NewClient(cfg,NewClientOptionVHost("sweetlucky"))
		if err != nil{
			t.Fatal(err)
		}
	fmt.Println("loop to confirm")
	ch,err := cli.Channel()
	if err != nil{
		t.Fatal(err)
	}
	err = ch.Confirm(false)
	if err != nil{
		t.Fatal(err)
	}
	err = ch.ExchangeDeclare(
		"testExchange","fanout",true,false,false,false,nil,
	)
	if err != nil{
		t.Log(err)
	}
	publicConfirm := make(chan amqp.Confirmation,100)
	publicConfirm = ch.NotifyPublish(publicConfirm)
	for{
		select{
		case confirm,ok:=<-publicConfirm:
			fmt.Println(len(publicConfirm))
			if ok{
				fmt.Printf("ok: %+v\n",confirm)
			}else{
				fmt.Printf("not ok, %+v\n",confirm)
			}
			err = ch.Nack(confirm.DeliveryTag,false,false)
			if err != nil{
				log.Println(err)
			}
			//err = ch.Ack(confirm.DeliveryTag,false)
			//if err != nil{
			//	log.Println(err)
			//}
			time.Sleep(time.Second)
		}
	}

}
