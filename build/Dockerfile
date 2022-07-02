
FROM alpine:3.12

RUN domain="mirrors.aliyun.com" \
  && echo "http://$domain/alpine/v3.12/main" > /etc/apk/repositories \
  && echo "http://$domain/alpine/v3.12/community" >> /etc/apk/repositories

RUN set -ex \
    && apk add --no-cache ca-certificates

# COPY ./htpasswd /htpasswd
COPY ./build/registry.yml /etc/docker/registry/config.yml
COPY ./docker-registry /bin/docker-registry
COPY ./build/htpasswd /tmp/htpasswd
VOLUME ["/var/lib/registry"]
EXPOSE 5000
ENTRYPOINT ["/bin/docker-registry"]
CMD ["serve", "/etc/docker/registry/config.yml"]
