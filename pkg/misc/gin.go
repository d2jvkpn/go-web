package misc

import (
	// "fmt"
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Cors(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")

	ctx.Header(
		"Access-Control-Allow-Headers", "Content-Type, Authorization",
	)
	// X-CSRF-Token

	ctx.Header(
		"Access-Control-Expose-Headers",
		"Access-Control-Allow-Origin, Access-Control-Allow-Headers, "+
			"Content-Type, Content-Length",
	)

	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}

func WsUpgrade(ctx *gin.Context) {
	if ctx.GetHeader("Upgrade") != "websocket" && ctx.GetHeader("Connection") != "Upgrade" {
		// fmt.Printf("~~~~ Headers: %v\n", ctx.Request.Header)
		ctx.String(http.StatusUpgradeRequired, "Upgrade Required")
		ctx.Abort()
		return
	}

	ctx.Next()
}

// https://github.com/thinkerou/favicon/blob/master/favicon.go
func ServeFile(bts []byte, typ, name string, ts ...time.Time) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var t time.Time

		if len(ts) > 0 {
			t = ts[0]
		}

		reader := bytes.NewReader(bts)

		ctx.Header("Content-Type", typ)
		http.ServeContent(ctx.Writer, ctx.Request, name, t, reader)
	}
}
