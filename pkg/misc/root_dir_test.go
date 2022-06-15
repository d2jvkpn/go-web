package misc

import (
	"fmt"
	"testing"
)

func TestRootDir(t *testing.T) {
	p, err := RootDir()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(p)
}

func TestRootFile(t *testing.T) {
	p, err := RootFile("config", "a.toml")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(p)
}
