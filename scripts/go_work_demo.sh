#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

# https://go.googlesource.com/proposal/+/master/design/45713-workspace.md

####
ls go-web

mkdir -p releases
cd releases && go mod init releases

####
cd ${_path}
go work init go-web releases
ls go.work go.work.sum

####
cd releases

cat > main.go << EOF
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
EOF

go mod tidy
# touch configs/local.yaml
go build main.go
# go mod download, go mod graph, go mod verify, go mod why
# go mod edit, go mod init, go mod tidy, go mod vendor
