package resp

import (
	// "fmt"
	// "encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (logger *Logger) HandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			code   int
			err    *HttpError
			fields []zap.Field
		)

		fields = make([]zap.Field, 0, 8)
		appendString := func(key, val string) {
			fields = append(fields, zap.String(key, val))
		}

		start := time.Now()
		requestId := uuid.NewString()
		ctx.Set(KeyRequestId, requestId)
		appendString("ip", ctx.ClientIP())
		appendString("method", ctx.Request.Method)
		appendString("path", ctx.Request.URL.Path)
		appendString("query", ctx.Request.URL.RawQuery)

		ctx.Next()
		// TODO: handle panic

		// ctx.Request.Referer()
		// ctx.GetHeader("User-Agent")
		latency := time.Since(start).Milliseconds()
		appendString("userId", ctx.GetString(KeyUserId))
		fields = append(fields, zap.Int("status", ctx.Writer.Status()))
		fields = append(fields, zap.Int64("latency", latency))

		if intf, ok := ctx.Get(KeyError); ok {
			if err, ok = intf.(*HttpError); ok && err != nil {
				code = err.Code
				// bts, _ := json.Marshal(err)
				// fields[7] = zap.String("error", string(bts))
				fields = append(fields, zap.Any("error", err))
			}
		}

		switch {
		case code == 0:
			logger.Info(requestId, fields[0:]...)
		case code < 100:
			logger.Warn(requestId, fields[0:]...)
		default:
			logger.Error(requestId, fields[0:]...)
		}
	}
}
