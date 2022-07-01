#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

convert -resize x32 -gravity center -crop 16x16+0+0 internal/static/assets/favicon.png \
  -flatten -colors 256 -background transparent internal/static/assets/favicon.ico
