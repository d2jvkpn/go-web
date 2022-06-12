package misc

import (
	"bytes"
	"fmt"
	"io"
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
	fp    string
	w     io.Writer
	buf   *bytes.Buffer
	file  *os.File
	mutex *sync.Mutex
}

type LogPrinter struct{}

func (w *LogPrinter) Write(bts []byte) (int, error) {
	return fmt.Printf("%s %s\n", time.Now().Format(RFC3339ms), bytes.TrimSpace(bts))
}

func RegisterLogPrinter() {
	w := new(LogPrinter)
	log.SetFlags(log.Lshortfile | log.Lmsgprefix)
	log.SetOutput(w)
}

func NewLogWriter(prefix string, w io.Writer) (lw *LogWriter, err error) {
	tag, bts := time.Now().Format("2006-01-02_15-04-05.000"), make([]byte, 0, 1024)

	lw = &LogWriter{
		fp:  prefix + "." + strings.Replace(tag, ".", "_", 1) + ".log",
		w:   w,
		buf: bytes.NewBuffer(bts),
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
	if lw.w != nil {
		lw.w.Write(lw.buf.Bytes())
	}
	lw.buf.Reset()

	return n, err
}

// set as output of log pkg
func (lw *LogWriter) Register() {
	log.SetFlags(log.Lshortfile | log.Lmsgprefix)
	log.SetOutput(lw)
}

func (lw *LogWriter) Close() (err error) {
	return lw.file.Close()
}
