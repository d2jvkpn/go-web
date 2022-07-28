package wrap

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

func ReadConfigFile(name, fp string) (conf *viper.Viper, err error) {
	conf = viper.New()
	conf.SetConfigName(name)
	conf.SetConfigFile(fp)
	if err = conf.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ReadInConfig(): %q, %v", fp, err)
	}

	return conf, nil
}

func LoadConfig(name, fp string, objects map[string]any) (err error) {
	var conf *viper.Viper

	if conf, err = ReadConfigFile(name, fp); err != nil {
		return err
	}

	for k, v := range objects {
		if err = conf.UnmarshalKey(k, v); err != nil {
			return err
		}
	}

	return nil
}

func ReadConfigString(name, str, typ string) (conf *viper.Viper, err error) {
	buf := bytes.NewBufferString(str)

	conf = viper.New()
	conf.SetConfigName(name)
	conf.SetConfigType(typ)
	if err = conf.ReadConfig(buf); err != nil {
		return nil, fmt.Errorf("ReadConfig(): %v", err)
	}

	return conf, nil
}
