package kafka

import(
	kf "github.com/Shopify/sarama"
)

type Option func(cfg *kf.Config)

func OptionSASLPlainAuth(userName string,password string)Option{
	return func(cfg *kf.Config){
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.Mechanism =  kf.SASLTypePlaintext 
		cfg.Net.SASL.User = userName
		cfg.Net.SASL.Password= password
	}
}
