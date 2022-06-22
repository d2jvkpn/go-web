package api

import (
	"time"

	. "github.com/d2jvkpn/go-web/pkg/resp"

	"github.com/gin-gonic/gin"
)

func Load(rg *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	open := rg.Group("/api/v1/open", handlers...)

	open.GET("/timeout", func(ctx *gin.Context) {
		time.Sleep(20 * time.Second)
		JSON(ctx, gin.H{"code": 0, "msg": "ok"}, nil)
	})

	open.GET("/panic", func(ctx *gin.Context) {
		a, b := 1, 0
		result := a / b
		JSON(ctx, gin.H{"result": result}, nil)
	})

	open.POST("/login", login)
	open.GET("/hello", hello)
	open.GET("/hello/:name", hello)
}
