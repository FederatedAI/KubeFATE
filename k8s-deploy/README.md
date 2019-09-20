# Deployment on Kubernetes
In a multi-node deployment scenario, a user can use [Kubernetes](https://kubernetes.io/) as their underlying infrastructure to create and manage the FATE cluster. To facilitate the deployment on Kubernetes, FATE provides scripts to generate deployment files automatically for users.

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

By default, the script pulls the images from [Docker Hub](https://hub.docker.com/search?q=federatedai&type=image) during the deployment. A user could also modify `KubeFATE/.env` to specify a registry to pull images from.

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
To verify the deployment, the user can log in the `python` pod of any party and run example cases.
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
The above example also shows that communication between two parties is working as intended, since the guest and host of the example are `party-10000` and `party-9999`, respectively.

## Custom Deployment (Optional)
By default, the Kubernetes scheduler will balance the deployment among the whole Kubernetes cluster. However, a user can also deploy a service to a specified node by using the [Node Seloctor](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector). This is useful when a service requires resource like GPU, Huge size hard disk ect. that only exists on a few machines.

View your nodes by using this command:  
`$ kubectl get nodes`
```bash
NAME      STATUS    AGE       VERSION
master    Ready     5d        v1.15.3
node-0    Ready     5d        v1.15.3
node-1    Ready     5d        v1.15.3
node-2    Ready     5d        v1.15.3
node-3    Ready     5d        v1.15.3
```

A user can also tag a specified node with lables by using this command:  
`$ kubectl label nodes <node-name> <label-key>=<label-value>`
```bash
$ kubectl label nodes node-0 fedai.org=egg0

node "node-0" labeled
```

After tagging all nodes, verify that they are worked by running:  
`$ kubectl get nodes --show-labels`
```bash
NAME      STATUS    AGE       VERSION   LABELS
master    Ready     5d        v1.15.3   kubernetes.io/arch=amd64,kubernetes.io/hostname=master,kubernetes.io/os=linux,name=master,node-role.kubernetes.io/master=
node-0    Ready     5d        v1.15.3   ..., fedai.org=egg0, ...
node-1    Ready     5d        v1.15.3   ..., fedai.org=egg1, ...
node-2    Ready     5d        v1.15.3   ..., fedai.org=worker1, ...
node-3    Ready     5d        v1.15.3   ..., fedai.org=worker2, ...
```

With the info of the node labels, a user could customize the deployment by configuring the "KubeFATE/k8s-deploy/kube.cfg". A sample is as follows:
```bash
# Specify k8s node selector, default use fedai.org
nodeLabel=fedai.org
# Please fill in multiple label value for multiple eggs, and split with spaces
eggList=(egg0 egg1 egg1) # This will deploy two egg services in node-1 and an egg module in node-0. If you only need one egg service, just fill one value.
federation=worker1
metaService=worker1
mysql=worker1
proxy=worker1
python=worker2
redis=worker2
roll=worker2
servingServer=worker2
```
The above sample will deploy `federation`, `metaService`, `mysql`, `proxy` service to node-2 and `python`, `redis`, `roll`, `servingServer` to node-3. If no value is filled in for the services will be deployed randomly among the cluster according to the strategy of the scheduler.

By default, only one egg service will be deployed. To deploy multiple egg services, please fill in the `eggList` with the label of the Kubernetes nodes (Split with spaces), the Helm will deploy one egg to each node.
