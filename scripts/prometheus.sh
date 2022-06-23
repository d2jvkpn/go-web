#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

# https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/
# https://hub.docker.com/r/prom/prometheus
# https://hub.docker.com/r/grafana/grafana/tags

#### create prometheus job configs
mkdir -p configs

cat > configs/prometheus.yaml << EOF
global:
  scrape_interval:     15s
  evaluation_interval: 15s

scrape_configs:
- job_name: prometheus
  static_configs:
  - targets: ['prometheus:9090']

- job_name: go-web_dev
  metrics_path: /prometheus
  static_configs:
  - targets:
    - go-web_dev:8080
EOF

#### run services
cp scripts/docker-compose.yaml docker-compose.yaml

docker-compose up -d

docker exec -it go-web_dev sh
# curl prometheus:9090/metrics
# curl go-web_dev:8080/prometheus

#### reset grafana password
docker exec -it grafana grafana-cli admin reset-admin-password PASSWORD

# grafana url: http:localhost:3000, admin PASSWORD
# add prometheus data source: http://prometheus:9090

curl http://localhost:8080/api/v1/open/hello

ls -alh /var/lib/docker/volumes/go-web_grafana-storage/_data
