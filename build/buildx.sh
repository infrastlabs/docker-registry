#!/bin/bash
cur=$(cd "$(dirname "$0")"; pwd)
cd $cur/..

# ref _ct\fk-agent\_build.sh
buidImg=registry.cn-shenzhen.aliyuncs.com/infrastlabs/golang:1.16.8-alpine3.14 #ref: dvp-ci-mgr>> maven:3.6.3-jdk-8-slim
function buildRegistry(){

  # build-agent > pack
  # CGO_ENABLED=0
  # go build -o agent -v -ldflags "-s -w $flags" ./cmd/agent

  # ENV GO111MODULE=on
  # ENV GOPROXY=https://goproxy.cn
  # GOPATH: -v to cache
  test -d /mnt/mnt/data && pref=/mnt/data/
  cur2=$(echo $(pwd)|sed "s/_ext/dbox_ext/g") #barge253
  echo "==src: $pref$cur2"
  docker run -v $pref$cur2:/src \
    -e GO111MODULE=on -e GOPROXY=https://goproxy.cn \
    $buidImg \
    sh -c "cd /src && ls -lh && sh build/go-build.sh; ls -lh build/*.tar.gz"

}
buildRegistry

