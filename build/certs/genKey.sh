#!/bin/bash
cur=$(cd "$(dirname "$0")"; pwd)
DOMAIN=$1 #"test.registry.ssl"
test -z "$DOMAIN" && DOMAIN="test.registry.ssl" #exit 1

# docker 20.10+
# 问题
# Get "https://test.registry.net/v2/": x509: certificate relies on legacy Common Name field, use SANs or temporarily enable Common Name matching with GODEBUG=x509ignoreCN=0
# 分析
# 由于docker20.10.8及以上版本编译使用的go版本过高（>1.15.1）。go 1.15版本开始废弃CommonName，需要使用SAN证书

######################################################################
# 一：创建一个自建域名证书。
# 以下是registry服务端的服务器生成SAN证书的方式
# # mkdir -p $cur/certs; cd $cur/certs
# touch openssl.cnf
# cat /etc/ssl/openssl.cnf > openssl.cnf

# # 修改 操作目录的openssl.cnf
# vi $cur/certs/openssl.cnf
# # 修改如下 （去掉129 行的# 添加 两行，可以写多个域名 ）
# 129 req_extensions = v3_req # The extensions to add to a certificate request
# [ alt_names ]
# DNS.1 = test.registry.net


# 生成默认 ca
# https://blog.csdn.net/k_young1997/article/details/104425743
cd /root
openssl rand -writerand .rnd
# cd $cur/certs
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -subj "/CN=example.ca.com" -days 5000 -out ca.crt

# 生成证书请求(domain.csr)
# DOMAIN="deploy.xx.com.ssl" #VALIDATE1: 重复在ten-vm1上生成（复用ca, openssl.cnf内的DNS.1 = test.registry.net不影响）
# 指定key?: https://www.amd5.cn/atang_3828.html  #VALIDATE2: (22.6.24中午：可用)
#(umask 077; openssl genrsa 1024 >hub.key)  #创建一对1024位长度的密钥
#openssl req -new -key hub.key -out hub.csr   #生成证书颁发请求
openssl genrsa 1024 >$DOMAIN.key
csrKey=$DOMAIN.key
# csrKey=ca.key
openssl req -new -sha256 \
-key $csrKey \
-subj "/C=CN/ST=Beijing/L=Beijing/O=UnitedStack/OU=Devops/CN=$DOMAIN" \
-reqexts SAN \
-config <(cat ./openssl.cnf \
<(printf "[SAN]\nsubjectAltName=DNS:$DOMAIN")) \
-out $DOMAIN.csr

# 生成证书（这里用的ca）
DOMAIN02=$DOMAIN #*.ssl #*.*.ssl
openssl x509 -req -days 365000 \
-in $DOMAIN.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
-extfile <(printf "subjectAltName=DNS:$DOMAIN02") \
-out $DOMAIN.crt

# 当前复用的ca.key
# cp ca.key $DOMAIN.key


# +genKey.md