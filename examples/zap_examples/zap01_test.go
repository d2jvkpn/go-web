package main

import (
	// "fmt"
	"log"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZap01(t *testing.T) {
	var (
		err    error
		logger *zap.Logger
	)

	// zap.NewDevelopment()
	//	if logger, err = zap.NewProduction(); err != nil {
	//		log.Fatal(err)
	//	}
	//	defer func() {
	//		if err := logger.Sync(); err != nil {
	//			log.Printf("zap.Sync: %v\n", err)
	//		}
	//	}()

	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	config.EncoderConfig.FunctionKey = "func"
	// opts := []zap.Option{zap.AddCallerSkip(0), zap.AddCaller()}
	// loggerConfig.Build(opts...)

	if logger, err = config.Build(); err != nil {
		log.Fatal(err)
	}

	logger.Info("hello", zap.Int64("code", 100), zap.String("entity", "rover"))

	logger.Warn(
		"hello",
		zap.Int64("code", 101), zap.String("entity", "rover"),
		zap.Any("data", map[string]int{"hello": 2022, "world": 203}),
	)

	lg := logger.Named("membership")
	lg.Info("WORLD", zap.Int64("code", 100), zap.String("entity", "rover"))
}
