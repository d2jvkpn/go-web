package internal

import (
	"context"
	// "fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/d2jvkpn/go-web/internal/settings"
	"github.com/d2jvkpn/go-web/pkg/wrap"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ServeOption func(*gin.RouterGroup) error

func StaticDir(dir, local string, listDir bool) ServeOption {
	return func(rg *gin.RouterGroup) (err error) {
		if listDir {
			rg.StaticFS(dir, http.Dir(local))
		} else {
			rg.Static(dir, local)
		}
		return nil
	}
}

func AppendServeOptions(opts ...ServeOption) {
	if len(opts) == 0 {
		return
	}

	_ServeOptions = append(_ServeOptions, opts...)
}

func Load(fp string, release bool) (err error) {
	var (
		engi *gin.Engine
	)

	if _Config, err = wrap.ReadConfigFile("config", fp); err != nil {
		return
	}
	settings.Config = _Config

	if !release {
		_ = os.Setenv("APP_DebugMode", "true")
	} else {
		_ = os.Setenv("APP_DebugMode", "false")
	}
	_Release = release

	settings.Logger = wrap.NewLogger("logs/go-web_api.log", zap.InfoLevel, 256, nil)

	if err = _SetupCrons(); err != nil {
		return err
	}

	setExpvars()

	if engi, err = NewEngine(_Release); err != nil {
		return err
	}
	_Server = &http.Server{ // TODO: set consts in base.go
		ReadTimeout:       HTTP_ReadTimeout,
		WriteTimeout:      HTTP_WriteTimeout,
		ReadHeaderTimeout: HTTP_ReadHeaderTimeout,
		MaxHeaderBytes:    HTTP_MaxHeaderBytes,
		// Addr:              addr,
		Handler: engi,
	}

	return
}

func Serve(addr string, parameters map[string]interface{}) (err error) {
	_Cron.Start()
	logStartup(settings.Logger.Logger, parameters)

	log.Printf(">>> HTTP server listening on %s\n", addr)
	_Server.Addr = addr
	if err = _Server.ListenAndServe(); err == http.ErrServerClosed {
		err = nil
	}

	return
}

func Down() {
	var err error

	log.Println("<<< Stop Cron Job")
	_Cron.Stop()

	if _Server != nil {
		log.Println("<<< Shutdown HTTP Server")
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		if err = _Server.Shutdown(ctx); err != nil {
			log.Printf("!!! _Server.Shutdown: %v\n", err)
		}
		cancel()
	}

	log.Println("<<< Close Loggers")
	// close other goroutines or services
	settings.Logger.Down()
}
