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

func TestStsUploadLocal_t1(t *testing.T) {
	var (
		link   string
		err    error
		result *StsResult
	)

	if result, err = testStsClient.GetSTS("xxxxxxxx", ""); err != nil {
		t.Fatal("GetSTS:", err)
	}

	if link, err = result.UploadLocal("hello.txt", "test/sts/hello1.txt"); err != nil {
		t.Fatal(err)
	}
	fmt.Println(">>> TestStsUploadLocal_t1:", link)
}

func TestStsUploadLocal_t2(t *testing.T) {
	var (
		link   string
		err    error
		result *StsResult
	)

	if result, err = testStsClient.GetSTS("xxxxxxxx", ""); err != nil {
		t.Fatal("GetSTS:", err)
	}

	link, err = result.UploadLocal("hello.txt", "aaaa/sts/hello1.txt")
	if err == nil {
		fmt.Println(">>> TestStsUploadLocal_t2:", link)
		t.Fatal("unexpected upload successed")
	}
}
