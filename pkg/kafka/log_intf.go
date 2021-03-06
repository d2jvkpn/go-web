package kafka

import (
	"fmt"
	"strings"
	"time"
)

// customize yourself, no concurrency safe guaranteed
type LogIntf interface {
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

type Logger struct{}

func NewLogger() (logger *Logger) {
	logger = new(Logger)
	// logger.l = log.New(os.Stdout, prefix, log.Lmsgprefix | log.Lshortfile)
	// log.SetOutput(logger.l)
	return logger
}

func (logger *Logger) Printf(format string, a ...any) (int, error) {
	t := time.Now().Format(RFC3339ms)
	return fmt.Printf(t+" "+strings.TrimSpace(format)+"\n", a...) // bytes.TrimSpace(bts)
}

func (logger *Logger) Info(format string, a ...any) {
	logger.Printf(" [INFO] "+format, a...)
}

func (logger *Logger) Warn(format string, a ...any) {
	logger.Printf(" [WARN] "+format, a...)
}

func (logger *Logger) Error(format string, a ...any) {
	logger.Printf(" [ERROR] "+format, a...)
}
