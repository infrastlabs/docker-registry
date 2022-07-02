#!/bin/bash
cur=$(cd "$(dirname "$0")"; pwd)
cd $cur

#echo "export DOCKER_REGISTRY_USER_sdsir=xxx" >> /etc/profile
#echo "export DOCKER_REGISTRY_PW_sdsir=xxx" >> /etc/profile

export GO111MODULE=on
export GOPROXY=https://goproxy.cn
# -x -v
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -v -ldflags "-s -w $flags" ./cmd/docker-registry/

source /etc/profile
export |grep DOCKER_REG
repo=registry.cn-shenzhen.aliyuncs.com
echo "${DOCKER_REGISTRY_PW_infrastSubUser2}" |docker login --username=${DOCKER_REGISTRY_USER_infrastSubUser2} --password-stdin $repo

ns=infrastlabs
# cache="--no-cache"
# pull="--pull"
ver=latest #02: +full; 04: bins;
img="docker-registry:$ver"

docker build $cache $pull -t $repo/$ns/$img -f Dockerfile . 
docker push $repo/$ns/$img
