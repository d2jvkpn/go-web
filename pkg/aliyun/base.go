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
	testOssClient *OssClient
	testStsClient *StsClient
	testFlag      *flag.FlagSet
)

func init() {
}
