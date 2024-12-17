#!/bin/bash
cur=$(cd "$(dirname "$0")"; pwd)

apk add upx #alpine-go
# ref _ct\fk-agent\_build.sh
# apt -y install upx
upx -V > /dev/null 2>&1
errCode=$?
test "0" != "$errCode" && sudo apt -y install upx

# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
#   go build -o docker-registry-x64 -v -ldflags "-s -w $flags" ./cmd/docker-registry/

# seq=220213
version="v3.0.0"
seq=$(date +%Y%m%d |sed "s/^20//g"); echo "seq: $seq"
os=linux
onePack(){
  arch=$1
  cd $cur/..
    # go build -x -v -ldflags "-s -w $flags" ./cmd/registry/*.go
    CGO_ENABLED=0 GOOS=linux GOARCH=$arch \
      go build -o build/bin/docker-registry-$arch -v -ldflags "-s -w $flags" ./cmd/docker-registry/
    echo "errCode: $?"

  rm -rf $cur/../build/docker-registry-$arch; \cp -a $cur/../build/docker-registry $cur/../build/docker-registry-$arch #cp for batchBuild
  cd $cur/../build/docker-registry-$arch
    rm -f ./docker-registry; upx -7 ../bin/docker-registry-$arch -o ./docker-registry
    # \cp -a ../docker-registry-$arch ./docker-registry;
    \cp -a ../../README.md ./; chmod +x *.sh
    tar --exclude-from=../../.tarignore -zcvf ../docker-registry-$version-$seq-$os-$arch.tar.gz *
    # clear
    # rm -f ../bin/docker-registry-$arch
    rm -rf $cur/../build/docker-registry-$arch
}
onePack arm & #TODO batchMode> -o ./docker-registry
onePack arm64 &
onePack amd64 &
wait

ls -lh $cur |grep "docker-registry-"