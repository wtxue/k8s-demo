FROM alpine:3.11.3

MAINTAINER xk <xk.136@qq.com>

RUN apk update \
    && apk upgrade \
    && apk add --no-cache busybox-extras tcpdump tzdata curl wget lsof iotop bind-tools strace file iproute2 drill \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  \
    && echo "Asia/Shanghai">>/etc/timezone \
    && rm -rf /var/cache/apk/*

COPY Dockerfile /Dockerfile

CMD ["sh"]
