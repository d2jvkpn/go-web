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
- golang generic, generate, benchmark, build information
- logging, telemetry, data collection(BI)
- monitoring
- databases: relational database, cache system, NoSQL, migration
- toolchains: api test, ws test, subcommands, configurations

- sevices:
  - http: register, login, logoff, basic auth, jwt-token, file upload and download, etc.
  - websocket
  - grpc
  - messaging queue: kafka

- deployments:
  - docker (build image, docker-compose)
  - self-host docker image registry
  - ansible
  - kubernetes


#### 3. Subsystems
- logging: api, cron job, bussiness data;
- cron job;
- telemetry, monitoring and alerting;
- messaging queue;
- memory cache;
- rpc;
