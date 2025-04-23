FROM dockerproxy.com/neilpang/acme.sh:3.1.1 AS acme
FROM ccr.ccs.tencentyun.com/storezhang/alpine:3.20.0 AS builder

COPY --from=acme /root/.acme.sh /docker/opt/neilpang/acme
COPY ssl /docker/usr/local/bin/



FROM ccr.ccs.tencentyun.com/storezhang/alpine:3.20.0


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
    && apk update \
    \
    # 安装依赖库 \
    && apk --no-cache add socat openssl curl \
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
