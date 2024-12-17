# https://gitee.com/g-golang/fk-distribution/blob/main/Dockerfile
FROM alpine:${ALPINE_VERSION}
RUN apk add --no-cache ca-certificates
# COPY cmd/registry/config-dev.yml /etc/docker/registry/config.yml
# COPY --from=binary /registry /bin/registry
VOLUME ["/var/lib/registry"]
EXPOSE 5000
ENTRYPOINT ["registry"]
CMD ["serve", "/etc/docker/registry/config.yml"]
