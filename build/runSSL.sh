#!/bin/bash
cur=$(cd "$(dirname "$0")"; pwd)

certs=$cur/certs-san
domain=test.registry.ssl

export REGISTRY_HTTP_ADDR=:8143 #0.0.0.0:8143
export REGISTRY_HTTP_TLS_CERTIFICATE=$certs/$domain.crt
export REGISTRY_HTTP_TLS_KEY=$certs/$domain.key 

# mkdir -p $cur/logs
# docker run --entrypoint htpasswd registry.cn-shenzhen.aliyuncs.com/infrastlabs/registry:2.4.1 -Bbn admin admin123  > /tmp/htpasswd
echo 'admin:$2y$05$4bDXoc2Xm5DgxTLIy2eG..BA0NOxyX6ADHv5Iwj3AVtszn8W.3wE6'  > /tmp/htpasswd
# exec ./docker-registry serve $cur/cnf-registry.yml #> $cur/logs/console.log #tee #优先环境变量

cd $cur/..
go run ./cmd/registry-list/ serve $cur/cnf-registry.yml
