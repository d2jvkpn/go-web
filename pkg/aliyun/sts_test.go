package aliyun

import (
	"fmt"
	"testing"
)

func TestStsResult(t *testing.T) {
	var (
		err    error
		result *StsResult
	)

	if result, err = testStsClient.GetSTS("xxxxxxxx", ""); err != nil {
		t.Fatal(err)
	}

	fmt.Println(">>> TestStsResult:", result)
}

func TestStsUpload(t *testing.T) {
	var (
		link   string
		err    error
		result *StsResult
	)

	if result, err = testStsClient.GetSTS("xxxxxxxx", ""); err != nil {
		t.Fatal("GetSTS:", err)
	}

	if link, err = result.Upload("hello.txt", "test/sts/hello1.txt"); err != nil {
		t.Fatal(err)
	}
	fmt.Println(">>> TestStsUpload:", link)
}
