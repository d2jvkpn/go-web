package aliyun

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var (
		configFile         string
		ossField, stsField string
		err                error
	)

	testFlag = flag.NewFlagSet("testFlag", flag.ExitOnError)
	flag.Parse() // must do

	testFlag.StringVar(&configFile, "config", "wk_01/test.yaml", "config filepath")
	testFlag.StringVar(&ossField, "ossField", "aliyun_oss", "aliyun oss field in config")
	testFlag.StringVar(&stsField, "stsField", "aliyun_sts", "aliyun sts field in config")

	testFlag.Parse(flag.Args())

	//	if testConfig, err = NewConfig(configFile, field); err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}

	if testOssClient, err = NewOssClient(configFile, ossField); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if testStsClient, err = NewStsClient(configFile, stsField); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m.Run()
}

func TestConfigDemo(t *testing.T) {
	fmt.Printf(">>> TestConfigDemo:\n%s\n", ConfigDemo())
}
