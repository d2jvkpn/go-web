package internal

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/d2jvkpn/goapp/pkg/misc"

	"github.com/spf13/viper"
)

var (
	_Relase   bool
	_Config   *viper.Viper
	_Server   *http.Server
	BuildInfo [][2]string
)

func Load(fp string, release bool) (err error) {
	if _Config, err = misc.ReadConfigFile("config", fp); err != nil {
		return
	}
	_Relase = release

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
