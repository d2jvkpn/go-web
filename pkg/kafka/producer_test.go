package kafka

import (
	"fmt"
	"log"
	"testing"

	"github.com/Shopify/sarama"
)

func TestProducer(t *testing.T) {
	var (
		err      error
		config   *sarama.Config
		producer sarama.AsyncProducer
	)

	config = sarama.NewConfig()
	if config.Version, err = sarama.ParseKafkaVersion(_KafkaVersion); err != nil {
		t.Fatal(err)
	}

	if producer, err = sarama.NewAsyncProducer(_Addrs, config); err != nil {
		t.Fatal(err)
	}

	for i := _Index; i < _Index+_Num; i++ {
		msg := fmt.Sprintf("hello message: %d", i)
		log.Println("--> send msg:", msg)

		pmsg := &sarama.ProducerMessage{
			Topic: _Topic,
			Key:   sarama.StringEncoder("e0001"),
			Value: sarama.ByteEncoder([]byte(msg)),
		}

		producer.Input() <- pmsg
	}

	if err = producer.Close(); err != nil {
		t.Fatal(err)
	}
}
