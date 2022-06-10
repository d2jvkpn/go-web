package internal

import (
	// "fmt"
	"embed"
	"net/http"

	"github.com/d2jvkpn/goapp/pkg/misc"

	"github.com/spf13/viper"
)

var (
	//go:embed static
	_Static embed.FS
	//go:embed templates
	_Templates embed.FS

	_Release    bool
	_InstanceId string
	_Config     *viper.Viper
	_Server     *http.Server
	_ApiLogger  *misc.Logger
	BuildInfo   [][2]string
)

func init() {
	_InstanceId = misc.RandString(16)
}
