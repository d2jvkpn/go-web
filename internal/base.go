package internal

import (
	// "fmt"
	"embed"
	"net/http"

	"github.com/d2jvkpn/go-web/pkg/misc"

	"github.com/robfig/cron/v3"
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

	Cron_At   = "0 0 1 * * *" // everyday at 1 o'clock
	Cron_Name = "Cron Test"
	_Cron     *cron.Cron
)

func init() {
	_InstanceId = misc.RandString(16)
	_Cron = cron.New(cron.WithSeconds())

	misc.RegisterDefaultLogFmt()
}
