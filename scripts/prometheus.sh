#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})

# https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/
# https://hub.docker.com/r/prom/prometheus
# https://hub.docker.com/r/grafana/grafana/tags
# https://povilasv.me/prometheus-go-metrics/

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
cat > docker-compose.yaml << EOF
version: "3"

services:
  grafana:
    image: grafana/grafana:main
    container_name: grafana
    networks: [go-web]
    ports: ["3000:3000"]
    volumes:
    - grafana-storage:/var/lib/grafana

  go-web:
    #build:
    #  context: ./
    #  dockerfile: Dockerfile
    image: registry.cn-shanghai.aliyuncs.com/d2jvkpn/go-web:dev
    container_name: go-web_dev
    restart: always
    networks: [go-web]
    ports: ["8080:8080"]
    environment: ["TZ=Asia/Shanghai"]
    volumes:
    - ./configs/:/opt/go-web/configs/
    - ./logs/:/opt/go-web/logs/
    command: ["./main", "--config=configs/dev.yaml", "--addr=:8080", "--release"]

  prometheus:
    image: prom/prometheus:main
    container_name: prometheus
    networks: [go-web]
    ports:
    - 9090:9090
    volumes:
    - ./configs/prometheus.yaml:/etc/prometheus/prometheus.yaml
    - prometheus_data:/prometheus
    command:
    - '--config.file=/etc/prometheus/prometheus.yaml'
    - '--storage.tsdb.path=/prometheus'
    - '--web.console.libraries=/usr/share/prometheus/console_libraries'
    - '--web.console.templates=/usr/share/prometheus/consoles'
    restart: always

volumes:
  grafana-storage:
  prometheus_data:

networks:
  go-web:
    name: go-web_dev
    driver: bridge
    external: false
EOF

docker-compose up -d

docker exec -it go-web_dev curl prometheus:9090/metrics
docker exec -it go-web_dev curl go-web_dev:8080/prometheus

#### reset grafana password
docker exec -it grafana grafana-cli admin reset-admin-password PASSWORD

# grafana url: http:localhost:3000, admin PASSWORD
# add prometheus data source: http://prometheus:9090

curl http://localhost:8080/api/v1/open/hello

ls -alh /var/lib/docker/volumes/go-web_grafana-storage/_data
