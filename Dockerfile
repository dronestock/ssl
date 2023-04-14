FROM dockerproxy.com/neilpang/acme.sh:3.0.5 AS acme



FROM storezhang/alpine:3.17.2


LABEL author="storezhang<华寅>" \
    email="storezhang@gmail.com" \
    qq="160290688" \
    wechat="storezhang" \
    description="SSL证书自动更新插件，使用ACME生成证书，自动替换各个平台的证书，包括：1、创世云"


# 复制文件
COPY plugin /bin


RUN set -ex \
    \
    \
    \
    && apk update \
    && apk --no-cache add docker \
    \
    \
    \
    # 增加执行权限
    && chmod +x /bin/plugin \
    \
    \
    \
    && rm -rf /var/cache/apk/*


# 执行命令
ENTRYPOINT /bin/plugin
