#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})


APP_ENV="$1"
gitBranch="$2"
PORT=$3

export APP_ENV=${APP_ENV} gitBranch=${gitBranch} PORT=${PORT}
envsubst < ${_path}/deploy.yaml > docker-compose.yaml

# docker-compose pull
[[ ! -z "$(docker ps --all --quiet --filter name=goapp_$APP_ENV)" ]] &&
  docker rm -f goapp_$APP_ENV
# docker-compose down for running containers only, not stopped containers

docker-compose up -d
