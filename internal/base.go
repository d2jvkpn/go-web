package internal

import (
	// "fmt"
	"embed"
	"expvar"
	"net/http"
	"runtime"
	"time"

	"github.com/d2jvkpn/go-web/pkg/misc"
	"github.com/d2jvkpn/go-web/pkg/wrap"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	HTTP_MaxHeaderBytes     = 2 << 20
	HTTP_ReadTimeout        = 10 * time.Second
	HTTP_ReadHeaderTimeout  = 5 * time.Second
	HTTP_WriteTimeout       = 10 * time.Second
	HTTP_IdleTimeout        = 60
	HTTP_MaxMultipartMemory = 8 << 20 // 8M
)

var (
	//go:embed static
	_Static embed.FS
	//go:embed templates
	_Templates  embed.FS
	_StaticDirs []wrap.StaticDir

	_Release  bool
	_Instance string
	// _InstanceId string
	_BuildInfo [][2]string

	Cron_At   = "0 0 1 * * *" // everyday at 1 o'clock
	Cron_Name = "Cron Test"

	_Cron   *cron.Cron
	_Config *viper.Viper
	_Server *http.Server
)

func init() {
	// _InstanceId = wrap.RandString(16)
	_Cron = cron.New(cron.WithSeconds())
	_Instance = misc.RandString(4)
	_StaticDirs = make([]wrap.StaticDir, 0)
}

func LoadBuildInfo(info [][2]string) {
	_BuildInfo = info
}

func logStartup(logger *zap.Logger, parameters map[string]interface{}) {
	fields := make([]zap.Field, 0, len(_BuildInfo))
	for _, v := range _BuildInfo {
		fields = append(fields, zap.String(v[0], v[1]))
	}

	logger.Info("startup", zap.Any(
		"event",
		map[string]interface{}{"buildInfo": fields, "parameters": parameters},
	))
}

func setExpvars() {
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Format(time.RFC3339)
	}))

	// export memstats and cmdline by default
	//	expvar.Publish("memStats", expvar.Func(func() any {
	//		memStats := new(runtime.MemStats)
	//		runtime.ReadMemStats(memStats)
	//		return memStats
	//	}))
}
