package misc

import (
	"bytes"
	"fmt"
	"net/http"
	// "strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Cors(origin string, allowHeaders ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", origin)

		//		if len(allowHeaders) == 0 {
		//			allowHeaders = []string{"Content-Type", "Authorization"}
		//		}

		//		ctx.Header(
		//			"Access-Control-Allow-Headers", strings.Join(allowHeaders, ", "),
		//		)
		//		// Content-Type, Authorization, X-CSRF-Token

		//		exposeHeaders := []string{
		//			"Access-Control-Allow-Origin",
		//			"Access-Control-Allow-Headers",
		//			"Content-Type",
		//			"Content-Length",
		//		}

		//		ctx.Header("Access-Control-Expose-Headers", strings.Join(exposeHeaders, ", "))

		//		ctx.Header("Access-Control-Allow-Credentials", "true")
		//		ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.Next()
	}
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
func ServeFavicon(bts []byte, ts ...time.Time) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var t time.Time

		if len(ts) > 0 {
			t = ts[0]
		}

		reader := bytes.NewReader(bts)

		ctx.Header("Content-Type", "image/x-icon")
		http.ServeContent(ctx.Writer, ctx.Request, "favicon.ico", t, reader)
	}
}

func CacheControl(seconds int) gin.HandlerFunc {
	cc := fmt.Sprintf("public, max-age=%d", seconds)
	// strconv.FormatInt(time.Now().UnixMilli(), 10)
	etag := fmt.Sprintf(`"%d"`, time.Now().UnixMilli()) // must be a quoted string

	return func(ctx *gin.Context) {
		if ctx.Request.Method != "GET" {
			ctx.Next()
			return
		}

		ctx.Header("Cache-Control", cc)
		// browser send If-None-Match: etag, if unchanged, response 304
		ctx.Header("ETag", etag)
		ctx.Next()
	}
}

func GetCtxValue[T any](ctx *gin.Context, key string) (v T, ok bool) {
	var intf interface{}

	if intf, ok = ctx.Get(key); !ok {
		return
	}
	v, ok = intf.(T)
	return
}
