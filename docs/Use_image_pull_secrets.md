## How to configure username and password for a images repository/registry
### 1. Create `imagePullSecrets`
This solution can be used in both [Dockerhub](https://hub.docker.com/) or other private/public image repositories/registries. The following example take Dockerhub as an example,
```bash
DOCKER_REGISTRY_SERVER=<URL of Dockerhub>
DOCKER_USER=<username of registry>
DOCKER_EMAIL=<email of registry>
DOCKER_PASSWORD=<password of registry>

kubectl create secret docker-registry myregistrykey \
  --docker-server=$DOCKER_REGISTRY_SERVER \
  --docker-username=$DOCKER_USER \
  --docker-password=$DOCKER_PASSWORD \
  --docker-email=$DOCKER_EMAIL
```

> Note: The URL of Dockerhub is：https://index.docker.io/v1/

### 2. Add the created secret to KubeFATE config
Make sure the secrete is created in the same namespace going to deploy FATE, add it in the `imagePullSecrets` of `cluster.yaml` as following,

```bash
imagePullSecrets: 
- name: myregistrykey
```

> Adding account information to registry can solve the problem of traffic limitation by Dockerhub