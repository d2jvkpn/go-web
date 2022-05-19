package ws

import (
	"log"

	"github.com/gin-gonic/gin"
)

func hello(ctx *gin.Context) {
	client := NewClient(
		ctx.Request.RemoteAddr, // ctx.ClientIP(),
		ctx.DefaultQuery("name", "World"),
	)

	log.Printf("================ %s\n", client)

	// _ = _MelHello.HandleRequest(ctx.Writer, ctx.Request)
	_ = _MelHello.HandleRequestWithKeys(
		ctx.Writer, ctx.Request, map[string]interface{}{"client": client},
	)
}
