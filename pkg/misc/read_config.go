package misc

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
