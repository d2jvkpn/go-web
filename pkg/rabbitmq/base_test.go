package x

import (
	"flag"
	"fmt"
	"testing"
)

// go test -run TestTopicQueue -- URLURL
func TestMain(m *testing.M) {
	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		_TestUri = args[0]
	}

	fmt.Println("~~~ _TestUri:", _TestUri)
}
