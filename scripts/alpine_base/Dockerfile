####
FROM alpine
ENV TZ="Asia/Shanghai"
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN { apk --no-cache update && apk --no-cache upgrade && apk --no-cache add tzdata; }
