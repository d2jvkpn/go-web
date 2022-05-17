package ws

import (
	"github.com/gin-gonic/gin"
)

func Load(rg *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	rg.GET(`/:a/`, hello) // router ws://localhost:8080//

	open := rg.Group("/ws/v1", handlers...)

	open.GET("/hello", hello)
}
