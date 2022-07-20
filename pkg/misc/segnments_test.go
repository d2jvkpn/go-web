package misc

import (
	"fmt"
	"testing"
)

func TestSegnmentsDiv(t *testing.T) {
	list := [][2]int{
		{100, 20},
		{93, 10},
		{108, 10},
	}

	for _, v := range list {
		result := SegnmentsDiv(v[0], v[1])
		fmt.Printf(">>> %v: %t\n    %v\n", v, v[1] == len(result), result)
	}
}
