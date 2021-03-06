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
	"github.com/d2jvkpn/go-web/pkg/wrap"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func NewLogHandler(logger *wrap.Logger, name string) gin.HandlerFunc {
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
		appendString("requestId", requestId)
		appendString("ip", ctx.ClientIP())
		appendString("method", ctx.Request.Method)
		appendString("path", ctx.Request.URL.Path)
		appendString("query", ctx.Request.URL.RawQuery)

		saveLog := func() {
			// ctx.Request.Referer()
			// ctx.GetHeader("User-Agent")
			latencyUs := time.Since(start).Microseconds()
			appendString("userId", ctx.GetString(KeyUserId))
			fields = append(fields, zap.Int("status", ctx.Writer.Status()))
			fields = append(fields, zap.Int64("latencyUs", latencyUs))

			if err, ok = wrap.GetCtxValue[*HttpError](ctx, KeyError); ok {
				fields = append(fields, zap.Any(KeyError, err))
				code = err.Code
			}

			if event, ok = wrap.GetCtxValue[any](ctx, KeyEvent); ok {
				fields = append(fields, zap.Any(KeyEvent, event))
			}

			switch {
			case code <= 0:
				logger.Info(name, fields...) // array fields[0:]...
			case code < 100:
				logger.Warn(name, fields...)
			default:
				logger.Error(name, fields...)
			}
		}

		defer func() {
			var intf interface{}
			if intf = recover(); intf == nil {
				return
			}

			stacks := misc.Stack(gomod)
			err = ErrServerError(fmt.Errorf("%v", intf), Skip(5))
			ctx.Set(KeyError, err)
			ctx.Set(KeyEvent, gin.H{"kind": "panic", "stacks": stacks})
			// TODO: alerting the developers
			JSON(ctx, nil, err)
			saveLog()
		}()

		ctx.Status(1000)
		ctx.Next()

		select {
		case <-ctx.Done():
			// fmt.Println("~~~~")
		default:
		}

		saveLog()
	}
}

func Log2Tsv(fp string, w io.Writer, times ...time.Time) (err error) {
	type Record struct {
		Time      string `json:"time"`
		Level     string `json:"level"`
		RequestId string `json:"requestId"`
		Ip        string `json:"ip"`
		Msg       string `json:"msg"`
		Method    string `json:"method"`
		Path      string `json:"path"`
		Query     string `json:"query"`
		UserId    string `json:"userId"`
		Status    int64  `json:"status"`
		LatencyUs int64  `json:"latencyUs"`
		// error any
		// event any
	}

	var (
		line    int64
		bts     []byte
		tm      time.Time
		start   time.Time
		end     time.Time
		file    *os.File
		scanner *bufio.Scanner
		buf     *bytes.Buffer
		record  Record
	)

	record2Str := func(r *Record) string {
		strs := []string{
			r.Time, r.Level, r.Ip, r.Msg, r.RequestId,
			r.Method, r.Path, r.Query, r.UserId,
			strconv.FormatInt(r.Status, 10), strconv.FormatInt(r.LatencyUs, 10),
		}

		return strings.Join(strs, "\t")
	}

	if len(times) > 1 {
		start, end = times[0], times[1]
	} else if len(times) == 1 {
		start = times[0]
	}

	filter := func(t time.Time) bool {
		switch {
		case start.IsZero() && end.IsZero():
			return true
		case !start.IsZero() && end.IsZero():
			return t.Sub(start) >= 0
		case start.IsZero() && !end.IsZero():
			return t.Sub(end) <= 0
		default:
			return t.Sub(start) >= 0 && t.Sub(end) <= 0
		}
	}

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

		if tm, err = time.Parse(time.RFC3339, record.Time); err != nil {
			return err
		}
		if !filter(tm) {
			continue
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
