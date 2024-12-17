FROM registry.cn-shenzhen.aliyuncs.com/infrasync/library-alpine:3.13.12 AS files1
  # files1: 三方文件, files2: src编译
  # xz-utils > xz
  ARG TARGETPLATFORM
  RUN sed -i "s^dl-cdn.alpinelinux.org^mirrors.tuna.tsinghua.edu.cn^g" /etc/apk/repositories; \
    apk add bash curl xz
  ADD docker-compose*.tar.gz /build
  RUN \
  # test -z "$(echo $TARGETPLATFORM |grep arm)" && arch=amd64 || arch=arm64; \
  case ${TARGETPLATFORM} in \
    "linux/amd64")  arch=amd64    ;; \
    "linux/arm64")  arch=arm64     ;; \
    "linux/arm/v7") arch=arm   ;; \
  esac; \
  mkdir -p /app/docker-registry; \
  tar -zxf /build/docker-registry-*-$arch.tar.gz -C /app/docker-registry; \
  find /app/docker-registry -type f |sort; 

# ref https://gitee.com/g-golang/fk-distribution/blob/main/Dockerfile
# FROM alpine:${ALPINE_VERSION}
FROM registry.cn-shenzhen.aliyuncs.com/infrasync/library-alpine:3.13.12
RUN apk add --no-cache ca-certificates
# COPY cmd/registry/config-dev.yml /etc/docker/registry/config.yml
# COPY --from=binary /registry /bin/registry
COPY --from=files1 /app/docker-registry /app/docker-registry
VOLUME ["/var/lib/registry"]
EXPOSE 5000
ENTRYPOINT ["registry"]
CMD ["serve", "/etc/docker/registry/config.yml"]
