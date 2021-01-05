## 如何配置使用私有镜像仓库部署FATE

### 1. 生成`imagePullSecrets`

生成docker hub的Secrets

```bash
DOCKER_REGISTRY_SERVER=<registry的URL>
DOCKER_USER=<你的docker用户名>
DOCKER_EMAIL=<你的docker邮箱>
DOCKER_PASSWORD=<你的docker密码>

kubectl create secret docker-registry myregistrykey \
  --docker-server=$DOCKER_REGISTRY_SERVER \
  --docker-username=$DOCKER_USER \
  --docker-password=$DOCKER_PASSWORD \
  --docker-email=$DOCKER_EMAIL
```

> docker hub的registry的URL：https://index.docker.io/v1/

### 2. 配置

在要部署FATE的对应namespace下生成上述secret，然后将secret的name写入`cluster.yaml`的`imagePullSecrets`

```bash
# 例如这样
imagePullSecrets: 
- name: myregistrykey
```

配合`registry`可以使用私有的仓库镜像，也可以应对docker hub限流。