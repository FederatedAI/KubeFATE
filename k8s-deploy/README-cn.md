# kubernetes 部署 FATE

## 目标

在 kubernetes 集群上成功部署 FATE

## 在你开始之前

- 有一个 kubernetes 集群或者 minikube（[如何安装 kubernetes](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/)，[如何安装 minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)）
- 已经完成了 FATE [docker 镜像的制作](../)
- 已经安装 helm（[如何安装 helm](https://helm.sh/docs/using_helm/#installing-helm)）

## 修改配置文件 kube.cfg

    partylist=(10000 9999)
    partyiplist=(proxy.fate-10000 proxy.fate-9999)

> partA 实例 id 10000, namespace fate-10000<br>
> partB 实例 id 9999, namespace fate-9999

当前配置默认在同一个集群部署两个 FATE 实例，两个实例部署在不同的 namespace 上。

## 定制化部署（可选项）
通过填写 `nodeSelector` 指定某个模块安装在某个 kubernetes 节点上面，默认使用 hostname 作为 `label-key`。  
通过运行 `kubectl get nodes --show-labels` 来查看所有node的label。也可以运行 `kubectl describe node "nodename"` 来查看某个node的所有 `label` 。  
如果想使用自定义的 `label` ，运行 `kubectl label nodes <node-name> <label-key>=<label-value>` 给node增加新的 `label` 。  更多信息请参考[nodeSelector](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector)。

默认只部署一个egg节点。如果想部署多个egg节点，请填写多个与 egg 对应的 `nodeSelector` （用空格进行分割）。  

    # Specify k8s node selector, default use hostname
    # Label key
    nodeLabel=kubernetes.io/hostname
    # Please fill in multiple hostname for multiple eggs
    # Label value
    eggList=()
    federation=
    metaService=
    mysql=
    proxy=
    python=
    redis=
    roll=
    servingServer=

## 生成 helm 部署文件

根据 kube.cfg 生成 helm chart 部署文件
`bash create-helm-deploy.sh`
在当前目录生成若干个`fate-*`的目录，例如：

    drwxr-xr-x. 2 root root   60 9月   9 18:21 fate-10000
    drwxr-xr-x. 2 root root   60 9月   9 18:21 fate-9999

## 部署

执行 helm 部署命令

partA:

    $ helm install --name=fate-10000 --namespace=fate-10000 ./fate-10000/

partB:

    $ helm install --name=fate-9999 --namespace=fate-9999 ./fate-9999/

执行上述命令后如果出现

    NAME: fate-10000
    LAST DEPLOYED: Mon Sep 9 18:50:49 2019
    NAMESPACE: fate-10000
    STATUS: DEPLOYED

说明部署成功。

## 持久化

如果需要永久化 MySQL、redis 和 egg 的数据,
你需要一个 nfs 服务。

安装 nfs 服务,服务端代表存储数据的节点，客户端是所有 kubernetes 节点。

    # 服务端
    $ yum install -y nfs-utils rpcbind
    # 客户端
    $ yum install -y nfs-utils

    # 服务端
    # 创建共享目录
    $ mkdir -p /data/fate-data
    $ chmod 755 /data/fate-data
    # 修改 NFS 配置文件 `/etc/exports`
    $ vim /etc/exports
    /data/fate-data *(rw,sync,insecure,no_subtree_check,no_root_squash)

    # 启动 RPC 服务
    $ systemctl start rpcbind

    # 启动 NFS 服务
    $ systemctl start nfs

    # 客户端
    $ showmount -e 192.168.0.2
    Export list for 192.168.0.2:
    /data/fate-data *

    $ NfsPath=/data/fate-data
    $ NfsIp=192.168.0.2

部署可以持久化的 FATE 实例

    $ helm install --set nfspath=${NfsPath} --set nfsserver=${NfsIp} --name=fate-10000 --namespace=fate-10000 ./fate-10000/
然后就可以启动 MySQL 持久化的 FATE
提示：当前持久化是 MySQL 和 Redis 的数据存储在 nfs 服务器，egg 数据存储在 pod 所在节点的本地磁盘上。
如果对部署有更进一步的需求可以手动修改 values.yaml 配置文件

## 测试部署成功

查看 namespace 部署

    $ kubectl get all -n fate-10000

进入 python 节点

    $ kubectl exec -it svc/python bash -n fate-10000

当只部署 partA 时执行测试

    $ source /data/projects/fate/venv/bin/activate
    $ cd /data/projects/fate/python/examples/toy_example/
    $ python run_toy_example.py 10000 10000 1

如果同时部署 partA partB

    $ source /data/projects/fate/venv/bin/activate
    $ cd /data/projects/fate/python/examples/toy_example/
    $ python run_toy_example.py 10000 9999 1

如果没有返回错误，说明 FATE 实例已经成功部署。

## 删除部署

如果需要删除部署

    $ helm del --purge fate-10000

## 常见问题

- **Q: python 的 pod 状态一直不是 Running。**<br>
  A: `pod/python-\*` 的运行会依赖 MySQL。MySQL 第一次启动会有初始化数据库的过程，必须等待 MySQL 服务正常运行。
  等待一会即可。

- **Q: 所有的 pod 一直处于 ContainerCreating 状态**<br>
  A: 由于 image 体积较大，第一次下载需要一点时间。如果长时间没有改变，请检查网络。
