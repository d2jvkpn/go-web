package aliyun

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/d2jvkpn/go-web/pkg/misc"
)

// default config: wk_config/test.yaml
func TestMain(m *testing.M) {
	var (
		configFile, testYaml string
		ossField, stsField   string
		err                  error
	)

	if testYaml, err = misc.RootFile("configs", "aliyun_test.yaml"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	testFlag = flag.NewFlagSet("testFlag", flag.ExitOnError)
	flag.Parse() // must do

	testFlag.StringVar(&configFile, "config", testYaml, "config filepath")
	testFlag.StringVar(&ossField, "oss", "aliyun_oss", "aliyun oss field in config")
	testFlag.StringVar(&stsField, "sts", "aliyun_sts", "aliyun sts field in config")

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

func TestConfig(t *testing.T) {
	var (
		config Config
		err    error
	)

	if config, err = NewConfig("config_demo.toml", "aliyun_oss"); err != nil {
		t.Fatal(err)
	}

	if config, err = NewConfig("config_demo.toml", "aliyun_sts"); err != nil {
		t.Fatal(err)
	}

	if config, err = NewConfig("config.demo.yaml", "aliyun_oss"); err != nil {
		t.Fatal(err)
	}

	if config, err = NewConfig("config.demo.yaml", "aliyun_sts"); err != nil {
		t.Fatal(err)
	}

	fmt.Println(config)
}

func TestConfigDemo(t *testing.T) {
	fmt.Printf(">>> TestConfigDemo:\n%s\n", ConfigDemo())
}
