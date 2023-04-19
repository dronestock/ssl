FROM dockerproxy.com/library/golang:1.20.3-alpine AS golang
FROM dockerproxy.com/neilpang/acme.sh:3.0.5 AS acme


FROM ccr.ccs.tencentyun.com/storezhang/alpine:3.17.2 AS builder

COPY --from=golang /usr/local/go/bin/go /docker/usr/local/go/bin/go
COPY --from=golang /usr/local/go/pkg /docker/usr/local/go/pkg
COPY --from=golang /usr/local/go/src /docker/usr/local/go/src
COPY --from=acme /root/.acme.sh /docker/opt/neilpang/acme
COPY ssl /docker/usr/local/bin/



FROM ccr.ccs.tencentyun.com/storezhang/alpine:3.17.2


LABEL author="storezhang<华寅>" \
    email="storezhang@gmail.com" \
    qq="160290688" \
    wechat="storezhang" \
    description="SSL证书自动更新插件，使用ACME生成证书，自动替换各个平台的证书，包括：1、创世云"


# 复制文件
COPY --from=builder /docker /


RUN set -ex \
    \
    \
    \
    # 增加执行权限
    && chmod +x /usr/local/bin/ssl \
    \
    \
    \
    && rm -rf /var/cache/apk/*


# 执行命令
ENTRYPOINT /usr/local/bin/ssl


# 配置环境变量
ENV PATH ${PATH}:/usr/local/go/bin:/opt/neilpang/acme
