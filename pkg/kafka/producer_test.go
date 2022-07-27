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
	if config.Version, err = sarama.ParseKafkaVersion(testKafkaVersion); err != nil {
		t.Fatal(err)
	}

	if producer, err = sarama.NewAsyncProducer(testAddrs, config); err != nil {
		t.Fatal(err)
	}

	for i := testIndex; i < testIndex+testNum; i++ {
		msg := fmt.Sprintf("hello message: %d", i)
		log.Println("--> send msg:", msg)

		pmsg := &sarama.ProducerMessage{
			Topic: testTopic,
			Key:   sarama.StringEncoder("key0001"),
			Value: sarama.ByteEncoder([]byte(msg)),
		}

		producer.Input() <- pmsg
	}

	if err = producer.Close(); err != nil {
		t.Fatal(err)
	}
}
