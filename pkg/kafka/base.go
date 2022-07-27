package kafka

import (
	"flag"
	"log"

	"github.com/Shopify/sarama"
)

const (
	RFC3339ms = "2006-01-02T15:04:05.000Z07:00"
)

var (
	testKafkaVersion string
	testAddrs        []string
	testTopic        string
	testGroupId      string

	testIndex, testNum int
	testOffset         int64

	testFlag *flag.FlagSet
)

func init() {
	// testAddrs = []string{"127.0.0.1:9093"}
	// testTopic = "test"
	// testGroupId = "default"
	// testKafkaVersion = "3.2.0"
	// testIndex = 0
	// testNum = 10
	// testOffset = 0
}

func defaultProcess(msg *sarama.ConsumerMessage) (metadata string, err error) {
	tmpl := "<-- msg.Timestamp=%q, msg.Topic=%q, msg.Partition=%d, msg.Offset=%v,\n" +
		"    key=%q, value=%q\n"

	// msg.BlockTimestamp
	log.Printf(
		tmpl, msg.Timestamp.Format(RFC3339ms), msg.Topic, msg.Partition, msg.Offset,
		msg.Key, msg.Value,
	)

	return "consumed", nil
}
