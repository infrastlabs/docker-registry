#!/bin/bash
cur=$(cd "$(dirname "$0")"; pwd)
DOMAIN=$1 #"test.registry.ssl"
test -z "$DOMAIN" && DOMAIN="registry.local" #test.registry.ssl
keylen=2048 #1024
days=365000 #1k年; (cfssl: <300年)

# docker 20.10+
# 问题
# Get "https://test.registry.net/v2/": x509: certificate relies on legacy Common Name field, use SANs or temporarily enable Common Name matching with GODEBUG=x509ignoreCN=0
# 分析
# 由于docker20.10.8及以上版本编译使用的go版本过高（>1.15.1）。go 1.15版本开始废弃CommonName，需要使用SAN证书

# ref:
#   //arch-draft//2022\registry\genKey\tls_sanKey.md; tls_rsaKey.sh >> tls_sanKey.sh
#   ref1: https://blog.csdn.net/jeccisnd/article/details/106896440 ##Docker Registry 支持自建证书的Https访问 ##(非san证书)
#   ref2: https://www.jianshu.com/p/cad3377692c9  @2024 ##openssl为IP签发证书（支持多IP/内外网）
#   ref3: https://zhuanlan.zhihu.com/p/36981565 浅谈SSL/TLS工作原理 @2018
######################################################################
# 一：创建一个自建域名证书。
echo -e "\n###[openssl.cnf]############################################"
    # 01.以下是registry服务端的服务器生成SAN证书的方式（手动）
        # # mkdir -p $cur/certs; cd $cur/certs
        # touch openssl.cnf
        # cat /etc/ssl/openssl.cnf > openssl.cnf

        # # 修改 操作目录的openssl.cnf
        # vi $cur/openssl.cnf
        # # 修改如下 （去掉129 行的# 添加 两行，可以写多个域名 ）
        # 129 req_extensions = v3_req # The extensions to add to a certificate request
        # [ alt_names ]
        # DNS.1 = test.registry.net
    # 02.cmd: https://www.jianshu.com/p/cad3377692c9

cat > $cur/openssl.cnf <<EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req

[req_distinguished_name]
countryName = Country Name (2 letter code)
countryName_default = CH
stateOrProvinceName = State or Province Name (full name)
stateOrProvinceName_default = GD
localityName = Locality Name (eg, city)
localityName_default = GZ
organizationalUnitName  = Organizational Unit Name (eg, section)
organizationalUnitName_default  = ORGZ
commonName = Internet Widgits Ltd
commonName_max  = 64

[ v3_req ]
# Extensions to add to a certificate request
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
# 改成自己的域名
#DNS.1 = *.registry.local
#DNS.2 = helpdesk.example.org
#DNS.3 = systems.example.net

# 改成自己的ip
IP.1 = 172.16.24.143
IP.2 = 172.29.50.253
EOF



echo -e "\n###[CA: pri, pub]############################################"
# 生成默认 ca
# https://blog.csdn.net/k_young1997/article/details/104425743
# cd /root
#     openssl rand -writerand .rnd #拷贝版openssl.cnf>> RANDFILE		= $ENV::HOME/.rnd
cd $cur
    test -s ca.key || openssl genrsa -out ca.key 2048
    test -s ca.crt || openssl req -x509 -new -nodes -key ca.key -subj "/CN=private.ca.com" -days $days -out ca.crt

echo -e "\n###[server: pri]############################################"
openssl genrsa $keylen >$DOMAIN.key


echo -e "\n###[server: csr]############################################"
# 生成证书请求(domain.csr)
# DOMAIN="deploy.xx.com.ssl" #VALIDATE1: 重复在ten-vm1上生成（复用ca, openssl.cnf内的DNS.1 = test.registry.net不影响）
# 指定key?: https://www.amd5.cn/atang_3828.html  #VALIDATE2: (22.6.24中午：可用)
#(umask 077; openssl genrsa 1024 >hub.key)  #创建一对1024位长度的密钥
#openssl req -new -key hub.key -out hub.csr   #生成证书颁发请求
csrKey=$DOMAIN.key
# csrKey=ca.key
openssl req -new -sha256 \
    -key $csrKey \
    -subj "/C=ZH/ST=GD/L=GZ/O=ORGZ/OU=Devops/CN=$DOMAIN" \
    -reqexts SAN \
    -config <(cat ./openssl.cnf \
    <(printf "[SAN]\nsubjectAltName=DNS:$DOMAIN")) \
    -out $DOMAIN.csr

echo -e "\n###[server: pubSign]############################################"
# 生成证书（这里用的ca）
DOMAIN02="*.$DOMAIN" #*.ssl #*.*.ssl
openssl x509 -req -days $days \
    -in $DOMAIN.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -extfile <(printf "subjectAltName=DNS.1:$DOMAIN02,DNS.2:localhost,IP.1:172.29.50.253,IP.2:172.29.50.254") \
    -out $DOMAIN.crt

# 当前复用的ca.key
# cp ca.key $DOMAIN.key


# +genKey.md