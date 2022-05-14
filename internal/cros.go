package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors(ctx *gin.Context) {
	method := ctx.Request.Method
	ctx.Header("Access-Control-Allow-Origin", "*")

	ctx.Header(
		"Access-Control-Allow-Headers", "Content-Type, Authorization",
	)
	// X-CSRF-Token

	ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	ctx.Header(
		"Access-Control-Expose-Headers",
		"Access-Control-Allow-Origin, Access-Control-Allow-Headers, "+
			"Content-Type, Content-Length",
	)

	ctx.Header("Access-Control-Allow-Credentials", "true")

	if method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
