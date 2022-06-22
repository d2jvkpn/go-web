package misc

import (
	"testing"
)

func TestNTP(t *testing.T) {
	ser, _ := NewNetworkTimeServer(":8080", 1000)
	ser.Run()
}
