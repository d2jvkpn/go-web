package misc

import (
	"errors"
	"fmt"
	// "log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	_Rand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	_LetterRunes []rune = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func BasenameWithoutExt(cf string) (base string) {
	base = filepath.Base(cf)
	ext := filepath.Ext(base)
	base = base[:len(base)-len(ext)]
	return
}

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = _LetterRunes[_Rand.Intn(len(_LetterRunes))]
	}

	return string(b)
}

func ListenOSSignal(ch chan int, sgs ...os.Signal) {
	// linux support syscall.SIGUSR2
	quit := make(chan os.Signal, 1)

	if len(sgs) == 0 {
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	} else {
		signal.Notify(quit, sgs...)
	}

	select {
	case _ = <-quit: // sig := <-quit
		// log.Printf("received os signal: %v\n", sig)
		ch <- 0
	case _ = <-ch:
		// log.Printf("received value: %v\n", value)
		ch <- -1
	}

	return
}

func FileSize2Str(n int64) string {
	switch {
	case n <= 0:
		return "0"
	case n < 1<<10:
		return fmt.Sprintf("%dB", n)
	case n >= 1<<10 && n < 1<<20:
		return fmt.Sprintf("%dK", n>>10)
	case n >= 1<<20 && n < 1<<30:
		return fmt.Sprintf("%dM", n>>20)
	default:
		return fmt.Sprintf("%dG", n>>30)
	}
}

func CheckDuplicateFilename(p string) (out string, err error) {
	var (
		i         int
		base, ext string
	)

	ext = filepath.Ext(p)
	base = p[0:(len(p) - len(ext))]
	i, out = 1, p
	for {
		// fmt.Println(i, out)
		if _, err = os.Stat(out); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return out, nil
			}
			return "", err
		}
		i++
		out = fmt.Sprintf("%s-%d%s", base, i, ext)
	}
}
