package kafka

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)

// $ go test -run TestHandler -- -index=10
func TestMain(m *testing.M) {
	var addrs string

	testFlag = flag.NewFlagSet("testFlag", flag.ExitOnError)
	flag.Parse() // must do

	testFlag.StringVar(&addrs, "addrs", "127.0.0.1:9093", "kakfa brokers address seperated by comma")
	testFlag.StringVar(&testTopic, "topic", "test", "kafka topic")
	testFlag.StringVar(&testGroupId, "groupId", "default", "kakfa group id")
	testFlag.StringVar(&testKafkaVersion, "kafkaVersion", "3.2.0", "kakfa version")

	testFlag.IntVar(&testIndex, "index", 0, "first message index")
	testFlag.IntVar(&testNum, "num", 10, "number of messages")
	testFlag.Int64Var(&testOffset, "offset", 0, "offset number")

	testFlag.Parse(flag.Args())

	for _, v := range strings.Split(addrs, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			testAddrs = append(testAddrs, v)
		}
	}

	fmt.Printf(
		"==> TestMain: testAddrs=%v, testTopic=%q, testGroupId=%q, testKafkaVersion=%q\n"+
			"    testIndex=%d, testNum=%d, testOffset=%d\n",
		testAddrs, testTopic, testGroupId, testKafkaVersion, testIndex, testNum, testOffset,
	)
	if testNum == 0 {
		fmt.Println("invalid num:", testNum)
		os.Exit(1)
	}

	m.Run()
}
