#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

curl -i -X POST localhost:8080/api/v1/open/login

curl -i -X POST localhost:8080/api/v1/open/login -H "X-Token: aa"

curl -i -X POST localhost:8080/api/v1/open/login -H "X-Token: aaaaaaaa"
