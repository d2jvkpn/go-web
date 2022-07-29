package misc

import (
	"fmt"
	"testing"
)

func TestSigningUrlMd5(t *testing.T) {
	sign := NewSigningUrlMd5("xxxxxxxx", "sign", true)
	query := sign.Sign(map[string]string{"a": "1", "b": "zzzz"})

	fmt.Println(">>> query:", query)

	if err := sign.Verify(query); err != nil {
		t.Fatal(err)
	}
}
