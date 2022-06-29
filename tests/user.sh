#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})


#### user register
curl -X POST http://localhost:8080/api/open/register -d '{"user":"rover","password":"aa123456"}'

#### user upload file(s)
curl -X PUT -u rover:aa123456 http://localhost:8080/api/auth/user/upload \
  -F "files=@./tests/user.sh" -H "Content-Type: multipart/form-data"

### curl http://localhost:8080/uploads/rover/2022-06-29/user.1656467925590_ebJp5ykmx1eKcwFn.sh

#### user unregister
curl -X POST -u rover:aa123456 http://localhost:8080/api/auth/user/unregister
