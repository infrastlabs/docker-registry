# docker-regsitry

基于Docker官方Registry（官方已更名`https://github.com/distribution/distribution`）扩展实现：

- 增加首页镜像列表，支持大小统计，支持MultiArch镜像
- 支持HTTP/HTTPS（Docker客户端、CURL/浏览器访问指引，详见[Detail.md](./Detail.md)）
- 二进制单文件运行，支持AMD64/ARM64/ARMV7

**文档**

- [01-auth.md](./_docs/01-auth.md) 帐号认证、TLS证书配置
- [02-store.md](./_docs/02-store.md) Minio存储(s3)、HA高可用部署
- [03-client.md](./_docs/03-client.md) 浏览器访问、Docker客户端配置

![](./_docs/res/01-imgList.jpg)

**运行**

```bash
# HTTP模式
export REGISTRY_HTTP_ADDR=:5000
./docker-registry serve ./registry.yml 

# HTTPS_TLS模式
certs=./certs; domain=registry.local
bash genKey.sh $domain #"*.$domain"
export REGISTRY_HTTP_TLS_CERTIFICATE=$certs/$domain.crt
export REGISTRY_HTTP_TLS_KEY=$certs/$domain.key 

echo 'admin:$2y$05$4bDXoc2Xm5DgxTLIy2eG..BA0NOxyX6ADHv5Iwj3AVtszn8W.3wE6'  > /tmp/htpasswd.txt #AUTH: admin admin123
export REGISTRY_HTTP_ADDR=:5000 #0.0.0.0:5000
./docker-registry serve ./registry.yml 
```

**访问**

```bash
headless @ barge in ~ |17:53:09  
$ ss -ntlp |grep regis
LISTEN0 128  *:5000   *:*users:(("docker-registry",pid=1142444,fd=3))
LISTEN0 128  *:8000   *:*users:(("docker-registry",pid=1142444,fd=8))

http://localhost:8000  #首页镜像列表 (HTTP)
https://localhost:5000 #首页镜像列表 (HTTPS)
admin admin123 #pull/push默认帐号

# Api v2:
https://localhost:5000/v2/_catalog
https://localhost:5000/v2/abc/busybox/tags/list
https://localhost:5000/v2/abc/busybox/manifests/v2

$ curl -k https://admin:admin123@localhost:5000/v2/_catalog
{"repositories":[]}

```
