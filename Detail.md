

## Detail

**Items**

- distribution/v3 v3.0.0-20220620080156; google/go-containerregistry v0.4.0
  - go.mod `module github.com/distribution/distribution/v3`
  - https://github.com/distribution/distribution/issues?q=v3+in%3Atitle
  - https://github.com/distribution/distribution/pull/3225
- build
  - cert/genKey.sh[tls_san] (生成*.registry.local域名SAN证书); cfssl
  - run.sh 生成htpasswd(多用户) (admin, admin123); auth: 配置用户密码,~~动态添加用户(支持)~~ +ldap_nginx?
  - 多域名，通配域名: [泛域名(*.ssl可用; *.xx.ssl可用; *.*.ssl不可用)]
- TODO
  - ~~multiArch列表显示~~ @9.19
  - tls: tls访问>hosts_domain,导入证书用例/面板cert下载(key+cmd)\ index页认证?(open); 双端口;
  - cli: authLogin/pull

**TODO**

- ~~官方registry(cnf文件)~~> ~~+首页列表查看~~
  - ~~conf: user,pass,size~~
  - UI查看_static打包
- 简配: `docker-registry re-init`
  - auth: htpasswd生成, 多个动态生成
  - tls: 证书生成
  - ha: local/minio存储
- 简用:
  - /etc/docker/certs.d/xxDomain:port/cert.crt导入; 
  - login验证, pull拉取
  - uiDash面板(删/gabage回收/用户权限)
- 
- 单文件模式:(简配) >>conf生成|免生成?(直接用当前目录下的配置<local/s3>)
  - auth: admin帐号(htpasswd); +users/web页加用户(tokenAuth);
  - tls: 证书域名/证书生成; tpl导入cert脚本(简用);
- TODO
  - ~~size: cachedMap镜像大小~~
  - ~~store: s3配置~~
  - ~~listConf: 基于原来的定义的配置，扩展user/pass/size~~ Done. (~~参数~~及env做传入)

```bash
# TODO
简用：md2html,bindata, goTplCertsImport;
  1.cli: httpConf/tlsCertImport, authLogin, pull/push; 展示指令用法即可;
  2.ns/img:tag 删(GET:Auth)>> `curl -u AUTH ns/img:tag |原生API即可`
  3.程序内置证书compile; 远端url证书:crt一个即可;
简配：
#  1)单文件运行
  释放生成./data/{certs.domain,htpasss}
  1.环境变量/约定写死;  domain不同,生成时存于dir文件;  local/minio设定;
  2.auth:追加htpasswd(动态生效?);   list:本地sock文件通信/免accessLog;

#  2)auth.htpasswd:
  1.htpasswd文件多用户(defaultAdmin + curl/cmd:`docker-registry htpasswd user:pass:pull@match >> ./htpasswd.etxt`)
  2.token认证服务:用户/权限 (匿名用户[ns:public])

#  ~~3)auth2.token.page:~~ (暂不做，复杂度/实现成本)
  2.1 token认证服务:组/用户/权限(匿名用户[ns:public],image/tag名正则匹配); <auth.db>
  2.2 admin页:repos列表/删;
  2.3 admin页:组/用户/权限-列表、增改删; (UI选择:reg-admin/jc21/dvp-web)
```

## 附

**Ref**

- https://github.com/distribution/distribution #官方registry
- https://github.com/orcaman/concurrent-map
- https://github.com/jc21/docker-registry-ui
- https://github.com/Quiq/docker-registry-ui
- https://github.com/zebox/registry-admin #仓库管理+用户权限认证
- 
- https://github.com/google/go-containerregistry #取layer信息，计算imgSize (v1 api)
- https://github.com/opencontainers/image-spec #OCI
- https://github.com/containers/image #/v5 `syncer:v5.7.0`
- 
- https://github.com/AliyunContainerService/image-syncer
- https://gitee.com/infrastlabs/fk-image-syncer #同步工具增强：支持OCI格式镜像的同步
- https://github.com/crazy-max/undock/ #从镜像内取出文件/目录
- https://crazymax.dev/undock/usage/examples/
- :3000/g-dev2/dcp-deploy/src/branch/dev/harbor-online-v264

**Dev**

```bash
fk-distribution: docker-registry-220620; (org: v271> v284EOF; v3.x.x待出)

# refCode:
  registry/root.go CMD; copy
  registry/registry.go 改:flag二次解析/index页入口(status>imageList); TODO logrus改回;
  conf/configuration.go #cut简化; +list.param x3;
  conf/parse.go #copy


## Dev
  # run
  $ go run  ./cmd/docker-registry/*.go serve registry.yml 
  # build
  $ export CGO_ENABLED=0
  $ go build ./cmd/docker-registry
```
