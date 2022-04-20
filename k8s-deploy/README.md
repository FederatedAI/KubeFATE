# Kubernetes Deployment
We recommend using [Kubernetes](https://kubernetes.io/) as an underlying infrastructure to create and manage the FATE clusters in a production environment. KubeFATE supports deploying multiple FATE clusters in an instance of Kubernetes with different namespaces for the purposes of development, testing and production. Considering the different IT designs and standards in each company, the actual deployment should be customized. KubeFATE is flexibile for the FATE configuration.

If you focus on how to quickly use KubeFATE, please jump to [Use Scenarios](#use-scenarios).

## High-level architecture of multiple federated learning parties
The high-level architecture of a multi-party federated learning deployment (e.g. two parties) is shown as follows:
<div align="center">
  <img src="./images/hamflp.PNG">
</div>

* KubeFATE: Orchestrates a FATE cluster of a party. It offers APIs for FATE-Cloud Manager and other management portals.
* Harbor (Optional): Versioned FATE deployments and images management.
* Kubernetes: Container orchestration engine.

KubeFATE is responsible for:
* Day 1 initialization: Provision a FATE cluster on Kubernetes
* Day 2 operations: Provides RESTful APIs to manage FATE clusters

## High-level architecture of KubeFATE
The high-level architecture of KubeFATE is shwon as follows:
<div align="center">
  <img src="./images/kfha.PNG">
</div>

The numbers depicted in the diagram:
1. Accepting external API calls of Authentication & authorization
2. Rendering templates via Helm;
3. Storing jobs and configuration of a FATE deployment
4. KubeFATE is running as a service of Kubernetes

There are two parts of KubeFATE:
* The KubeFATE CLI. KubeFATE CLI is an executable helps to quickly initialize and manage a FATE cluster in an interactive mode. It does not rely on Kubernetes. Eventually, KubeFATE CLI calls KubeFATE Service for operations with a KubeFATE user token.
* The KubeFATE Service. The KubeFATE service provides RESTful APIs for managing FATE clusters. The KubeFATE service is deployed in Kubernetes, and exposes APIs via [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/). For the authentication and authorization, the KubeFATE service implements [JWT](https://jwt.io/introduction/), and neutral to other security solutions which can be added to Kubernetes ingress.

KubeFATE is designed to handle different versions FATE. Normally, KubeFATE CLI and KubeFATE service can work with several FATE releases.

## User scenarios
Suppose in an organization, there are two roles:
* System Admin: who is responsible for the infrastructure management as well as Kubernetes administration
* ML Infrastructure Operators: who is responsible for managing the machine learning cluster like FATE

<div align="center">
  <img src="./images/swim.PNG">
</div>
### Initializing a FATE deployment

Recommended version of dependent software:

Kubernetes: [v1.23.5](https://github.com/kubernetes/kubernetes/releases/tag/v1.23.5)

Ingress-nginx: [v1.1.3](https://github.com/kubernetes/ingress-nginx/releases/tag/controller-v1.1.3)

#### Creating role, namespace and other resource in Kubernetes
The example yaml can be found in [rbac-config.yaml](./rbac-config.yaml). In this example, we create a kube-fate namespace for KubeFATE service. Resource constraints can be applied to kube-fate namespace, refer to [Kubernetes Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/), [Configure Memory and CPU Quotas for Namespace](https://kubernetes.io/docs/tasks/administer-cluster/manage-resources/quota-memory-cpu-namespace/).

Run the following command to create the namespace:
```
$ kubectl apply -f ./rbac-config.yaml
```
Note that, the default username and password of KubeFATE service can be set in `rbac-config.yaml` Secret->kubefate-secret->stringData :

```
stringData:
  kubefateUsername: admin
  kubefatePassword: admin
```

#### Preparing domain name and deploying KubeFATE in Kubernetes
Because KubeFATE service exposes RESTful APIs for external access, system admin needs to prepare a domain name for KubeFATE service. In our example, the domain name is `example.com` . Moreover, system admin should create a namespace (e.g. fate-9999) for FATE deployment.
```
$ kubectl apply -f ./kubefate.yaml
$ kubectl create namespace fate-9999
```
For more about the configuration of KubeFATE service, please refer to: [KubeFATE service Configuration Guild](../docs/configurations/kubefate_service_configuration.md).

#### Preparing cluster configuration and deploying FATE
After the system admin deployed the KubeFATE service and prepared the namespace for FATE. The ML Infrastructure Operator is able to start the deployment of FATE. The `config.yaml` for `kubefate` CLI is required. It contains the username and password of KubeFATE access, and the KubeFATE service URL:

```
log:
  level: info
user:
  username: admin
  password: admin

serviceurl: example.com
```

|Name       |Type    |Description                                                       |
|-----------|--------|------------------------------------------------------------------|
|log        |scalars |The log level of command line.                                    |
|user       |mappings|User name and password when logging into KubeFATE service.        |
|serviceurl |scalars |KubeFATE service's ingress domain name, defined in kubefate.yaml. |

Create a `cluster.yaml` for FATE cluster configuration. The details of configuration can be found here: [FATE Cluster Configuration Guide](../docs/configurations/FATE_cluster_configuration.md). 

**NOTE:** For Chinese user, specifying a local image registry in `cluster.yaml` can accelerate the download of images. The details are as follows:
```
registry: "hub.c.163.com/federatedai"
```

Next, install the FATE cluster,

```
$ kubefate cluster install -f ./cluster.yaml
create job success, job id=d92d7a56-7002-46a4-9363-da9c7346e05a
```
*NOTE: If you want to deploy **FATE on Spark**, you can use `cluster-spark.yaml`.*

#### Checking the status of "Installing Cluster" job
After the above command has finished, a job is created for installing a FATE cluster. Run the command `kubefate job describe` to check the status of the job, until the "Status" turns to `Success`.

```bash
$ kubefate job describe d92d7a56-7002-46a4-9363-da9c7346e05a
UUID     	d92d7a56-7002-46a4-9363-da9c7346e05a
StartTime	2022-04-12 07:34:09
EndTime  	2022-04-12 07:48:14
Duration 	14m
Status   	Success
Creator  	admin
ClusterId	24bb75ff-f636-4c64-8c04-1b9073f89a2f
States   	- update job status to Running
         	- create Cluster in DB Success
         	- helm install Success
         	- checkout Cluster status [794]
         	- job run Success

SubJobs  	nodemanager-0        ModuleStatus: Available, SubJobStatus: Success, Duration:    13m, StartTime:
         	2022-04-12 07:34:09, EndTime: 2022-04-12 07:47:26
         	nodemanager-1        ModuleStatus: Available, SubJobStatus: Success, Duration:    13m, StartTime:
         	2022-04-12 07:34:09, EndTime: 2022-04-12 07:47:18
         	python               ModuleStatus: Available, SubJobStatus: Success, Duration:    14m, StartTime:
         	2022-04-12 07:34:09, EndTime: 2022-04-12 07:48:14
         	rollsite             ModuleStatus: Available, SubJobStatus: Success, Duration:    13m, StartTime:
         	2022-04-12 07:34:09, EndTime: 2022-04-12 07:47:24
         	client               ModuleStatus: Available, SubJobStatus: Success, Duration:    11m, StartTime:
         	2022-04-12 07:34:09, EndTime: 2022-04-12 07:45:22
         	clustermanager       ModuleStatus: Available, SubJobStatus: Success, Duration:    13m, StartTime:
         	2022-04-12 07:34:09, EndTime: 2022-04-12 07:47:11
         	mysql                ModuleStatus: Available, SubJobStatus: Success, Duration:    13m, StartTime:
         	2022-04-12 07:34:09, EndTime: 2022-04-12 07:47:11
```
#### Describing the cluster and finding FATE access information
After the `installing cluster` job succeeded, use `kubefate cluster describe` to check the FATE access information:
```bash
$ kubefate cluster describe 24bb75ff-f636-4c64-8c04-1b9073f89a2f
UUID        	24bb75ff-f636-4c64-8c04-1b9073f89a2f
Name        	fate-9999
NameSpace   	fate-9999
ChartName   	fate
ChartVersion	v1.8.0
Revision    	1
Age         	44h
Status      	Running
Spec        	backend: eggroll
            	chartName: fate
            	chartVersion: v1.8.0
            	imagePullSecrets:
            	- name: myregistrykey
            	imageTag: 1.8.0-release
            	ingress:
            	  client:
            	    hosts:
            	    - name: party9999.notebook.example.com
            	  fateboard:
            	    hosts:
            	    - name: party9999.fateboard.example.com
            	ingressClassName: nginx
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
            	podSecurityPolicy:
            	  enabled: false
            	pullPolicy: null
            	python:
            	  grpcNodePort: 30092
            	  httpNodePort: 30097
            	  logLevel: INFO
            	  type: NodePort
            	registry: ""
            	rollsite:
            	  nodePort: 30091
            	  partyList:
            	  - partyId: 10000
            	    partyIp: 192.168.10.1
            	    partyPort: 30101
            	  type: NodePort
            	servingIp: 192.168.9.2
            	servingPort: 30095

Info        	dashboard:
            	- party9999.notebook.example.com
            	- party9999.fateboard.example.com
            	ip: 192.168.9.1
            	port: 30091
            	status:
            	  containers:
            	    client: Running
            	    clustermanager: Running
            	    fateboard: Running
            	    mysql: Running
            	    nodemanager-0: Running
            	    nodemanager-0-eggrollpair: Running
            	    nodemanager-1: Running
            	    nodemanager-1-eggrollpair: Running
            	    python: Running
            	    rollsite: Running
            	  deployments:
            	    client: Available
            	    clustermanager: Available
            	    mysql: Available
            	    nodemanager-0: Available
            	    nodemanager-1: Available
            	    python: Available
            	    rollsite: Available
```

#### Access the UI of FATEBoard and Notebook

If the components of fateboard and client are installed, you can use the information `party9999.fateboard.example.com` and `party9999.notebook.example.com` obtained in the previous step to access FATEBoard and Notebook UI, and configure the resolution of these two domain names It can be opened in the browser.

##### FATEBoard

 http://party9999.fateboard.example.com

Access to FATEBoard UI requires a login user name and password, which can be found in `cluster.yaml` [Configuration](../docs/configurations/FATE_cluster_configuration.md#fateboard mappings).

![fate_board](../docs/tutorials/images/tkg_fate_board.png)

##### Notebook

 http://party9999.fateboard.example.com

![notebook](../docs/tutorials/images/tkg_notebook.png)

### Other user scenarios
#### [Manage FATE and FATE-Serving Version](../docs/Manage_FATE_and_FATE-Serving_Version.md)
#### [Update and Delete a FATE Cluster](../docs/Update_and_Delete_a_FATE_Cluster.md)
#### [KubeFATE Examples](examples)

#### [KubeFATE Command Line User Guide](../docs/KubeFATE_command_line_user_guide.md)

## KubeFATE service RESTful APIs reference
#### [API Reference](docs/KubeFATE_API_Reference_Swagger.md)
