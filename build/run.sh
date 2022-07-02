#!/bin/bash
cur=$(cd "$(dirname "$0")"; pwd)

certs=$cur/certs; domain=test.registry.ssl
export REGISTRY_HTTP_ADDR=:8143 #0.0.0.0:8143
export REGISTRY_HTTP_TLS_CERTIFICATE=$certs/$domain.crt
export REGISTRY_HTTP_TLS_KEY=$certs/$domain.key 
# REGISTRY_AUTH
# export REGISTRY_AUTH=htpasswd
export REGISTRY_AUTH_HTPASSWD_REALM="basic-realm"
export REGISTRY_AUTH_HTPASSWD_PATH=$cur/htpasswd
export REGISTRY_STORAGE_FILESYSTEM_ROOTDIRECTORY=$cur/data/

# export REGISTRY_AUTH_HTPASSWD_PATH=/tmp/htpasswd
# docker run --entrypoint htpasswd registry.cn-shenzhen.aliyuncs.com/infrastlabs/registry:2.4.1 -Bbn admin admin123  > /tmp/htpasswd
# echo 'admin:$2y$05$4bDXoc2Xm5DgxTLIy2eG..BA0NOxyX6ADHv5Iwj3AVtszn8W.3wE6'  > /tmp/htpasswd

cd $cur
# mkdir -p $cur/logs
# go run ./cmd/docker-registry/ serve $cur/registry.yml
exec ./docker-registry serve $cur/registry.yml #> $cur/logs/console.log #tee #优先环境变量
