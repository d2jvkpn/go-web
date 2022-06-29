#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

image=registry.cn-shanghai.aliyuncs.com/d2jvkpn/alpine:base

docker pull alpine
docker build -f ${_path}/Dockerfile --no-cache -t "$image" .
docker push $image
