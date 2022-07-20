package misc

import (
// "fmt"
)

func SegnmentsDiv(length, num int) (segs [][2]int) {
	if length == 0 || num == 0 {
		return nil
	}

	size := length / num
	segs = make([][2]int, num)

	for i := 0; i < num; i++ {
		segs[i] = [2]int{i * size, i*size + size}
	}

	if segs[len(segs)-1][1] != length {
		segs[len(segs)-1][1] = length
	}

	return
}
