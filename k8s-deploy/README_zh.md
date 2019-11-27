# kubernetes 部署 FATE

## 目标

在 kubernetes 集群上成功部署 FATE

## 在你开始之前

- 有一个 kubernetes 集群或者 minikube [v1.9+]（[如何安装 kubernetes](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/)，[如何安装 minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)）
- 已经完成了 FATE [docker 镜像的制作](../)
- 已经安装 helm（[如何安装 helm](https://helm.sh/docs/using_helm/#installing-helm)）

## 部署概要

<div style="text-align:center", align=center>
<img src="./images/k8s-summary.jpg" />
</div>

将FATE的组件封装成Pod，部署两个FATE parties到两个namespaces上，每个party有8个pod。
FATE组件和pod的关系如下：

Pod            | Service URL                 | FATE component          | Expose Port
---------------|-----------------------------|-------------------------|------------
egg            | egg.\<namespace>            | egg/Storage-Service-cxx | 7888,7778
federation     | federation.\<namespace>     | federation              | 9394
meta-service   | meta-service.\<namespace>   | meta-service            | 8590
proxy          | proxy.\<namespace>          | proxy                   | 9370
roll           | roll.\<namespace>           | roll                    | 8011
redis          | redis.\<namespace>          | redis                   | 6379
mysql          | mysql.\<namespace>          | mysql                   | 3306
python         | fateflow.\<namespace><br>fateboard.\<namespace> | fate-flow/fateboard     | 9360,9380,8080

## 准备FATE镜像

如果你可以连上Docker Hub，可以直接从Docker Hub下载FATE镜像。如果你的集群环境不能连接互联网。可以自己构建镜像或者采用离线镜像，[请参文章](https://github.com/FederatedAI/FATE/tree/master/docker-build)

## 克隆KubeFATE项目

通过以下命令克隆远端代码块：
```bash
$ git clone git@github.com:FederatedAI/KubeFATE.git
```

## 使用第三方Docker仓库
非互联网集群建议使用[Harbor](https://goharbor.io/)作为第三方仓库。安装Harbor请参考[文章](https://github.com/FederatedAI/KubeFATE/blob/master/registry/install_harbor.md)。在`.env`文件中，将`THIRDPARTYPREFIX`更改为Harbor的IP。 192.168.10.1是Harbor IP的示例。
```bash
$ cd KubeFATE/k8s-deploy/
$ vi .env

THIRDPARTYPREFIX=192.168.10.1/federatedai
```

## 修改配置文件

KubeFATE项目将大部分的配置项放在了KubeFATE/k8s-deploy/kube.cfg里面，下面是一个简单的配置：
```bash
partylist=(10000 9999)                              # partyid
partyiplist=(192.168.11.2:30010 192.168.11.3:30009) # 部署partyid的集群任一node的iP和Port
exchangeip=192.168.11.4:30000                       # 部署exchange的集群任一node的iP和Port
```

```bash
partA 实例 id 10000, namespace fate-10000
partB 实例 id 9999, namespace fate-9999
exchange 实例 id 0000, namespace fate-exchange
```
当前配置默认在同一个集群部署两个 FATE party实例，一个exchange实例，所有实例部署在不同的 namespace 上。

## 定制化部署（可选项）

在一些实际的部署中，Kubernetes集群可能有很多节点。这些节点的配置也许不同。例如有的节点有GPU，有的节点内存大等等。默认情况下，Kubernetes会自动把服务分别部署到各个节点上。如果想在特定节点上部署特定的服务，KubeFATE通过Kubernetes的 [nodeSelector](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector) 来实现这个目标，当一些资源（比如GPU，大容量存储）只能在某个节点上使用的时候，定制化部署很有用。

用这个命令查看Kubernetes节点:  
`$ kubectl get nodes -o wide`
```bash
NAME      STATUS    AGE       VERSION       INTERNAL-IP
master    Ready     5d        v1.15.3       192.168.11.1
node-0    Ready     5d        v1.15.3       192.168.11.2
node-1    Ready     5d        v1.15.3       192.168.11.3
node-2    Ready     5d        v1.15.3       192.168.11.4
node-3    Ready     5d        v1.15.3       192.168.11.5
```

选定一个节点，添加一个标签（label）：
```bash
$ kubectl label nodes node-0 fedai.hostname=fate-node-0

node "node-0" labeled
```
这个命令会给node-0添加一个 fedai.hostname=fate-node-0的标签，标签是键值对的形式，fedai.hostname是key，它的value是fate-node-0。

当所有的工作节点都添加过标签后，用下面的命令来查看：  
`$ kubectl get nodes --show-labels`
```bash
NAME      STATUS    AGE       VERSION   LABELS
master    Ready     5d        v1.15.3   kubernetes.io/arch=amd64,kubernetes.io/hostname=master,kubernetes.io/os=linux,name=master,node-role.kubernetes.io/master=
node-0    Ready     5d        v1.15.3   ..., fedai.hostname=fate-node-0, ...
node-1    Ready     5d        v1.15.3   ..., fedai.hostname=fate-node-1, ...
node-2    Ready     5d        v1.15.3   ..., fedai.hostname=fate-node-2, ...
node-3    Ready     5d        v1.15.3   ..., fedai.hostname=fate-node-3, ...
```

配置文件kube.cfg里面还有一部分配置项，可定制服务部署的节点，例如下面配置：
```bash
...

# 指定打算使用Kubernetes的label，默认用fedai.hostname
nodeLabel=fedai.hostname
# 如果想部署多个egg服务，需要添多个egg值，用空格分割
eggList=(fate-node-0 fate-node-1) #这个例子会在node-0和node-1上分别部署一个egg服务，如果只需要一个egg服务，填一个值就好了
federation=fate-node-2
metaService=fate-node-2
mysql=fate-node-2
proxy=fate-node-2
python=fate-node-3
redis=fate-node-3
roll=fate-node-3
```

在这个例子里面，会在节点node-0上面部署一个`egg`服务，节点node-1上面部署一个`egg`服务，在节点node-2上面部署`federation`,`metaService`,`mysql`,`proxy`服务，在node-3上面部署`python`,`redis`,`roll`服务。示意图如下：

<div style="text-align:center", align=center>
<img src="./images/k8s-cluster.jpg" />
</div>
如果没有给服务配置节点，这个服务会交给Kubernetes选择一个节点来部署。

## 分模块部署

DMZ部署参考文档[DMZ部署](DMZ-deploy_zh.md)

## 生成 helm 部署文件

根据 kube.cfg 生成 helm chart 部署文件
```bash
$ cd KubeFATE/k8s-deploy/
$ bash create-helm-deploy.sh
```
如果使用第三方仓库，请使用这个命令：
```bash
$ cd KubeFATE/k8s-deploy/
$ bash create-helm-deploy.sh useThirdParty
```
根据kube.cfg的内容，将会在当前目录生成两个文件夹，fate-9999/和fate-10000/。结构是这样的：
```bash
fate-*
|-- templates   
|-- Chart.yaml   
|-- values.yaml
```

- templates: 目录里面包含用来部署fate集群的Helm模版。
- Chart.yaml: 描述这个Helm chart的信息。
- values.yaml: 声明用于渲染Helm模版的变量。


## 部署
先确保Kubernetes集群有fate-9999、fate-10000和fate-exchange三个namespaces，如果没有相应的namespace，可以用下面的命令创建：
```bash
$ kubectl create namespace fate-9999
$ kubectl create namespace fate-10000
$ kubectl create namespace fate-exchange
```

执行 helm 部署命令
- Party-10000:
```
$ helm install --name=fate-10000 --namespace=fate-10000 ./fate-10000/
```

- Party-9999:
```
$ helm install --name=fate-9999 --namespace=fate-9999 ./fate-9999/
```

- Party-exchange:
```
$ helm install --name=fate-exchange --namespace=fate-exchange ./fate-exchange/
```

运行完这三个命令之后，可以用`helm list`来查看部署的状态：
```bash
NAME         	REVISION	UPDATED                 	STATUS  	CHART              	APP VERSION	NAMESPACE    
fate-10000   	1       	Tue Oct 29 03:47:05 2019	DEPLOYED	fate-party-0.3.0   	1.1      	fate-10000   
fate-9999    	1       	Tue Oct 29 03:46:58 2019	DEPLOYED	fate-party-0.3.0   	1.1      	fate-9999    
fate-exchange	1       	Tue Oct 29 03:46:53 2019	DEPLOYED	fate-exchange-0.3.0	1.1      	fate-exchange
```

在这次部署中，”MySQL”, ”Redis”, ”egg”的数据将留在服务所在的本地节点上。如果以后某个服务迁移到其他节点上了，那么以前的数据就不能用了，这是因为数据不会同步迁移。  
针对这个问题，使用NFS共享存储是一个简单的持久化方案。可以搭建一个NFS服务器，然后用下面的命令部署FATE集群。
```bash
helm install --set nfspath=${NfsPath} --set nfsserver=${NfsIp} --name=fate-* --namespace=fate-* ./fate-*/
```
需要注意的一点，NFS的路径需要no_root_squash权限。

## 验证部署

登录到名称为python的pod中跑一些例子来验证是否部署成功。
- 登录到python container
    ```bash
     $ kubectl exec -it -c python svc/fateflow bash -n fate-10000
    ```
- 运行toy_example
    ```bash
    $ source /data/projects/python/venv/bin/activate
    $ cd /data/projects/fate/python/examples/toy_example/
    $ python run_toy_example.py 10000 9999 1
    ```
- 查看输出的内容，看到这个就证明成功了
    ```bash
    "2019-09-10 07:21:34,118 - secure_add_guest.py[line:121] – INFO: success to calculate secure_sum, it 2000.0000000000002"
    ```

## 删除部署

如果需要删除部署

    $ helm del --purge fate-10000

## 可视化

如果你的k8s集群部署了ingress controller ( [ingress-nginx](https://kubernetes.github.io/ingress-nginx/deploy/) )，那么就可以通过 URL http://\<partyid\>.fateboard.fedai.org 访问可视化组件fateboard。

在这之前，你需要修改自己的hosts文件，

```bash
<node-ip> <party-id>.fateboard.fedai.org     # 增加这条记录
```

> <node-ip>：集群任一节点的IP
><party-id>：部署FATE的partyId

## 常见问题

- **Q: 所有的 pod 一直处于 ContainerCreating 状态**<br>
  A: 由于 image 体积较大，第一次下载需要一点时间。如果长时间没有改变，请检查网络。
