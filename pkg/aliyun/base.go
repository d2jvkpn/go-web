package aliyun

import (
	_ "embed"
	"flag"
)

const (
	ALIYUN_Domain         = "aliyuncs.com"
	ALIYUN_Code_NoSuchKey = "NoSuchKey"
)

var (
	//go:embed config_demo.md
	configDemo string

	testConfig    *Config
	testFlag      *flag.FlagSet
	testOssClient *OssClient
	testStsClient *StsClient
)

func init() {
}
