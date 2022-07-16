package settings

import (
	"github.com/d2jvkpn/go-web/pkg/wrap"

	"github.com/spf13/viper"
)

var (
	Config *viper.Viper
	// other global workers...
	Logger *wrap.Logger
)
