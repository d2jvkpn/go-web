#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

for i in {1..3}; do
    mkdir -p wk_v$i

    go test  -run none -bench ^BenchmarkLimiterV${i}_b1$ -outputdir wk_v$i \
      -cpuprofile=cpu.out -memprofile=mem.out -blockprofile=block.out

    for f in wk_v$i/*.out; do
        go tool pprof -svg $f > ${f%.out}.svg
    done
done
