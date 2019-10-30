# Deployment on Kubernetes
In a multi-node deployment scenario, a user can use [Kubernetes](https://kubernetes.io/) as their underlying infrastructure to create and manage the FATE cluster. To facilitate the deployment on Kubernetes, FATE provides scripts to generate deployment files automatically for users.

## Summary

<div style="text-align:center", align=center>
<img src="./images/k8s-summary.jpg" />
</div>

Package the FATE component into a Pod, deploy two FATE parties to two namespaces, and each party has 8 pods.
The relationship between the FATE component and the pod is as follows:

Pod            | Service URL                 | FATE component          | Expose Port
---------------|-----------------------------|-------------------------|------------
egg            | egg.\<namespace>            | egg/Storage-Service-cxx | 7888,7778,50001,50002,50003,50004
federation     | federation.\<namespace>     | federation              | 9394
meta-service   | meta-service.\<namespace>   | meta-service            | 8590
proxy          | proxy.\<namespace>          | proxy                   | 9370
roll           | roll.\<namespace>           | roll                    | 8011
redis          | redis.\<namespace>          | redis                   | 6379
serving-server | serving-server.\<namespace> | serving-server          | 8001
mysql          | mysql.\<namespace>          | mysql                   | 3306
python         | python.\<namespace>         | fate-flow/fateboard     | 9360,9380,8080

## Prerequisites
- A Linux laptop can run the installation command
- A working Kubernetes cluster.
- [The FATE Images](https://github.com/FederatedAI/FATE/tree/contributor_1.0_docker/docker-build) have been built and downloaded by nodes of Kubernetes cluster.
- Helm v2.14.0 or above installed

## Helm Introduction
The Helm is a package management tool of Kubernetes, it simplifies the deployment and management of applications on Kubernetes. Before using the script, a user needs to install on his machine first. For more details about Helm and installation please refer to the [official page](https://helm.sh/docs/using_helm/).

## Deploying FATE
Use the following command to clone repo if you did not clone before: 
```bash
$ git clone git@github.com:FederatedAI/KubeFATE.git
```

By default, the script pulls the images from [Docker Hub](https://hub.docker.com/search?q=federatedai&type=image) during the deployment.

### Use Third Party Registry (Optional)
It is recommended that non-Internet clusters use [Harbor](https://goharbor.io/) as a third-party registry. Please refer to [this guide](https://github.com/FederatedAI/KubeFATE/blob/master/registry/install_harbor.md) to install Harbor. Change the `THIRDPARTYPREFIX` to Harbor hostname in the `.env` file. `192.168.10.1` is an example of Harbor ip.

```bash
$ cd KubeFATE/
$ vi .env

THIRDPARTYPREFIX=192.168.10.1/federatedai
```

### Configure Parties
Before deployment, a user needs to define the FATE parties in `KubeFATE/k8s-deploy/kube.cfg`, a sample is as follows:
```bash
partylist=(10000 9999)
partyiplist=(proxy.fate-10000 proxy.fate-9999)
```
The above sample defines two parties, these parties will be deployed on the same Kubernetes cluster but isolated by the namespace. Moreover, each party contains one Egg service.

### Generating Deployment Files
After finished the definition, use the following command to generate deployment files:
```bash
$ cd KubeFATE/k8s-deploy/
$ bash create-helm-deploy.sh
```
If using a third-party registry, use the following command to generate deployment files:
```bash
$ cd KubeFATE/k8s-deploy/
$ bash create-helm-deploy.sh useThirdParty
```
According to the `kube.cfg`, the script creates two directories “fate-10000” and “fate-9999” under the current path. The structure of each directory is as follows:
```
fate-*
|-- templates   
|-- Chart.yaml   
|-- values.yaml
```

- The "templates" directory contains template files to deploy FATE components. 
- The "Chart.yaml" file describes the Chart's information.
- The "values.yaml" file defines the value used to render the templates.

### Launching Deployment

First make sure that the Kubernetes cluster has two namespaces, fate-9999 and fate-10000. If there is no corresponding namespace, you can create it with the following command：
```bash
$ kubectl create namespace fate-9999
$ kubectl create namespace fate-10000
```

Use the following commands to deploy parties.

- Party-10000:
```
$ helm install --name=fate-10000 --namespace=fate-10000 ./fate-10000/ 
```

- Party-9999:
```
$ helm install --name=fate-9999 --namespace=fate-9999 ./fate-9999/ 
```

After the command returns, use `helm list` to fetch the status of deployment, an example output is as follows:
```
NAME          REVISION    UPDATED                     STATUS      CHART         APP VERSION    NAMESPACE 
fate-10000    1           Tue Sep 10 10:48:47 2019    DEPLOYED    fate-0.1.0    1.0            fate-10000
fate-9999     1           Tue Sep 10 10:49:18 2019    DEPLOYED    fate-0.1.0    1.0            fate-9999 
```

In the above deployment, the data of "mysql", "redis" and "egg" will be persisted to the worker node that hosting the services(Pod). Which means if a service shifted to the other worker node, the service will be unable to read the previous data.

A simple solution to persist the data is to use a NFS as the shared storage, so that the services can read/wirte data from/to the NFS directly. An user need to setup [NFS](https://help.ubuntu.com/lts/serverguide/network-file-system.html) first, then use the following command to deploy FATE:
```
$ helm install --set nfspath=${NfsPath} --set nfsserver=${NfsIp} --name=fate-* --namespace=fate-* ./fate-*/

# NfsPath: The NFS exposed the path
# NfsIp: The NFS IP address
```

### Verifying the Deployment
To verify the deployment, the user can log in the `python` pod of his or her party and runs example cases.
The following steps illustrate how to perform a test on `party-10000`:
1. Log into the python container
```bash
$ kubectl exec -it svc/python bash -n fate-10000
```
2. Run the test toy_example
```bash
$ source /data/projects/fate/venv/bin/activate
$ cd /data/projects/fate/python/examples/toy_example/
$ python run_toy_example.py 10000 9999 1
```
3. Verify the output, a successful example is as follows:
```
"2019-08-29 07:21:34,118 - secure_add_guest.py[line:121] - INFO: success to calculate secure_sum, it is 2000.0000000000002"
```
The above example also shows that communication between two parties is working as intended, since the guest and the host of the example are `party-10000` and `party-9999`, respectively.

## Custom Deployment (Optional)
By default, the Kubernetes scheduler will balance the workload among the whole Kubernetes cluster. However, a user can deploy a service to a specified node by using [Node Selector](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector). This is useful when a service requires resources like GPU, or large size hard disk which are only available on some hosts.

View your nodes by this command:  
`$ kubectl get nodes`
```bash
NAME      STATUS    AGE       VERSION
master    Ready     5d        v1.15.3
node-0    Ready     5d        v1.15.3
node-1    Ready     5d        v1.15.3
node-2    Ready     5d        v1.15.3
node-3    Ready     5d        v1.15.3
```

A user can tag a specified node with labels, for example:  
```bash
$ kubectl label nodes node-0 fedai.hostname=egg0

node "node-0" labeled
```
The above command tagged node-0 with a label `fedai.hostname=egg`.

After tagging all nodes, verify that they are worked by running:  
`$ kubectl get nodes --show-labels`
```bash
NAME      STATUS    AGE       VERSION   LABELS
master    Ready     5d        v1.15.3   kubernetes.io/arch=amd64,kubernetes.io/hostname=master,kubernetes.io/os=linux,name=master,node-role.kubernetes.io/master=
node-0    Ready     5d        v1.15.3   ..., fedai.hostname=egg0, ...
node-1    Ready     5d        v1.15.3   ..., fedai.hostname=egg1, ...
node-2    Ready     5d        v1.15.3   ..., fedai.hostname=worker1, ...
node-3    Ready     5d        v1.15.3   ..., fedai.hostname=worker2, ...
```

With the use of node labels, a user can customize the deployment by configuring the "KubeFATE/k8s-deploy/kube.cfg". A sample is as follows:
```bash
...

# Specify k8s node selector, default use fedai.hostname
nodeLabel=fedai.hostname
# Please fill in multiple label value for multiple eggs, and split with spaces
eggList=(egg0 egg1) # This will deploy an egg service in node-0 and an egg service in node-1. If you only need one egg service, just fill one value.
federation=worker1
metaService=worker1
mysql=worker1
proxy=worker1
python=worker2
redis=worker2
roll=worker2
servingServer=worker2
```

<div style="text-align:center", align=center>
<img src="./images/k8s-cluster.jpg" />
</div>

The above sample will deploy an `egg` service in node-0 and, an `egg` service in node-1, `federation`, `metaService`, `mysql`, `proxy` services to node-2 and `python`, `redis`, `roll`, `servingServer` services to node-3. If no value is given, a service will be deployed in the cluster according to the strategy of the scheduler.

By default, only one egg service will be deployed. To deploy multiple egg services, please fill in the `eggList` with the label of the Kubernetes nodes (Separated with spaces). Helm will deploy one egg service to each node.
