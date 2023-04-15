FROM dockerproxy.com/library/golang:1.20.3-alpine AS golang



FROM ccr.ccs.tencentyun.com/storezhang/alpine:3.17.2


LABEL author="storezhang<华寅>" \
    email="storezhang@gmail.com" \
    qq="160290688" \
    wechat="storezhang" \
    description="自动SSL证书插件运行环境"


# 复制文件
COPY --from=golang /usr/local/go/bin/go /usr/local/go/bin/go
COPY --from=golang /usr/local/go/pkg /usr/local/go/pkg
COPY --from=golang /usr/local/go/src /usr/local/go/src


# 配置环境变量
ENV PATH ${PATH}:/usr/local/go/bin
ENV GO /var/lib/go
ENV GOPATH ${GO}/path
ENV GOCACHE ${GO}/cache
ENV GOLANGCI_LINT_CACHE ${GO}/linter
ENV GOPROXY https://goproxy.cn,https://mirrors.aliyun.com/goproxy,https://proxy.golang.com.cn,direct
