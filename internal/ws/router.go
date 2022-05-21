package ws

import (
	"github.com/gin-gonic/gin"
)

func Load(rg *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	rg.GET(`/:a/`, hello) // router ws://localhost:8080//

	ws := rg.Group("/ws", handlers...)
	ws.GET("/v1/hello", hello)
}
