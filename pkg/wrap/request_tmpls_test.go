package wrap

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestLoadRequestTmpls(t *testing.T) {
	item, err := LoadRequestTmpls("config", "request_tmpls.yaml")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v", item)
}

func TestDoRequest(t *testing.T) {
	item, err := LoadRequestTmpls("config", "request_tmpls.yaml")
	if err != nil {
		t.Fatal(err)
	}

	tmpls, err := item.Match("hello")
	if err != nil {
		t.Fatal(err)
	}

	if statusCode, body, err := item.Request(tmpls[0]); err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("statusCode: %d\nbody: %s\n", statusCode, body)
	}
}

func TestIsJSON(t *testing.T) {
	{
		var js json.RawMessage
		fmt.Println(json.Unmarshal([]byte(`aaa`), &js))
	}

	{
		var js json.RawMessage
		fmt.Println(json.Unmarshal([]byte(`{"a": 1}`), &js))
	}

	{
		var js json.RawMessage
		fmt.Println(json.Unmarshal([]byte(`{a: 1}`), &js))
	}
}
