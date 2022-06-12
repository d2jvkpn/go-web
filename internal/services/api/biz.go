package api

import (
	// "log"
	"fmt"
	"net/http"

	. "github.com/d2jvkpn/go-web/pkg/resp"

	"github.com/gin-gonic/gin"
)

func ping(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
}

func hello(ctx *gin.Context) {
	// key := "Authorization"
	// log.Printf("~~~ Header %s: %s\n", key, ctx.GetHeader(key))
	// ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "data": gin.H{}})
	// Ok(ctx)
	name := "Jane Doe"

	if v := ctx.Param("name"); v != "" {
		name = v
	}

	JSON(ctx, gin.H{"name": name}, nil)
}

func login(ctx *gin.Context) {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"code": 0, "msg": "ok", "data": gin.H{"token": "xxxxxxxx"},
	//	})

	key := "X-Token" // "Authorization"
	val := ctx.GetHeader(key)
	// log.Printf("~~~ Header %s: %s\n", key, ctx.GetHeader(key))
	ctx.Set(KeyUserId, "xxxxxxxx")

	switch {
	case val == "":
		BadRequest(ctx, fmt.Errorf("missing header: %s", key))
	case len(val) != 8:
		Error(ctx, ErrBadRequest(fmt.Errorf("invalid X-Token"), Msg("invalid X-Token")))
	default:
		ctx.Set(KeyEvent, "user logined success")
		JSON(ctx, gin.H{key: val}, nil)
	}
}
