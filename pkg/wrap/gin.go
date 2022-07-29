package wrap

import (
	"bytes"
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Cors(origin string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", origin)

		//		allowHeaders = []string{"Content-Type", "Authorization"}
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

//func WrapF(f http.HandlerFunc) gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		f(ctx.Writer, ctx.Request)
//	}
//}

//func WrapH(h http.Handler) gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		h.ServeHTTP(ctx.Writer, ctx.Request)
//	}
//}

func WriteJSON(ctx *gin.Context, bts []byte) (int, error) {
	ctx.Header("StatusCode", strconv.Itoa(http.StatusOK))
	ctx.Header("Status", http.StatusText(http.StatusOK))
	ctx.Header("Content-Type", "application/json") // ; charset=utf-8
	return ctx.Writer.Write(bts)
}

func Healthy(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
}

func Pprof(rg *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	dbg := rg.Group("/debug", handlers...)

	///
	buildInfo, _ := debug.ReadBuildInfo()
	dbg.GET("/build_info", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"buildInfo": buildInfo})
	})

	dbg.GET("/stats", gin.WrapH(expvar.Handler()))

	///
	dbg.GET("/pprof/", gin.WrapF(pprof.Index))
	for _, v := range []string{
		"allocs", "block", "goroutine", "heap", "mutex", "threadcreate",
	} {
		dbg.GET("/pprof/"+v, gin.WrapH(pprof.Handler(v)))
	}

	dbg.GET("/pprof/profile", gin.WrapF(pprof.Profile))
	dbg.GET("/pprof/trace", gin.WrapF(pprof.Trace))
	dbg.GET("/pprof/cmdline", gin.WrapF(pprof.Cmdline))
	dbg.GET("/pprof/symbol", gin.WrapF(pprof.Symbol))

	return
}

type StaticDir func(*gin.RouterGroup) error

func ServeStatic(dir, local string, listDir bool) StaticDir {
	return func(rg *gin.RouterGroup) (err error) {
		if listDir {
			rg.StaticFS(dir, http.Dir(local))
		} else {
			rg.Static(dir, local)
		}
		return nil
	}
}

func GinJwtHSAuth(auth *JwtHSAuth, handle func(*gin.Context, map[string]any) error,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			head string
			data map[string]any
			err  error
		)

		if head = ctx.Request.Header.Get("Authorization"); head == "" {
			// ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": -100, "msg": "unauthorized", "data": gin.H{},
			})
			return
		}
		if !strings.HasPrefix(head, "Bearer ") {
			// ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": -101, "msg": "unauthorized", "data": gin.H{},
			})
			return
		}

		if data, err = auth.Parse(head[7:]); err != nil {
			// ctx.AbortWithStatus(http.StatusForbidden)
			ctx.JSON(http.StatusForbidden, gin.H{
				"code": -102, "msg": "forbidden", "data": gin.H{},
			})
			return
		}

		if handle != nil {
			// use ctx.Set in handle, response inside handle if an error occurs
			if err = handle(ctx, data); err != nil {
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
