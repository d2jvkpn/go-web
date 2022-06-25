package main

import (
	// "fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/d2jvkpn/go-web/pkg/wrap"

	"github.com/gin-gonic/gin"
)

var (
	counter int64
	ch      chan bool
)

func init() {
	ch = make(chan bool)
}

func main() {
	var (
		addr    string
		engine  *gin.Engine
		iroutes gin.IRoutes
	)

	addr = os.Args[1]
	engine = gin.Default()

	iroutes = engine.RouterGroup.Use(wrap.NewPrometheusMonitor(""))

	iroutes.GET("/prometheus", wrap.PrometheusFunc)

	iroutes.GET("/open", func(ctx *gin.Context) {
		c := atomic.AddInt64(&counter, 1)
		go func() {
			ch <- true
		}()
		ctx.JSON(http.StatusOK, gin.H{"counter": c})
	})

	iroutes.GET("/close", func(ctx *gin.Context) {
		c := atomic.LoadInt64(&counter)
		for i := int64(0); i < c; i++ {
			<-ch
			_ = atomic.AddInt64(&counter, -1)
		}

		ctx.JSON(http.StatusOK, gin.H{"counter": atomic.LoadInt64(&counter)})
	})

	engine.Run(addr)
}
