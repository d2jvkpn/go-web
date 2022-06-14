# go-web
a golang prototype application


#### 1. Up and Running
```bash
# go env -w GOPROXY="https://goproxy.cn,direct"

go get

# bash scripts/go_build.sh

make run
```

#### 2. Tips
- test: unit test, integrity test
- golang generic, generate, benchmark
- logging, telemetry, data collection(BI)
- monitoring
- databases: relational database, cache system, NoSQL, migration
- toolchains: api test, ws test, subcommand

- sevices:
  - http: register, login, logoff, basic auth, jwt-token, file upload and download, etc.
  - websocket
  - grpc
  - messaging queue: kafka

- deployments:
  - docker (build image, docker-compose)
  - clould provode mirror service, self-host registry
  - ansible
  - kubernetes
