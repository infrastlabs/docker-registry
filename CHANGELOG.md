
**v22.12**

- 默认https服务 tls :8143
- index页列表，index页顶部新加官网API链接@22.12.27 `docker-registry-v3.0.0-221228-linux.tar.gz`
- 程序包，添加go-htpasswd

**v23.04**

- 新编译arm64版 `docker-registry-v3.0.0-230426-linux-x64/arm64.tar.gz`

**v23.12**

- index页列表，后台跑数 @23.7.23
- 查取仓库信息依赖库，`image/v5 v5.7.0 >> google/go-containerregistry v0.4.0`（tar pack 21M>6M）@23.9.14
- 代码整理
- 双端口服务，同时绑定https/http @23.9.19

**v24.08**

- Extend配置项 `extend.httpaddr, extend.list.user/pass/size` @24.1.3
- 列表页: 实时取所有tags, 只size走后台任务取数 @24.8.15
- multiArch排序: amd64, arm64, arm @24.8.19
- Fix: 列表页密码错误时，cpu100%及内存泄漏问题 @24.8.19
