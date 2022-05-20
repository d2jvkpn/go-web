package ws

import (
	"github.com/gin-gonic/gin"
)

func Load(rg *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	ws := rg.Use(handlers...)

	ws.GET(`/:a/`, hello) // router ws://localhost:8080//
	ws.GET("/ws/v1/hello", hello)
}
