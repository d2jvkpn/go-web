package main

import (
	"flag"
	// "fmt"
	"path/filepath"

	"github.com/d2jvkpn/go-web/pkg"
)

func main() {
	var (
		config  string
		addr    string
		release bool
	)

	flag.StringVar(&config, "config", filepath.Join("configs", "local.yaml"), "config file path")
	flag.StringVar(&addr, "addr", ":8080", "http serve address")
	flag.BoolVar(&release, "release", false, "run in release mode")

	flag.Parse()

	pkg.Run(config, addr, release)
}
