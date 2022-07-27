package kafka

import (
	// "fmt"
	"context"
	"flag"
	"log"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

func TestConsumer(t *testing.T) {
	var (
		topics     []string
		partitions []int32
		err        error
		cancel     chan struct{}

		config    *sarama.Config
		consumer  sarama.Consumer
		pconsumer sarama.PartitionConsumer
	)

	flag.Parse()
	config = sarama.NewConfig()
	if config.Version, err = sarama.ParseKafkaVersion(_KafkaVersion); err != nil {
		t.Fatal(err)
	}

	cancel = make(chan struct{})

	if consumer, err = sarama.NewConsumer(_Addrs, config); err != nil {
		t.Fatal(err)
	}

	if topics, err = consumer.Topics(); err != nil {
		t.Fatal(err)
	}

	if len(topics) == 0 {
		t.Fatal("no topics")
	}

	if partitions, err = consumer.Partitions(topics[0]); err != nil {
		t.Fatal(err)
	}
	if len(partitions) == 0 {
		t.Fatalf("topics %s has no partitions\n", topics[0])
	}
	log.Printf("~~~ partitions of %s: %v\n", topics[0], partitions)

	// topic string, partition int32, offset int64
	if pconsumer, err = consumer.ConsumePartition(_Topic, 0, 0); err != nil {
		t.Fatal(err)
	}

	go TestProducer(t)

	go func() {
		mc := pconsumer.Messages() // *sarama.ConsumerMessage

		tmpl := "<-- msg.Timestamp=%+v, msg.Topic=%v, msg.Partition=%v, msg.Offset=%v\n" +
			"    key: %q, value: %q\n"

		for {
			select {
			case msg := <-mc:
				log.Printf(
					tmpl, msg.Timestamp, msg.Topic, msg.Partition, msg.Offset,
					msg.Key, msg.Value,
				)
			case <-cancel:
				return
			}
		}
	}()

	time.Sleep(15 * time.Second)
	close(cancel)

	if err = consumer.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestHandler(t *testing.T) {
	var (
		err error
		ctx context.Context

		config  *sarama.Config
		group   sarama.ConsumerGroup
		handler *Handler // sarama.ConsumerGroupHandler
	)

	config = sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	if config.Version, err = sarama.ParseKafkaVersion(_KafkaVersion); err != nil {
		t.Fatal(err)
	}
	group, err = sarama.NewConsumerGroup(_Addrs, _GroupId, config)
	if err != nil {
		t.Fatal(err)
	}

	go TestProducer(t)

	ctx = context.Background()
	handler = NewHandler(ctx, group, defaultProcess)
	handler.Consume(_Topic)

	time.Sleep(15 * time.Second)
	log.Println("<<< Exit")

	handler.Close()
	if err = group.Close(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
}
