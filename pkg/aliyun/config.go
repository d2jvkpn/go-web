package aliyun

import (
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

type Config struct {
	AccessKeyId     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	RegionId        string `mapstructure:"region_id"`
	Bucket          string `mapstructure:"bucket"`

	RoleArn        string `mapstructure:"role_arn"`        // sts
	ExpiredSeconds int    `mapstructure:"expired_seconds"` // sts
}

func NewConfig(fp, field string) (config *Config, err error) {
	var conf *viper.Viper

	conf = viper.New()
	conf.SetConfigName("aliyun_config")

	switch {
	case strings.HasSuffix(fp, ".toml"):
		conf.SetConfigType("toml")
	case strings.HasSuffix(fp, ".yaml"):
		conf.SetConfigType("yaml")
	default:
		return nil, fmt.Errorf("unkonw config file, use .yaml or .toml")
	}
	conf.SetConfigFile(fp)

	if err = conf.ReadInConfig(); err != nil {
		return nil, err
	}

	config = new(Config)
	if err = conf.UnmarshalKey(field, config); err != nil {
		return nil, err
	}

	if err = config.Valid(); err != nil {
		return nil, err
	}
	return config, nil
}

func ConfigDemo() string {
	return configDemo
}

func (config *Config) Valid() (err error) {
	if config.AccessKeyId == "" || config.AccessKeySecret == "" {
		return fmt.Errorf("access_key_id or access_key_secret is empty")
	}

	if config.RegionId == "" || config.Bucket == "" {
		return fmt.Errorf("regionId or bucket is empty")
	}

	return nil
}

func NewOssClient(fp, field string) (client *OssClient, err error) {
	var config *Config

	if config, err = NewConfig(fp, field); err != nil {
		return nil, err
	}

	client = &OssClient{
		config: *config,
	}

	client.Client, err = oss.New(
		fmt.Sprintf("https://oss-%s.aliyuncs.com", config.RegionId),
		config.AccessKeyId, config.AccessKeySecret,
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewStsClient(fp, field string) (client *StsClient, err error) {
	var config *Config

	if config, err = NewConfig(fp, field); err != nil {
		return nil, err
	}

	client = &StsClient{config: *config}

	if config.RoleArn == "" {
		return nil, fmt.Errorf("role_arn is unset")
	}

	client.Client, err = sts.NewClientWithAccessKey(
		config.RegionId,
		config.AccessKeyId,
		config.AccessKeySecret,
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (config *Config) Url(ps ...string) (link string) {
	if len(ps) == 0 {
		return fmt.Sprintf("https://%s.oss-%s.%s", config.Bucket, config.RegionId, ALIYUN_Domain)
	}

	p := strings.TrimLeft(ps[0], "/")
	return fmt.Sprintf("https://%s.oss-%s.%s/%s", config.Bucket, config.RegionId, ALIYUN_Domain, p)
}

func ValidSubpath(p string) (out string, err error) {
	if out = strings.Trim(p, "/"); out == "" {
		return "", fmt.Errorf("invalid subpath")
	}

	return out, nil
}
