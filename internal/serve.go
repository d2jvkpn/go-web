package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/d2jvkpn/go-web/pkg/misc"

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

func Load(fp string, release bool) (err error) {
	if _Config, err = misc.ReadConfigFile("config", fp); err != nil {
		return
	}

	if !release {
		_ = os.Setenv("APP_DebugMode", "true")
	} else {
		_ = os.Setenv("APP_DebugMode", "false")
	}
	_Release = release

	_ApiLogger = misc.NewLogger(
		fmt.Sprintf("logs/go-web_api.%s.log", _InstanceId),
		zap.InfoLevel, 100, nil,
	)

	if err = _SetupCrons(); err != nil {
		return err
	}

	_Cron.Start()

	return
}

func Serve(addr string) (err error) {
	var engi *gin.Engine

	if engi, err = NewEngine(_Release); err != nil {
		return err
	}

	_Server = &http.Server{
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    4 << 20,
		Addr:              addr,
		Handler:           engi,
	}

	if err = _Server.ListenAndServe(); err == http.ErrServerClosed {
		err = nil
	}

	return
}

func Down() {
	var err error

	log.Println("<<< Stop Cron Job")
	_Cron.Stop()

	log.Println("<<< Close Loggers")
	// close other goroutines or services
	_ApiLogger.Down()

	if _Server != nil {
		log.Println("<<< Shutdown HTTP Server")
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		if err = _Server.Shutdown(ctx); err != nil {
			log.Printf("!!! _Server.Shutdown: %v\n", err)
		}
		cancel()
	}
}
