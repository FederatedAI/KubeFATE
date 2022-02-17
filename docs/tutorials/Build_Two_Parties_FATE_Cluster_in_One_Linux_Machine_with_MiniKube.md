# Tutorial Goal
In this tutorial, we will from scratch to install a MiniKube for Kubernetes and deploy KubeFATE service on it. Then, we will install a two-parties FATE cluster. Each of them is deployed in a given namespace. We are able to run federated learning with these two parties, and check FATE-Dashboard for the status of the learning job.

After the tutorial, the deployment architecture looks like the following diagram.

<div align="center">
  <img src="./images/goal.png">
</div>

# Prerequisites
1. A Linux machine. The verified OS is Ubuntu 18.04 LTS. <font color="red">* The demo machine is 8 core, 32G memory.</font>
2. A domain name for ingress of KubeFATE service, Jupyter Notebook, and FATE-Dashboard. An alternative is to set host both to deploying machine and client to access these endpoints. In this tutorial, we suppose to the latter case.  
3. Docker has been installed in the Linux machine. To install a Docker, please refer to [Install Docker in Ubuntu](https://docs.docker.com/install/linux/docker-ce/ubuntu/)
4. Configure username and password for a images repository/registry after the docker has been installed, please refer to [use image pull secrets](https://github.com/federatedai/KubeFATE/blob/master/docs/Use_image_pull_secrets.md).
5. Network connectivity to dockerhub or 163 Docker Image Registry, and google storage
6. Setup the global KubeFATE version using in the tutorial and create a folder for the whole tutorial. We use KubeFATE v1.6.0 in this tutorial, other versions should be similar.
```
export release_version=v1.6.0 && export kubefate_version=v1.4.1 && cd ~ && mkdir demo && cd demo
```

**<font color="red">!!!Note: in this tutorial, the IP of the machine we used is 192.168.100.123. Please change it to your machine's IP in all the following commands and config files.</font></div>**

# Start Tutorial
## Install related tools
The following tools and versions have been verified, which are the latest version by the date of drafting this tutorial.
1. MiniKube: v1.7.3
2. kubectl: v1.17.3
3. kubefate:
	* Release: v1.6.0
	* Service version: v1.4.1
	* Commandline version: v1.4.1

### Install kubectl
```
curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.17.3/bin/linux/amd64/kubectl && chmod +x ./kubectl && sudo mv ./kubectl /usr/bin
```
Try to verify if kubectl installed,
```
kubefate@machine:~/demo$ kubectl version
Client Version: version.Info{Major:"1", Minor:"17", GitVersion:"v1.17.3", GitCommit:"06ad960bfd03b39c8310aaf92d1e7c12ce618213", GitTreeState:"clean", BuildDate:"2020-02-11T18:14:22Z", GoVersion:"go1.13.6", Compiler:"gc", Platform:"linux/amd64"}
The connection to the server localhost:8080 was refused - did you specify the right host or port?
```
### Install MiniKube
```
curl -LO https://github.com/kubernetes/minikube/releases/download/v1.7.3/minikube-linux-amd64 && mv minikube-linux-amd64 minikube && chmod +x minikube && sudo mv ./minikube /usr/bin
```
Try to verify if MiniKube installed,
```
kubefate@machine:~/demo$ minikube version
minikube version: v1.7.3
commit: 436667c819c324e35d7e839f8116b968a2d0a3ff
```

### Install Kubernetes with MiniKube
In a Linux machine, we suggest using Docker as the hypervisor, which is easy. The details related to [Install MiniKube - Install a Hypervisor](https://kubernetes.io/docs/tasks/tools/install-minikube/#install-a-hypervisor). It is only one command,
```
sudo minikube start --vm-driver=none
```
Wait a few seconds for the command finish. Then run the code below to relocate kubectl or minikube.
```
sudo mv /home/vmware/.kube /home/vmware/.minikube $HOME
sudo chown -R $USER $HOME/.kube $HOME/.minikube
```
Try to verify if Kubernetes installed,
```
kubefate@machine:~/demo$ sudo minikube status
host: Running
kubelet: Running
apiserver: Running
kubeconfig: Configured
```
It means Kubernetes has been installed on your machine!

However, by default MiniKube will not enable the Ingress addon, which KubeFATE required, we need to enable it manually,
```
sudo minikube addons enable ingress
```
Till now, Kubernetes have been ready. 

### Download KubeFATE Release Pack, KubeFATE Server Image v1.4.1 and Install KubeFATE Command Lines
Go to [KubeFATE Release](https://github.com/FederatedAI/KubeFATE/releases), and find the latest kubefate-k8s release pack, which is `v1.6.0` as set to ENVs before. (replace ${release_version} with the newest version avaliable)
```
curl -LO https://github.com/FederatedAI/KubeFATE/releases/download/${release_version}/kubefate-k8s-${release_version}.tar.gz && tar -xzf ./kubefate-k8s-${release_version}.tar.gz
```
Then we will get the release pack of KubeFATE, verify it,
```
kubefate@machine:~/demo$ ls
cluster-serving.yaml  cluster.yaml  config.yaml  examples  kubefate  kubefate-k8s-v1.6.0.tar.gz  kubefate.yaml  rbac-config.yaml
```
Move the kubefate executable binary to path,
```
chmod +x ./kubefate && sudo mv ./kubefate /usr/bin
```
Try to verify if kubefate works,
```
kubefate@machine:~/demo$ kubefate version
* kubefate service connection error, Post http://example.com/v1/user/login: dial tcp: lookup example.com: no such host
* kubefate commandLine version=v1.4.1
```
It is fine only the command line version shows and get an error on KubeFATE service's version because we have not deployed KubeFATE service yet.

Then, we download the KubeFATE Server Image v1.4.1 as set to ENVs before,

${release_version} -- The release version of the KubeFATE.

${kubefate_version} -- The actual version of the KubeFATE.

You can find the newest version of the KubeFATE here: https://github.com/FederatedAI/KubeFATE/releases.

For this tutorials we are going to use ```release_version = v1.6.0 & kubefate_version = v.1.4.1```
```
curl -LO https://github.com/FederatedAI/KubeFATE/releases/download/${release_version}/kubefate-${kubefate_version}.docker
```
and load into local Docker. Please note that, because we are using MiniKube, which is an all-in-one deployment of Kubernetes, loading image to local is work for this tutorial. If you are running a cluster-installed Kubernetes, the image needs to be loaded into [Docker Registry](https://docs.docker.com/registry/introduction/) or [Harbor](https://goharbor.io/). For the details of using Harbor as a local image registry, please refer to: https://github.com/FederatedAI/KubeFATE/blob/master/registry/README.md.
```
kubefate@machine:~/demo$ docker load < ./kubefate-v1.4.1.docker
7a5b9c0b4b14: Loading layer [==================================================>]  3.031MB/3.031MB
8edfcca02080: Loading layer [==================================================>]  44.02MB/44.02MB
b7ffb386319e: Loading layer [==================================================>]  2.048kB/2.048kB
Loaded image: federatedai/kubefate:v1.4.1
```

## Deploy KubeFATE service
### Create kube-fate namespace and account for KubeFATE service
We have prepared the yaml for creating kube-fate namespace, as well as creating a service account in rbac-config.yaml in your working folder. Just apply it,
```
kubectl apply -f ./rbac-config.yaml
```

### (Optional) Use 163 Image Registory instead of Dockerhub
**Because the [Dockerhub latest limitation](https://docs.docker.com/docker-hub/download-rate-limit/), I suggest using 163 Image Repository instead.**
```
sed 's/mariadb:10/hub.c.163.com\/federatedai\/mariadb:10/g' kubefate.yaml > kubefate_163.yaml
sed 's/registry: ""/registry: "hub.c.163.com\/federatedai"/g' cluster.yaml > cluster_163.yaml
```


### Deploy KubeFATE serving to kube-fate Namespace

Apply the kubefate deployment YAML,
```
kubectl apply -f ./kubefate_163.yaml
```

We can verify it with `kubectl get all,ingress -n kube-fate`, if everything looks like,
```

kubefate@machine:~/demo$ kubectl get all,ingress -n kube-fate
NAME                            READY   STATUS    RESTARTS   AGE
pod/kubefate-5d97d65947-7hb2q   1/1     Running   0          51s
pod/mariadb-69484f8465-44dlw    1/1     Running   0          51s

NAME               TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/kubefate   ClusterIP   192.168.0.111   <none>        8080/TCP   50s
service/mariadb    ClusterIP   192.168.0.112   <none>        3306/TCP   50s

NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/kubefate   1/1     1            1           51s
deployment.apps/mariadb    1/1     1            1           51s

NAME                                  DESIRED   CURRENT   READY   AGE
replicaset.apps/kubefate-5d97d65947   1         1         1       51s
replicaset.apps/mariadb-69484f8465    1         1         1       51s

NAME                          HOSTS          ADDRESS          PORTS   AGE
ingress.extensions/kubefate   example.com   192.168.100.123   80      50s
```

It means KubeFATE service has been deployed. 

### (Optional) Add example.com to host file
Note: if we have the domain name setup, this step can be skipped.

Map the machine IP `192.168.100.123` （which is also the 'ADDRESS' field of 'ingress.extensions/kubefate'） above to `example.com`

```
sudo -- sh -c "echo \"192.168.100.123 example.com\"  >> /etc/hosts"
```

Verify if it works,
```
kubefate@machine:~/demo$ ping -c 2 example.com
PING example.com (192.168.100.123) 56(84) bytes of data.
64 bytes from example.com (192.168.100.123): icmp_seq=1 ttl=64 time=0.080 ms
64 bytes from example.com (192.168.100.123): icmp_seq=2 ttl=64 time=0.054 ms

--- example.com ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1006ms
rtt min/avg/max/mdev = 0.054/0.067/0.080/0.013 ms
```

### Verify KubeFATE service
When `example.com` well set, KubeFATE service version can be shown,
```
kubefate@machine:~/demo$ kubefate version
* kubefate service version=v1.4.1
* kubefate commandLine version=v1.4.1
```
Okay. The preparation has been done. Let's install FATE.

## Install two FATE parties: fate-9999 and fate-10000
Firstly, we need to prepare two namespaces: fate-9999 for party 9999, while fate-10000 for party 10000.
```
kubectl create namespace fate-9999
kubectl create namespace fate-10000
```
Then copy the cluster.yaml sample in the working folder. One for party 9999, the other one for party 10000,
```
cp ./cluster_163.yaml fate-9999.yaml && cp ./cluster_163.yaml fate-10000.yaml
```
They are how FATE cluster will be deployed. 

**NOTE: strongly recommend read following document**
For more what each field means, please refer to: https://github.com/FederatedAI/KubeFATE/blob/master/docs/configurations/FATE_cluster_configuration.md.

For the two files, pay extra attention of modify the partyId to the correct number otherwise you are not able to access the notebook or the fateboard.

For fate-9999.yaml, modify it as following,
```
name: fate-9999
namespace: fate-9999
chartName: fate
chartVersion: v1.6.0
partyId: 9999
registry: "hub.c.163.com/federatedai"
pullPolicy: 
persistence: false
istio:
  enabled: false
modules:
  - rollsite
  - clustermanager
  - nodemanager
  - mysql
  - python
  - fateboard
  - client

backend: eggroll

rollsite: 
  type: NodePort
  nodePort: 30091
  partyList:
  - partyId: 10000
    partyIp: 192.168.100.123
    partyPort: 30101

python:
  type: NodePort
  httpNodePort: 30097
  grpcNodePort: 30092
```

and for fate-10000:
```
name: fate-10000
namespace: fate-10000
chartName: fate
chartVersion: v1.6.0
partyId: 10000
registry: "hub.c.163.com/federatedai"
pullPolicy: 
persistence: false
istio:
  enabled: false
modules:
  - rollsite
  - clustermanager
  - nodemanager
  - mysql
  - python
  - fateboard
  - client

backend: eggroll

rollsite: 
  type: NodePort
  nodePort: 30101
  partyList:
  - partyId: 9999
    partyIp: 192.168.100.123
    partyPort: 30091

python:
  type: NodePort
  httpNodePort: 30107
  grpcNodePort: 30102
```

Okay, we can start to install these two FATE cluster via KubeFATE with the following command,
```
kubefate@machine:~/demo$ kubefate cluster install -f ./fate-9999.yaml
create job success, job id=2c1d926c-bb57-43d3-9127-8cf3fc6deb4b
kubefate@machine:~/demo$ kubefate cluster install -f ./fate-10000.yaml
create job success, job id=7752db70-e368-41fa-8827-d39411728d1b
```

There are two jobs created for deploying the FATE clusters. we can check the status of them with `kubefate job ls`. Or watch the clusters till their STATUS changing to `Running`. 
```

kubefate@machine:~/demo$ watch kubefate cluster ls
UUID                                    NAME            NAMESPACE       REVISION        STATUS  CHART   ChartVERSION    AGE
51476469-b473-4d41-b2d5-ea7241d5eac7    fate-9999       fate-9999       1               Running fate    v1.6.0          88s
dacc0549-b9fc-463f-837a-4e7316db2537    fate-10000      fate-10000      1               Running fate    v1.6.0          69s
```
We have about 10G Docker images that need to be pulled, this step will take a while for the first time. An alternative way is offline loading the images to the local environment.

To check the status of the loading, use the command,
```
kubectl get po -n fate-9999
kubectl get po -n fate-10000
```

When finished applying the image, the result will be similar to this,
```
NAME                             READY   STATUS    RESTARTS   AGE
clustermanager-bcfc6866d-nfs6c   1/1     Running   0          12m
mysql-c77b7b94b-zblt5            1/1     Running   0          12m
nodemanager-0-5599db57f4-2khcg   2/2     Running   0          12m
nodemanager-1-7c986f9454-qcscd   2/2     Running   0          12m
python-57b66d96bd-vj8kq          3/3     Running   0          12m
rollsite-7846898d6d-j2gb9        1/1     Running   0          12m
```

## Verify the deployment
### Access the cluster
From above `kubefate cluster ls` command, we know the cluster UUID of `fate-9999` is `51476469-b473-4d41-b2d5-ea7241d5eac7`, while cluster UUID of `fate-10000` is `dacc0549-b9fc-463f-837a-4e7316db2537`. Then, we can query there access information by,
```
kubefate@machine:~/demo$ kubefate cluster describe 51476469-b473-4d41-b2d5-ea7241d5eac7
UUID            51476469-b473-4d41-b2d5-ea7241d5eac7
Name            fate-9999
NameSpace       fate-9999
ChartName       fate
ChartVersion    v1.6.0
Revision        1
Age             2m22s
Status          Running
Spec            backend: eggroll
                chartName: fate
                chartVersion: v1.6.0
                istio:
                  enabled: false
                modules:
                - rollsite
                - clustermanager
                - nodemanager
                - mysql
                - python
                - fateboard
                - client
                name: fate-9999
                namespace: fate-9999
                partyId: 9999
                persistence: false
                pullPolicy: null
                python:
                  grpcNodePort: 30092
                  httpNodePort: 30097
                  type: NodePort
                registry: hub.c.163.com/federatedai
                rollsite:
                  nodePort: 30091
                  partyList:
                  - partyId: 10000
                    partyIp: 10.184.103.163
                    partyPort: 30101
                  type: NodePort

Info            dashboard:
                - party9999.notebook.example.com
                - party9999.fateboard.example.com
                ip: 10.184.103.163
                pod:
                - clustermanager-5fcbd4ccc6-fj6tq
                - mysql-7cf4d4dcb8-wvl4j
                - nodemanager-0-6cbbc86769-fk77x
                - nodemanager-1-5c6dd78f99-bgt2w
                - python-57668d4497-qwnbb
                - rollsite-f7476746-5cxh8
                status:
                  modules:
                    client: Running
                    clustermanager: Running
                    fateboard: Running
                    mysql: Running
                    nodemanager-0: Running
                    nodemanager-1: Running
                    python: Running
                    rollsite: Running
```
In `Info->dashboard` field, we can find there are two dashboards in current deployment: 
* Notebook in `party9999.notebook.example.com`, which is Jupyter Notebook integrated, where data scientists can write python or access shell in. We have pre-install FATE-clients to the Notebook.
* FATEBoard in `party9999.fateboard.example.com`, which we can inspect the status, job flows in FATE.

With similar command, we find Notebook for `fate-10000` is `party10000.notebook.example.com`, and FATEBoard for `fate-10000` is `party10000.fateboard.example.com`.

### Config dashboard's URLs in hosts
#### Note: if we have the domain name setup, this step can be skipped.

If no DNS service configured, we have to add these two url to our hosts file. In a Linux or macOS machine, 

```
sudo -- sh -c "echo \"192.168.100.123 party9999.notebook.example.com\"  >> /etc/hosts"
sudo -- sh -c "echo \"192.168.100.123 party9999.fateboard.example.com\"  >> /etc/hosts"
sudo -- sh -c "echo \"192.168.100.123 party10000.notebook.example.com\"  >> /etc/hosts"
sudo -- sh -c "echo \"192.168.100.123 party10000.fateboard.example.com\"  >> /etc/hosts"
```

In a Windows machine, you have to add them to `C:\WINDOWS\system32\drivers\etc\hosts`, please refer to [add host for Windows](https://github.com/ChrisChenSQ/KubeFATE/blob/master/docs/tutorials/Windows_add_host_tutorial.md).

### Run FATE example
If everything go well, you can access 4 dashboards now via the browser,
<div align="center">
  <img src="./images/fate-deploy-final.png" height = "500">
</div>

Click ```Pipeline/notebooks/usage_of_fate_client.ipynb```, `ipynb` is the format of Jupyter Notebook. For more info, please refer to: https://jupyter-notebook.readthedocs.io/en/stable/


Then, click on the button that showed in the image below to run the entire notebook automatically.
<div align="center">
  <img src="./images/fate-9999-run.png" height = "500">
</div>


When the notebook is running, you are able to track the process through FateBoard, 
<div align="center">
  <img src="./images/fate-9999-track-final.png" height = "500">
</div>


After the notebook finish running, if the last two lines of result shows message similar to this,
```
2021-07-07 05:31:50.784 | INFO     | pipeline.utils.invoker.job_submitter:monitor_job_status:129 - Job is success!!! Job id is 202107070529230126236
2021-07-07 05:31:50.788 | INFO     | pipeline.utils.invoker.job_submitter:monitor_job_status:130 - Total time: 0:02:24
```
This means that the job is successfully processed and KubeFate is running properly.

## Next Steps
1. The example showed above is the simplest of FATE's example. Please explore other Job examples in Notebook. But note that, the example is written in one party. You should modify it, and make sure the host party has run the `load_data` to upload data to the host party;
2. The FML_Manager will be merged to FATE-Clients soon, please check the FATE-Clients document: https://fate.readthedocs.io/en/latest/_build_temp/python/fate_client/flow_sdk/README.html. FATE-Clients has been installed and also can run in Jupyter Notebook;
3. Now you have deployed your first FATE cluster. We prepared example YAML files (https://github.com/FederatedAI/KubeFATE/tree/master/k8s-deploy/examples) for:
  * Deploy FATE-Serving
  * Deploy Spark-based FATE cluster

Check them, and we will publish more documents for above contents soon.
