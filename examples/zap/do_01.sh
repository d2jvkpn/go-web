#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

####
mv go.go main_01.go
go get go.uber.org/zap
go get go.uber.org/zap/zapcore
go get gopkg.in/natefinch/lumberjack.v2
go get github.com/google/uuid

go run main_01.go &> data.json
jq . data.json
jq ".ts, .caller" data.json

jq -s '.' data.json

jq -s '.' data.json | jq -r '.[] | .ts, .caller'

jq -sr '. | .[] | [.ts, .level, .caller, .func, .msg, .entity, .data | tostring] | @tsv' \
  data.json > data.tsv

# level, time, caller, msg, code, entity, data{}
