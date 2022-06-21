package resp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/d2jvkpn/go-web/pkg/misc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func NewLogHandler(logger *misc.Logger) gin.HandlerFunc {
	gomod, _ := misc.RootModule()

	return func(ctx *gin.Context) {
		var (
			ok     bool
			code   int
			err    *HttpError
			event  any // TODO: using a conrete type
			fields []zap.Field
		)

		fields = make([]zap.Field, 0, 9)
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

		saveLog := func() {
			// ctx.Request.Referer()
			// ctx.GetHeader("User-Agent")
			latency := time.Since(start).Milliseconds()
			appendString("userId", ctx.GetString(KeyUserId))
			fields = append(fields, zap.Int("status", ctx.Writer.Status()))
			fields = append(fields, zap.Int64("latency", latency))

			if err, ok = misc.GetCtxValue[*HttpError](ctx, KeyError); ok {
				fields = append(fields, zap.Any(KeyError, err))
				code = err.Code
			}

			if event, ok = misc.GetCtxValue[any](ctx, KeyEvent); ok {
				fields = append(fields, zap.Any(KeyEvent, event))
			}

			switch {
			case code <= 0:
				logger.Info(requestId, fields...) // array fields[0:]...
			case code < 100:
				logger.Warn(requestId, fields...)
			default:
				logger.Error(requestId, fields...)
			}
		}

		defer func() {
			var intf interface{}
			if intf = recover(); intf == nil {
				return
			}

			stacks := misc.Stack(4, gomod)
			err = ErrServerError(fmt.Errorf("%v", intf), Skip(5))
			ctx.Set(KeyError, err)
			ctx.Set(KeyEvent, gin.H{"kind": "panic", "stacks": stacks})
			// TODO: alerting the developers
			JSON(ctx, nil, err)
			saveLog()
		}()

		ctx.Next()
		saveLog()
	}
}

func Log2Tsv(fp string, w io.Writer) (err error) {
	type Record struct {
		Time    string `json:"time"`
		Level   string `json:"level"`
		Ip      string `json:"ip"`
		Msg     string `json:"msg"`
		Method  string `json:"method"`
		Path    string `json:"path"`
		Query   string `json:"query"`
		UserId  string `json:"userId"`
		Status  int64  `json:"status"`
		Latency int64  `json:"latency"`
		// error any
		// event any
	}

	record2Str := func(r *Record) string {
		strs := []string{
			r.Time, r.Level, r.Ip, r.Msg,
			r.Method, r.Path, r.Query, r.UserId,
			strconv.FormatInt(r.Status, 10), strconv.FormatInt(r.Latency, 10),
		}

		return strings.Join(strs, "\t")
	}

	var (
		line    int64
		bts     []byte
		file    *os.File
		scanner *bufio.Scanner
		buf     *bytes.Buffer
		record  Record
	)

	if file, err = os.Open(fp); err != nil {
		return err
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	buf = bytes.NewBuffer(nil)

	for scanner.Scan() {
		line++
		bts = scanner.Bytes()
		record = Record{}
		if err = json.Unmarshal(bts, &record); err != nil {
			return fmt.Errorf("readline %d: %w", line, err)
		}

		errorText := gjson.GetBytes(bts, KeyError).String()
		eventText := gjson.GetBytes(bts, KeyEvent).String()

		if _, err = buf.WriteString(record2Str(&record)); err != nil {
			return err
		}
		if err = buf.WriteByte('\t'); err != nil {
			return err
		}
		if _, err = buf.WriteString(errorText); err != nil {
			return err
		}
		if err = buf.WriteByte('\t'); err != nil {
			return err
		}
		if _, err = buf.WriteString(eventText); err != nil {
			return err
		}
		if err = buf.WriteByte('\n'); err != nil {
			return err
		}

		if _, err = w.Write(buf.Bytes()); err != nil {
			return err
		}
		buf.Reset()
	}

	return
}
