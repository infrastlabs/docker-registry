#!/bin/bash
cur=$(cd "$(dirname "$0")"; pwd)

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
      go build -o build/docker-registry-$arch -v -ldflags "-s -w $flags" ./cmd/docker-registry/

  cd $cur/../build/docker-registry
    # upx --help > /dev/null 2>&1 && echo "upx installed, skip" || sudo apt install upx
    # upx -7 ./psu -o /tmp/psu #段错误

    # cd /tmp; tar -zcf $cur/psu-$version-$seq-$os.tar.gz
    \cp -a ../docker-registry-$arch ./docker-registry;
    \cp -a ../../README.md ./; chmod +x *.sh
    tar --exclude-from=../../.tarignore -zcvf ../docker-registry-$version-$seq-$os-$arch.tar.gz *
}
onePack arm
onePack arm64
onePack amd64

ls -lh $cur |grep "docker-registry"