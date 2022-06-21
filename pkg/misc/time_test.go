package misc

import (
	"fmt"
	"testing"
)

func TestParseTime(t *testing.T) {
	strs := []string{
		"2022-06-22",
		"05:32:22",
		"2022-06-22T05:32:03",
		"2022-06-22 05:32:03",
	}

	for _, str := range strs {
		tm, err := ParseDatetime(str)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(tm)
	}
}
