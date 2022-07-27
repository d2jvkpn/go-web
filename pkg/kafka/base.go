package kafka

import (
	"log"

	"github.com/Shopify/sarama"
)

var (
	_KafkaVersion string
	_Addrs        []string
	_Topic        string
	_GroupId      string

	_Index, _Num int
)

func init() {
	// _Addrs = []string{"127.0.0.1:9093"}
	// _Topic = "test"
	// _GroupId = "default"
	// _KafkaVersion = "3.2.0"
	// _Index = 0
	// _Num = 10
}

func defaultProcess(msg *sarama.ConsumerMessage) (metadata string, err error) {
	tmpl := "<-- msg.Timestamp=%q, msg.Topic=%q, msg.Partition=%d, msg.Offset=%v\n" +
		"    key: %q, value: %q\n"

	log.Printf(
		tmpl, msg.Timestamp, msg.Topic, msg.Partition, msg.Offset,
		msg.Key, msg.Value,
	)

	return "consumed", nil
}
