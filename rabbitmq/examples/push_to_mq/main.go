package main

//import (
//	"fmt"	
//	rq "github.com/pigfall/tzzGoUtil/rabbitmq"
//	"flag"
//	"log"
//)
//
//func main() {
//	// {
//	var rqUser string
//	var rqPassword string
//	var rqServerHost string
//	var rqServerPort string
//	var vhost string
//	var exchange string
//	flag.StringVar(&rqUser,"user","guest","rabbitmq user")
//	flag.StringVar(&rqPassword,"password","demo_password","rabbitmq password")
//	flag.StringVar(&rqServerHost,"server_host","localhost","rabbitmq server host")
//	flag.StringVar(&rqServerPort,"server_port","5672","rabbitmq server port")
//	flag.StringVar(&vhost,"vhost","","vhost")
//	flag.StringVar(&exchange,"exchange","","exchange")
//	flag.Parse()
//	if len(exchange) == 0{
//		log.Fatal("exchange is nil")
//	}
//	// }
//	clientConnCfg := rq.NewConnCfg(
//			 true,
//			 rqUser,
//			 rqPassword,
//			 rqServerHost,
//			 rqServerPort,
//	)
//
//	opts := make([]rq.NewClientOption,0,)
//	if len(vhost) != 0{
//		opts = append(opts,rq.NewClientOptionVHost(vhost))
//	}
//
//	cli,err := rq.NewClient(
//		clientConnCfg,
//		opts...,
//	)
//	if err != nil{
//		log.Fatal(err)
//	}
//
//	ch,err := cli.Channel()
//	if err != nil{
//		log.Fatal(err)
//	}
//}
