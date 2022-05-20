package site

import (
	"github.com/gin-gonic/gin"
)

func Load(rg *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	route := rg.Group("/site/", handlers...)

	route.GET("/not_found", notFound)
}
