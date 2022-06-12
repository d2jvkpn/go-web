package misc

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	RFC3339ms = "2006-01-02T15:04:05.000Z07:00"
)

// a simple log writer
type LogWriter struct {
	fp          string
	send2stdout bool
	buf         *bytes.Buffer
	file        *os.File
	mutex       *sync.Mutex
}

type logWriter struct{}

func (w *logWriter) Write(bts []byte) (int, error) {
	return fmt.Printf("%s %s\n", time.Now().Format(RFC3339ms), bytes.TrimSpace(bts))
}

func RegisterDefaultLogFmt() {
	w := new(logWriter)
	log.SetFlags(0)
	log.SetOutput(w)
}

func NewLogWriter(prefix string, send2stdout bool) (lw *LogWriter, err error) {
	tag, bts := time.Now().Format("2006-01-02_15-04-05.000"), make([]byte, 0, 1024)

	lw = &LogWriter{
		fp:          prefix + "." + strings.Replace(tag, ".", "_", 1) + ".log",
		send2stdout: send2stdout,
		buf:         bytes.NewBuffer(bts),
	}

	if err = os.MkdirAll(filepath.Dir(prefix), 0755); err != nil {
		return nil, err
	}

	if lw.file, err = os.Create(lw.fp); err != nil {
		return nil, err
	}

	lw.mutex = new(sync.Mutex)
	return lw, nil
}

func (lw *LogWriter) Write(bts []byte) (int, error) {
	// time.RFC3339
	lw.mutex.Lock()
	defer lw.mutex.Unlock()

	// ?? check buffer size
	lw.buf.WriteString(time.Now().Format(RFC3339ms))
	lw.buf.WriteByte(' ')
	lw.buf.Write(bytes.TrimSpace(bts))
	lw.buf.WriteByte('\n')
	// bts = []byte(fmt.Sprintf("%s %s", , bts))
	n, err := lw.file.Write(lw.buf.Bytes())
	if lw.send2stdout {
		os.Stdout.Write(lw.buf.Bytes())
	}
	lw.buf.Reset()

	return n, err
}

// set as output of log pkg
func (lw *LogWriter) Register() {
	log.SetFlags(0)
	log.SetOutput(lw)
}

func (lw *LogWriter) Close() (err error) {
	err = lw.file.Close()
	w := new(logWriter)
	log.SetOutput(w)
	return err
}
