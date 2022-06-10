package resp

import (
	// "fmt"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	Writer *lumberjack.Logger
	config zapcore.EncoderConfig
	core   zapcore.Core
	*zap.Logger
}

func NewLogger(filename string, level zapcore.LevelEnabler, mbs int, skips ...int) (logger *Logger) {
	logger = new(Logger)

	logger.Writer = &lumberjack.Logger{
		Filename:  filename,
		LocalTime: true,
		MaxSize:   mbs, // megabytes
		// MaxBackups: 3,
		// MaxAge:     1, //days
		// Compress:   true, // disabled by default
	}

	logger.config = zapcore.EncoderConfig{
		MessageKey:   "msg",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		FunctionKey:  "func",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.RFC3339NanoTimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	logger.core = zapcore.NewCore(
		zapcore.NewJSONEncoder(logger.config),
		zapcore.AddSync(logger.Writer), level,
	)
	// zap.InfoLevel

	if len(skips) > 0 {
		logger.Logger = zap.New(logger.core, zap.AddCaller(), zap.AddCallerSkip(skips[0]))
	} else {
		logger.Logger = zap.New(logger.core)
	}

	return
}

func (logger *Logger) Down() {
	var err error

	if logger == nil {
		return
	}

	if err = logger.Sync(); err != nil {
		log.Printf("Logger.Sync: %v\n", err)
	}

	if err = logger.Writer.Close(); err != nil {
		log.Printf("Logger.Close: %v\n", err)
	}
}
