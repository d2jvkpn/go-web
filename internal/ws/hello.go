package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func hello0(ctx *gin.Context) {
	if ctx.Param("a") != "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	if ctx.GetHeader("Upgrade") != "websocket" && ctx.GetHeader("Connection") != "Upgrade" {
		ctx.String(http.StatusUpgradeRequired, "Upgrade Required")
		ctx.Abort()
		return
	}

	client := NewClient(
		ctx.Request.RemoteAddr, // ctx.ClientIP(),
		ctx.DefaultQuery("name", "World"),
		_MelHello,
	)

	// _ = _MelHello.HandleRequest(ctx.Writer, ctx.Request)
	_ = _MelHello.HandleRequestWithKeys(
		ctx.Writer, ctx.Request, map[string]interface{}{"client": client},
	)
}

func hello(ctx *gin.Context) {
	client := NewClient(
		ctx.Request.RemoteAddr, // ctx.ClientIP(),
		ctx.DefaultQuery("name", "World"),
		_MelHello,
	)

	// _ = _MelHello.HandleRequest(ctx.Writer, ctx.Request)
	_ = _MelHello.HandleRequestWithKeys(
		ctx.Writer, ctx.Request, map[string]interface{}{"client": client},
	)
}
