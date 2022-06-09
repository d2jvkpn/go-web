#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})


rm -r logs
go run main_02.go

cat logs/zap_example.log

jq -sr '. | .[] | [.time, .level, .caller, .func, .msg, .id, .code, .entity, .data | tostring] | @tsv' \
  logs/zap_example.log
