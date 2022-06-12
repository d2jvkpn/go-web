package api

import (
	"github.com/gin-gonic/gin"
)

func Load(rg *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	open := rg.Group("/api/v1/open", handlers...)

	open.POST("/login", login)
	open.GET("/hello", hello)
	open.GET("/hello/:name", hello)
}
