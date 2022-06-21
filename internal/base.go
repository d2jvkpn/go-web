package internal

import (
	// "fmt"
	"embed"
	"net/http"

	"github.com/d2jvkpn/go-web/pkg/misc"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	//go:embed static
	_Static embed.FS
	//go:embed templates
	_Templates embed.FS

	_Release bool
	// _InstanceId string
	_BuildInfo [][2]string

	Cron_At   = "0 0 1 * * *" // everyday at 1 o'clock
	Cron_Name = "Cron Test"

	_Cron      *cron.Cron
	_Config    *viper.Viper
	_Server    *http.Server
	_ApiLogger *misc.Logger
)

func init() {
	// _InstanceId = misc.RandString(16)
	_Cron = cron.New(cron.WithSeconds())
}

func LoadBuildInfo(info [][2]string) {
	_BuildInfo = info
}

func logBuildInfo(logger *zap.Logger) {
	fields := make([]zap.Field, 0, len(_BuildInfo))
	for _, v := range _BuildInfo {
		fields = append(fields, zap.String(v[0], v[1]))
	}

	logger.Info("BuildInfo", zap.Any(
		"event",
		map[string]interface{}{"kind": "info", "fields": fields},
	))
}
