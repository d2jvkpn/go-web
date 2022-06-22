#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})


mkdir -p configs

cat > configs/prometheus.yaml << EOF
global:
  scrape_interval:     15s
  evaluation_interval: 15s

scrape_configs:
- job_name: prometheus
  static_configs:
  - targets: ['localhost:9090']

- job_name: go-web_dev
  metrics_path: /prometheus
  static_configs:
  - targets:
    - go-web_dev:8080
EOF

cp scripts/docker-compose.yaml docker-compose.yaml

docker-compose up -d

docker exec -it grafana grafana-cli admin reset-admin-password PASSWORD

# grafana url: http:localhost:3000
# add prometheus data source: http://prometheus:9090
