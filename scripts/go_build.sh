#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

Program="go-web"

buildTime=$(date +'%FT%T%:z')
gitBranch="$(git rev-parse --abbrev-ref HEAD)" # current branch
gitCommit=$(git rev-parse --verify HEAD) # git log --pretty=format:'%h' -n 1
gitTime=$(git log -1 --format="%at" | xargs -I{} date -d @{} +%FT%T%:z)
gitTreeState="clean"

uncommitted=$(git status --short)
unpushed=$(git diff origin/$gitBranch..HEAD --name-status)
[[ ! -z "$uncommitted$unpushed" ]] && gitTreeState="dirty"

#if [[ $(printenv APP_GitForce) != "true" ]]; then
#    test -z "$uncommitted" || { echo "You have uncommitted changes!"; exit 1; }
#    #- test -z "$unpushed" || { echo "You have unpushed commits!"; exit 1; }
#fi

ldflags="\
  -X main.buildTime=${buildTime} \
  -X main.gitBranch=$gitBranch   \
  -X main.gitCommit=$gitCommit   \
  -X main.gitTime=$gitTime       \
  -X main.gitTreeState=$gitTreeState"

mkdir -p target
go build -ldflags="$ldflags" -o target/$Program main.go
echo "saved target/$Program"
# GOOS=windows GOARCH=amd64 go build -ldflags="$ldflags" -o target/$Program.exe main.go
ls -al target/
