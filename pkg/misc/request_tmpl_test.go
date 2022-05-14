package misc

import (
	"fmt"
	"testing"
)

func TestLoadRequestTmpls(t *testing.T) {
	item, err := LoadRequestTmpls("config", "../../cmd/api-test.yaml")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v", item)
}

func TestDoRequest(t *testing.T) {
	item, err := LoadRequestTmpls("config", "../../cmd/api-test.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if statusCode, body, err := item.Request("hello"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("statusCode: %d\nbody: %s\n", statusCode, body)
	}
}
