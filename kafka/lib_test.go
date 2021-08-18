package kafka

import(
	"strings"
	"os"
	"fmt"
	"testing"
)

func testKafkaCfg()(brokerAddrs []string,userName,password string){
	errWrapper := func(key string)error{
		return fmt.Errorf("not found %s in env variable",key)
	}
	const keyTestKafkaAddrs = "TEST_KAFKA_ADDR"
	const keyTestKafkaUser = "TEST_KAFKA_USER"
	const keyTestKafkaPassword = "TEST_KAFKA_PASSWORD"
	addrs, ok := os.LookupEnv(keyTestKafkaAddrs)
	if !ok {
		panic(errWrapper(keyTestKafkaAddrs))
	}
	userName, ok = os.LookupEnv(keyTestKafkaUser)
	if !ok {
		panic(errWrapper(keyTestKafkaUser))
	}
	password, ok = os.LookupEnv(keyTestKafkaPassword)
	if !ok {
		panic(errWrapper(keyTestKafkaPassword))
	}

	return strings.Split(addrs,","), userName,password
}

func TestReader(t *testing.T){
	brokerAddrs,userName,password:= testKafkaCfg()
	c,err := NewConsumer(brokerAddrs,0,"example",0,OptionSASLPlainAuth(userName,password))
	if err != nil{
		t.Fatal(err)
	}
	for{
			msg := <-c.Messages()
			t.Log(msg)
	}
}
