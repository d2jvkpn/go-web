package kafka

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	var (
		addrs  string
		mainfs *flag.FlagSet
	)

	mainfs = flag.NewFlagSet("mainfs", flag.ExitOnError)
	flag.Parse() // must

	mainfs.StringVar(&addrs, "addrs", "127.0.0.1:9093", "kakfa brokers address seperated by comma")
	mainfs.IntVar(&_Index, "index", 0, "first message index")
	mainfs.IntVar(&_Num, "num", 10, "number of messages")
	mainfs.Parse(flag.Args())

	for _, v := range strings.Split(addrs, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			_Addrs = append(_Addrs, v)
		}
	}

	fmt.Printf("==> TestMain: _Addrs=%v, _Index=%d, _Num=%d\n", _Addrs, _Index, _Num)
	if _Num == 0 {
		fmt.Println("invalid num:", _Num)
		os.Exit(1)
	}

	m.Run()
}
