package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func hello(ctx *gin.Context) {
	key := "Authorization"
	log.Printf("~~~ Header %s: %s\n", key, ctx.GetHeader(key))
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": gin.H{}})
}

func login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0, "msg": "ok", "data": gin.H{"token": "xxxxxxxx"},
	})
}
