#!/bin/sh
# cur=$(cd "$(readlink $(dirname "$0"))"; pwd)
cur=$(dirname $(readlink -f "$0"))

# AUTH
# export REGISTRY_AUTH=htpasswd
test -z "$REGISTRY_AUTH_HTPASSWD_PATH" && export REGISTRY_AUTH_HTPASSWD_PATH=$cur/htpasswd.txt
test -z "$REGISTRY_AUTH_HTPASSWD_REALM" && export REGISTRY_AUTH_HTPASSWD_REALM="basic-realm"
test -z "$REGISTRY_STORAGE_FILESYSTEM_ROOTDIRECTORY" && export REGISTRY_STORAGE_FILESYSTEM_ROOTDIRECTORY=$cur/data/

# HTTPS/HTTP
# export REGISTRY_HTTP_ADDR=:8143 #改绑定端口(默认5000); 0.0.0.0:8143 
domain=registry.local #test.registry.ssl
test -z "$REGISTRY_HTTP_TLS_CERTIFICATE" && export REGISTRY_HTTP_TLS_CERTIFICATE=$cur/certs/$domain.crt
test -z "$REGISTRY_HTTP_TLS_KEY" && export REGISTRY_HTTP_TLS_KEY=$cur/certs/$domain.key 

# EXTEND.HTTP,LIST
# export REGISTRY_EXTEND_HTTPADDR=:9000 #REGISTRY_HTTP_ADDR开启tls模式时，额外开启http端口, ":0"则不开启
# export REGISTRY_EXTEND_LIST_USER=admin
# export REGISTRY_EXTEND_LIST_PASS=admin123
# export REGISTRY_EXTEND_LIST_SIZE=false

cd $cur
# go run ./cmd/docker-registry/ serve $cur/registry.yml
args="$@"
test -z "$args" && args="serve $cur/registry.yml"
exec ./docker-registry $args
