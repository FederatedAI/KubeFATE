## DMZ 部署

很多企业考虑网络安全都会建立DMZ区，FATE也支持这种部署方式。

![](images/Proxy-Deployment-in-DMZ.jpg)

proxy组件采用docker-compose的部署方式部署在DMZ区，其他组件通过helm部署在k8s集群。

### 在你开始之前

- kubernetes集群[v1.9+]
- DMZ区支持docker-compose
- helm

### 开始部署

#### 修改配置文件

配置文件 `KubeFATE/k8s-deploy/kube.cfg`

```bash
partylist=(10000)
DeployPartyInternal=true
proxyIpList=(192.168.13.1)
```

#### k8s部署

根据 `kube.cfg` 生成 helm chart 部署文件

```bash
$ cd KubeFATE/k8s-deploy/
$ bash create-helm-deploy.sh
```

```bash
$ helm install --name=fate-10000 --namespace=fate-10000 ./fate-10000/
```

检查是否部署成功

```bash
$ kubectl get pod -n fate-10000
```

部署成功所有的pod状态都是Running

#### DMZ部署

查看k8s 节点Ip

```bash
$ kubectl get node -o wide
NAME     STATUS   ROLES    AGE   VERSION   INTERNAL-IP      EXTERNAL-IP   OS-IMAGE                KERNEL-VERSION           CONTAINER-RUNTIME
master   Ready    master   26d   v1.16.2   192.168.12.1   <none>        CentOS Linux 7 (Core)   3.10.0-1062.el7.x86_64   docker://19.3.4
node-0   Ready    <none>   26d   v1.16.2   192.168.12.2   <none>        CentOS Linux 7 (Core)   3.10.0-1062.el7.x86_64   docker://19.3.4
```

当前集群work节点 Ip是192.168.12.2

获得federation和fateflow的nodePort

```bash
$ export Fedetation_NODE_PORT=$(kubectl get services/federation -n fate-10000 -o go-template='{{(index .spec.ports 0).nodePort}}')
$ echo Fedetation_NODE_PORT=$Fedetation_NODE_PORT
Fedetation_NODE_PORT=30792
$ export Fateflow_NODE_PORT=$(kubectl get services/fateflow -n fate-10000 -o go-template='{{(index .spec.ports 0).nodePort}}')
$ echo Fateflow_NODE_PORT=$Fateflow_NODE_PORT
Fateflow_NODE_PORT=31840
```

生成proxy部署

```bash
$ bash ../docker-deploy/generate_config.sh splitting_proxy 10000 192.168.12.2 30792 192.168.12.2 31840 192.168.13.1
Handle Splitting Proxy
Splitting proxy of 10000 done!
```

部署proxy

```bash
$ bash ../docker-deploy/docker_deploy.sh splitting_proxy 10000 192.168.13.1
...
party 10000 deploy is ok!
```

测试组件通讯正常

```bash
$ kubectl exec -it -c python svc/python -n fate-10000 -- bash
$ source /data/projects/python/venv/bin/activate
$ cd examples/toy_example/
$ python run_toy_example.py 10000 10000 1
...
"2019-11-27 06:30:17,232 - secure_add_guest.py[line:134] - INFO: success to calculate secure_sum, it is 1999.9999999999998"
```

部署成功