package misc

import (
	"bytes"
	"fmt"
	"path/filepath"
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
		fmt.Println(">>>", Stack(2, ""))
	}()

	_fn1()
}

func _fn1() {
	var mySlice []int
	j := mySlice[0]

	fmt.Printf("Hello, playground %d", j)
}

func Stack(skip int, prefix string) (slice [][2]string) {
	bts := bytes.TrimSpace(debug.Stack())
	// fmt.Printf(">>>\n%s\n<<<\n", bts)
	re := regexp.MustCompile("\n.*\n\t.*")
	out := re.FindAllStringSubmatch(string(bts), -1)
	if skip < 2 {
		skip = 2
	}
	if len(out) < skip {
		return make([][2]string, 0)
	}
	slice = make([][2]string, 0, len(out)-skip)

	for i := skip; i < len(out); i++ {
		t := strings.Split(strings.TrimSpace(out[i][0]), "\n\t")
		if len(t) <= 1 {
			continue
		}
		if prefix != "" && !strings.HasPrefix(t[0], prefix) {
			continue
		}

		x := filepath.Base(strings.Fields(t[1])[0])
		slice = append(slice, [2]string{t[0], x})
	}

	return
}
