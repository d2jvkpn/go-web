package internal

import (
	"context"
	"embed"
	"log"
	"net/http"
	"time"

	"github.com/d2jvkpn/goapp/pkg/misc"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	_Release  bool
	_Config   *viper.Viper
	_Server   *http.Server
	BuildInfo [][2]string

	//go:embed static
	_Static embed.FS
	//go:embed templates
	_Templates embed.FS
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
	_Release = release

	return
}

func Shutdown() {
	var err error

	if _Server != nil {
		log.Println("<<< Shutdown HTTP Server")
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		if err = _Server.Shutdown(ctx); err != nil {
			log.Printf("!!! _Server.Shutdown: %v\n", err)
		}
		cancel()
	}
}
