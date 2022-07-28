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
	testConfig    *Config
	testOssClient *OssClient
	testStsClient *StsClient
	testFlag      *flag.FlagSet

	//go:embed config.demo.md
	configDemo string
)

func init() {
}
