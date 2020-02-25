# php_conf_agent

基于 Kubernetes 的 Sidecar 部署模式，使用 Go 语言开发的 PHP 配置更新服务。单个`agent`支持多项目，多 `namespace`。

**支持的配置中心**

- [x] Apollo

## 使用

**编译环境**

- Go 1.3及以上

**编译命令**

`go build -o conf_agent main.go`

配置文件  `app.yaml` 说明

```yaml
clusterName: dev # 集群
type: 2 # 请求配置中心类型,1 为每30秒请求一次配置中心缓存数据;2 为实时变更推送
address: http://localhost:8080 # Apollo 服务接口地址 
ip: 10.12.1.1 # 应用部署的机器ip 这个参数是可选的，用来实现灰度发布。
configs:
-
  path: /data/www/a.example.com/config # 生成的配置文件所存放的目录 配置文件名称以 namespace 来命名
  appId: a # 项目 AppId
  namespace: # 项目中的 Namespace
    - application
    - web
-
  path: /data/www/b.example.com/config
  appId: b
  namespace:
    - application 
```

> 配置文件支持热更新，如果 type 为 1 热更新周期约为30秒；type 为 2 最长热更新周期约为60秒

**运行**

赋予执行权限

`chmod +x conf_agent`

执行

`./conf_agent`