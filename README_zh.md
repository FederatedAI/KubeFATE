原文链接: [《KubeFATE README》](https://github.com/FederatedAI/KubeFATE/blob/master/README.md)

## 总览

由于联合学习需要多方协作来训练机器学习模型，因此，它通常基于分布式系统来实现。 KubeteFATE 使用云技术（例如容器）来分配联合学习工作负载
。 KubeFATE 的出现， 使联合学习作业能够在公有环境，私有环境和混合环境中进行。

目前，KubeFATE 支持通过 Docker Compose 和 Kubernetes 来部署 [FATE](https://github.com/FederatedAI/FATE)。

### FATE 组件的容器镜像

FATE 发布的所有组件都被提前构建在 Docker 镜像中。您可以直接从 Docker Hub 中提取它们。容器镜像是安装 FATE 的首选方法，它相比起从源代码构建 FATE 可以节省大量时间。

[Harbor](https://github.com/goharbor/harbor) 可以作为本地注册表来存储和提供 FATE 镜像，它可以从 Docker Hub 复制容器镜像以便于在本地环境中运行。 Harbor 可显着提高性能并降低对网络的依赖，因此推荐在容器中使用它。

请查看 [FATE 镜像的构建](https://github.com/FederatedAI/FATE/tree/master/docker-build) 来学习如何使用源代码构建 FATE 组件的镜像。要为开发环境设置 Harbor 注册表，请查看[这篇文章](https://github.com/FederatedAI/KubeFATE/blob/master/registry/README.md)。

### 使用 Docker Compose 部署

Docker Compose 可以将 FATE 组件部署在一台主机上。通过 Docker Compose， 您可以为使用联合方式进行协作的多方环境部署 FATE。您可以通过 [Docker Compose 部署](https://github.com/FederatedAI/KubeFATE/blob/master/docker-deploy/README.md) 来查阅详情。

### 使用 Kubernetes 部署

To deploy FATE in the cloud or in a multi-node environment, a convenient way is to use a Kubernetes cluster as the underlying infrastructure. Helm Charts can be used to deploy FATE on Kubernetes. Please refer to Kubernetes Deployment for more details.

如果您想在云环境或多节点环境中部署 FATE，使用 Kubernetes 集群作为基础架构是一种相当便捷的方法。 Helm Charts 可用于在 Kubernetes 上部署 FATE。您可以通过 [Kubernetes 部署](https://github.com/FederatedAI/KubeFATE/blob/master/k8s-deploy/README.md) 来查阅详情。

### ".env" 的使用须知

默认情况下，安装脚本在部署过程中从 [Docker Hub](https://hub.docker.com/u/federatedai) 提取映像。此外，用户还可以修改 `.env` 以指定本地注册表（例如 Harbor ）以从中提取镜像。

### 许可证

[Apache License 2.0](https://github.com/FederatedAI/FATE/blob/master/LICENSE)
