
- TLS使用: 证书生成| docker导入CA
- Auth认证: htpasswd(动态增删)


**auth**

```bash
# auth
# only v241, registry:2.7.1>> none;
# https://docs.docker.com/registry/configuration/
# The only supported password format is bcrypt.
docker run --entrypoint htpasswd registry:2.4.1 -Bbn admin xxx  >> htpasswd

# auth: htpasswd-test dyn-multi
$ docker run --entrypoint htpasswd registry:2.4.1 -Bbn u2 admin123
u2:$2y$05$IjbEbW0T14UlrG4ju0re1eEJataOAmJ4t8s/.vdVpMxcrFnUKUbYu
$ docker run --entrypoint htpasswd registry:2.4.1 -Bbn u3 admin456
u3:$2y$05$JxG0V7uThpZZHe/BaLQY0eVV06O.YCZx66Jq3Ub4IqwZ4iHU1CDze

docker run --entrypoint htpasswd registry.cn-shenzhen.aliyuncs.com/infrastlabs/registry:2.4.1 -Bbn admin admin123  > /tmp/htpasswd.txt
echo 'admin:$2y$05$4bDXoc2Xm5DgxTLIy2eG..BA0NOxyX6ADHv5Iwj3AVtszn8W.3wE6'  > /tmp/htpasswd.txt
```

**TLS** 泛域名

```bash
# SAN> 泛域名
# [TEST01] genKey.sh #DOMAIN02=*.ssl #$DOMAIN
# 验证yy.xx.ssl (不可用)
host-21-68:/etc/docker/certs.d/test.registry.ssl:8143 # : > cert.crt 
host-21-68:/etc/docker/certs.d/test.registry.ssl:8143 # vi cert.crt 
host-21-68:/etc/docker/certs.d/test.registry.ssl:8143 # echo admin123 |docker login test.registry.ssl:8143 --username=admin --password-stdin
Error response from daemon: Get https://test.registry.ssl:8143/v2/: x509: certificate is valid for *.ssl, not test.registry.ssl

# 验证xx.ssl (可用)
host-21-68:/etc/docker/certs.d # cp -a test.registry.ssl\:8143/ test2.ssl\:8143/
host-21-68:/etc/docker/certs.d # cat /etc/hosts
172.25.23.199 test.registry.ssl
172.25.23.199 test2.ssl
host-21-68:/etc/docker/certs.d # echo admin123 |docker login test2.ssl:8143 --username=admin --password-stdin
WARNING! Your password will be stored unencrypted in /root/.docker/config.json.
Login Succeeded


# [TEST02] genKey.sh #DOMAIN02=*.*.ssl #$DOMAIN
# 验证yy.xx.ssl (不可用)
host-21-68:../test.registry.ssl:8143 # echo admin123 |docker login test.registry.ssl:8143 --username=admin --password-stdin
Error response from daemon: Get https://test.registry.ssl:8143/v2/: x509: certificate is valid for *.*.ssl, not test.registry.ssl

# 验证xx.ssl (不可用)
host-21-68:/etc/docker/certs.d/test2.ssl:8143 # echo admin123 |docker login test2.ssl:8143 --username=admin --password-stdin
Error response from daemon: Get https://test2.ssl:8143/v2/: x509: certificate is valid for *.*.ssl, not test2.ssl
```
