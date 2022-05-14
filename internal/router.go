package internal

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoadAPI(rg *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	open := rg.Group("/api/v1/open", handlers...)

	open.GET("/hello", func(ctx *gin.Context) {
		key := "Authorization"
		fmt.Printf("~~~ Header %s: %s\n", key, ctx.GetHeader(key))
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": gin.H{}})
	})

	open.POST("/login", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0, "message": "ok", "data": gin.H{"token": "xxxxxxxx"},
		})
	})
}
