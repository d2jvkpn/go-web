package misc

import (
	"fmt"
	"testing"
)

func TestBase64Map(t *testing.T) {
	d1 := map[string]interface{}{
		"hello": "world",
		"value": 109,
	}
	str := Base64EncodeMap(d1)
	fmt.Println(">>> str:", str)

	d2, err := Base64DecodeMap(str)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(">>> d2:", d2)
}
