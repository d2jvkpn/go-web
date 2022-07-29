package aliyun

import (
	"fmt"
	"testing"
)

func TestOssUploadLocal(t *testing.T) {
	var (
		link string
		err  error
	)

	if link, err = testOssClient.UploadLocal("hello.txt", "test/oss/hello1.txt"); err != nil {
		t.Fatal(err)
	}

	fmt.Println(">>> TestOssUploadLocal:", link)
}

func TestOssUploadFromUrl(t *testing.T) {
	var (
		fsize int64
		link  string
		err   error
	)

	link, fsize, err = testOssClient.UploadFromUrl(
		"http://127.0.0.1:8000/hello.txt",
		"test/oss/hello2.txt",
	)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(">>> TestOssUploadFromUrl:", link, fsize)
}

func TestUploadDir(t *testing.T) {
	var (
		link string
		err  error
	)

	link, err = testOssClient.UploadDir("dir", "test/oss/dir-target", 4)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf(">>> TestUploadDir: %s\n", link)
}

func TestCopyFile(t *testing.T) {
	code, err := testOssClient.CopyFile("test/oss/hello1.txt", "test/oss/hello1.copy.txt")

	if err != nil {
		fmt.Println("!!! code =", code)
		t.Fatal(err)
	}
}
