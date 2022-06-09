package misc

import (
	"path/filepath"
	"time"

	"golang.org/x/exp/constraints"
)

func Basename(cf string) (base string) {
	base = filepath.Base(cf)
	ext := filepath.Ext(base)
	base = base[:len(base)-len(ext)]
	return
}

func NowMs() string {
	return time.Now().Format("2006-01-02T15:04:05.000Z07:00")
}

func VectorIndex[T constraints.Ordered](list []T, v T) int {
	for i := range list {
		if list[i] == v {
			return i
		}
	}

	return -1
}

func EqualVector[T constraints.Ordered](arr1, arr2 []T) (ok bool) {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}

func UniqVector[T constraints.Ordered](arr []T) (list []T) {
	n := len(arr)
	list = make([]T, 0, n)

	if len(arr) == 0 {
		return list
	}

	mp := make(map[T]bool, n)
	for _, v := range arr {
		if !mp[v] {
			list = append(list, v)
			mp[v] = true
		}
	}

	return list
}

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = _LetterRunes[_Rand.Intn(len(_LetterRunes))]
	}

	return string(b)
}
