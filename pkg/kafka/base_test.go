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
		addrs  string
		mainfs *flag.FlagSet
	)

	mainfs = flag.NewFlagSet("mainfs", flag.ExitOnError)
	flag.Parse() // must do

	mainfs.StringVar(&addrs, "addrs", "127.0.0.1:9093", "kakfa brokers address seperated by comma")
	mainfs.StringVar(&testTopic, "topic", "test", "kafka topic")
	mainfs.StringVar(&testGroupId, "groupId", "default", "kakfa group id")
	mainfs.StringVar(&testKafkaVersion, "kafkaVersion", "3.2.0", "kakfa version")

	mainfs.IntVar(&testIndex, "index", 0, "first message index")
	mainfs.IntVar(&testNum, "num", 10, "number of messages")
	mainfs.Int64Var(&testOffset, "offset", 0, "offset number")

	mainfs.Parse(flag.Args())

	for _, v := range strings.Split(addrs, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			testAddrs = append(testAddrs, v)
		}
	}

	fmt.Printf(
		"==> TestMain: testAddrs=%v, testIndex=%d, testNum=%d\n",
		testAddrs, testIndex, testNum,
	)
	if testNum == 0 {
		fmt.Println("invalid num:", testNum)
		os.Exit(1)
	}

	m.Run()
}
