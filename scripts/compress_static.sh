#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

# npm install terser -g
# apt install yui-compressor
terser app.js > app.min.js
yui-compressor style.css > static/css/style.min.css
