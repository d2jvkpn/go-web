package misc

import (
	"path/filepath"
)

func Basename(cf string) (base string) {
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
