package resp

import (
	// "fmt"
	"os"
	"testing"

	"github.com/d2jvkpn/go-web/pkg/misc"
)

func TestLog2Tsv(t *testing.T) {
	fp, err := misc.RootFile("logs", "go-web_api.log")
	if err != nil {
		t.Fatal(err)
	}

	if err := Log2Tsv(fp, os.Stdout); err != nil {
		t.Fatal(err)
	}
}
