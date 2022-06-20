package misc

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// get root path of project by recursively match go.mod in parent directory
func RootDir() (dir string, err error) {
	var (
		x, tmp string
		ms     []string
	)
	if x, err = os.Getwd(); err != nil {
		return "", err
	}

	for {
		if ms, err = filepath.Glob(filepath.Join(x, "go\\.mod")); err != nil {
			return "", nil
		}
		if len(ms) > 0 {
			return x, nil
		}

		if tmp = filepath.Dir(x); x == tmp {
			break
		}
		x = tmp
	}

	return "", fmt.Errorf("not found")
}

func RootFile(p2f ...string) (fp string, err error) {
	var dir string

	if dir, err = RootDir(); err != nil {
		return "", err
	}
	arr := make([]string, 0, len(p2f)+1)
	arr = append(arr, dir)
	arr = append(arr, p2f...)

	return filepath.Join(arr...), nil
}

func RootModule() (mod string, err error) {
	var (
		modf    string
		strs    []string
		file    *os.File
		scanner *bufio.Scanner
	)

	if modf, err = RootFile("go.mod"); err != nil {
		return "", err
	}

	if file, err = os.Open(modf); err != nil {
		return "", err
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	//    for scanner.Scan() {
	//        fmt.Println(scanner.Text())
	//    }
	_ = scanner.Scan()
	strs = strings.Fields(scanner.Text())

	return strs[len(strs)-1], nil
}
