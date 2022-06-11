#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})


gitBranch="$1"
app_env="$2"
port=$3

export app_env=${app_env} gitBranch=${gitBranch} port=${port}
envsubst < ${_path}/deploy.yaml > docker-compose.yaml

# docker-compose pull
[[ ! -z "$(docker ps --all --quiet --filter name=goapp_$app_env)" ]] &&
  docker rm -f goapp_$app_env
# docker-compose down for running containers only, not stopped containers

docker-compose up -d
