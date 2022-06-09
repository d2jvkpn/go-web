package misc

import (
	"math/rand"
	"time"
)

var (
	_Rand        = rand.New(rand.NewSource(time.Now().UnixNano()))
	_LetterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)
