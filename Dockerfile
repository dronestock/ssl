FROM dockerproxy.com/library/golang:1.20.3-alpine AS golang
FROM dockerproxy.com/neilpang/acme.sh:3.0.5 AS acme



FROM ccr.ccs.tencentyun.com/storezhang/alpine:3.17.2


LABEL author="storezhang<华寅>" \
    email="storezhang@gmail.com" \
    qq="160290688" \
    wechat="storezhang" \
    description="SSL证书自动更新插件，使用ACME生成证书，自动替换各个平台的证书，包括：1、创世云"


# 复制文件
COPY --from=golang /usr/local/go/bin/go /usr/local/go/bin/go
COPY --from=golang /usr/local/go/pkg /usr/local/go/pkg
COPY --from=golang /usr/local/go/src /usr/local/go/src
COPY --from=acme /root/.acme.sh /opt/neilpang/acme


# 配置环境变量
ENV PATH ${PATH}:/usr/local/go/bin:/opt/neilpang/acme
ENV GO /var/lib/go
ENV GOPATH ${GO}/path
ENV GOCACHE ${GO}/cache
ENV GOPROXY https://goproxy.cn,https://mirrors.aliyun.com/goproxy,https://proxy.golang.com.cn,direct
