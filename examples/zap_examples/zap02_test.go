package main

import (
	// "fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestZap02(t *testing.T) {
	var (
		// err    error
		logger *zap.Logger
	)

	config := zapcore.EncoderConfig{
		MessageKey:   "msg",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		FunctionKey:  "func",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.RFC3339NanoTimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	// opts := []zap.Option{zap.AddCallerSkip(0), zap.AddCaller()}
	// loggerConfig.Build(opts...)

	lf := &lumberjack.Logger{
		Filename:  "./logs/zap_example.log",
		LocalTime: true,
		MaxSize:   100, // megabytes
		// MaxBackups: 3,
		// MaxAge:     1, //days
		// Compress:   true, // disabled by default
	}

	w := zapcore.AddSync(lf)

	core := zapcore.NewCore(zapcore.NewJSONEncoder(config), w, zap.InfoLevel)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	defer func() {
		var err error

		if err = logger.Sync(); err != nil {
			log.Printf("zap.Logger.Sync: %v\n", err)
		}

		if err = lf.Close(); err != nil {
			log.Printf("lumberjack.Logger.Close: %v\n", err)
		}
	}()

	logger.Warn(
		"UserLogin",
		zap.String("id", uuid.NewString()), zap.Int64("code", 100),
		zap.String("entity", "rover"),
	)

	_ = lf.Rotate()

	logger.Warn(
		"UserLogout",
		zap.String("id", uuid.NewString()), zap.Int64("code", 101),
		zap.String("entity", "rover"),
		zap.Any("data", map[string]int{"hello": 2022, "world": 203}),
	)

	warn := func(msg string, code int64, entity string, data map[string]any) {
		logger.Warn(
			msg,
			zap.String("id", uuid.NewString()), zap.Int64("code", code),
			zap.String("entity", entity), zap.Any("data", data),
		)
	}

	warn("StreamerConnected", 0, "streamer0001", nil)
}
