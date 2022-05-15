package ws

import (
	"github.com/gin-gonic/gin"
)

func Load(rg *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	open := rg.Group("/ws/v1", handlers...)

	open.GET("/hello", hello)
}
