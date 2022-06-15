package misc

import (
	"bytes"
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"
)

func _fn2() {
	defer func() {
		var intf any

		if intf = recover(); intf == nil {
			return
		}
		fmt.Println("!!!", intf)

		slice := ParseStack(debug.Stack())
		fmt.Println(">>>", slice)
	}()

	_fn1()
}

func _fn1() {
	var mySlice []int
	j := mySlice[0]

	fmt.Printf("Hello, playground %d", j)
}

func ParseStack(bts []byte) (slice [][2]string) {
	bts = bytes.TrimSpace(bts)
	// fmt.Printf(">>>\n%s\n<<<\n", bts)
	re := regexp.MustCompile("\n.*\n\t.*")
	out := re.FindAllStringSubmatch(string(bts), -1)
	slice = make([][2]string, 0, len(out))

	for i := range out {
		v := strings.TrimSpace(out[i][0])
		if t := strings.Split(v, "\n\t"); len(t) > 1 {
			slice = append(slice, [2]string{t[0], strings.Fields(t[1])[0]})
		}
	}

	return
}
