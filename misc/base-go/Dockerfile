FROM alpine:3.11.3

MAINTAINER xk <xk.136@qq.com>

RUN apk update \
    && apk upgrade \
    && apk add --no-cache busybox-extras tcpdump tzdata curl wget lsof iotop bind-tools strace file iproute2 drill \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  \
    && echo "Asia/Shanghai">>/etc/timezone \
    && rm -rf /var/cache/apk/*

# used for go pprof
ENV GOLANG_VERSION 1.14
ENV GOLANG_DOWNLOAD_URL https://studygolang.com/dl/golang/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 08df79b46b0adf498ea9f320a0f23d6ec59e9003660b4c9c1ce8e5e2c6f823ca
RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
    && echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
    && tar -C /usr/local -xzf golang.tar.gz \
    && rm golang.tar.gz
ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH

COPY Dockerfile /Dockerfile

CMD ["sh"]
