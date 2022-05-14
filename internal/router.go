package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoadAPI(engi *gin.Engine, handlers ...gin.HandlerFunc) {
	open := engi.Group("/api/open", handlers...)

	open.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
	})
}
