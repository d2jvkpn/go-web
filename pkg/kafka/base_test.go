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
	var (
		addrs    string
		mainFlag *flag.FlagSet
	)

	mainFlag = flag.NewFlagSet("mainFlag", flag.ExitOnError)
	flag.Parse() // must do

	mainFlag.StringVar(&addrs, "addrs", "127.0.0.1:9093", "kakfa brokers address seperated by comma")
	mainFlag.StringVar(&testTopic, "topic", "test", "kafka topic")
	mainFlag.StringVar(&testGroupId, "groupId", "default", "kakfa group id")
	mainFlag.StringVar(&testKafkaVersion, "kafkaVersion", "3.2.0", "kakfa version")

	mainFlag.IntVar(&testIndex, "index", 0, "first message index")
	mainFlag.IntVar(&testNum, "num", 10, "number of messages")
	mainFlag.Int64Var(&testOffset, "offset", 0, "offset number")

	mainFlag.Parse(flag.Args())

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
