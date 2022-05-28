####
FROM golang:1.18-alpine AS builder
LABEL stage=goapp_builder

ARG BuildTime
ARG BuildBranch
ARG CommitHash
ARG CommitTime

ENV buildTime=$buildTime
ENV gitBranch=$gitBranch
ENV gitCommit=$gitCommit
ENV gitTime=$gitTime
ENV gitTreeState=$gitTreeState
ENV TZ="Asia/Shanghai"

RUN { apk --no-cache update && apk --no-cache upgrade && apk --no-cache add tzdata; } &> /dev/null

WORKDIR /opt/goapp
COPY ./main.go ./go.mod ./go.sum ./project.yaml  /opt/goapp/
COPY ./internal /opt/goapp/internal
COPY ./pkg      /opt/goapp/pkg
COPY ./cmd      /opt/goapp/cmd

# in alpine, date doesn't parse %:z
RUN go build  -ldflags="          \
  -X main.buildTime=${buildTime}  \
  -X main.gitBranch=${gitBranch}  \
  -X main.gitCommit=${gitCommit}  \
  -X main.gitTime=${gitTime}      \
  -X main.gitTreeState=${gitTreeState}" \
  -o main main.go

####
FROM alpine

ENV TZ="Asia/Shanghai"

WORKDIR /opt/goapp
COPY --from=builder /opt/goapp/main /opt/goapp/main

RUN { apk --no-cache update && apk --no-cache upgrade && apk --no-cache add tzdata; } &> /dev/null

EXPOSE 8080 1030
CMD ["./main", "serve", "--addr=:8080"]
