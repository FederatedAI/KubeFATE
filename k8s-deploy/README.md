# Kubernetes Deployment
We recommend use [Kubernetes](https://kubernetes.io/) as a underlying infrastructure to create and manage the FATE clusters in production environment. KubeFATE supports deploying multiple FATE clusters in one Kubernetes with different namespaces for development, test and production cases. Considering the various IT designed and standards in each company, the modules deployed should be customized. KubeFATE is isolated from the detail FATE configurations.

If you focus on how to quickly use KubeFATE, please jump to [Use Scenarios](#use-scenarios) section.

## Highlevel Architecture of multiple Federated Learning Parties
The very highlevel architecture of a multiple Federated Learning deployment (e.g. two parties) as follow image:
<div align="center">
  <img src="./images/hamflp.PNG">
</div>

* KubeFATE: Orchestrated FATE cluster inside one party, offer APIs for FATE-Cloud Manager and other management portals
* Harbor (Optional): Versioned FATE deployments and images management
* Kubernetes: Orchestration engine.

KubeFATE will responsible for:
* Day 1 initialization: One executable binary to deploy a FATE cluster
* Day 2 operations: Provides both executable binary and RESTful APIs to manage FATE clusters inside a party

## Highlevel Architecture of KubeFATE
The highlevel architecture of KubeFATE can be presented as follow image:
<div align="center">
  <img src="./images/kfha.PNG">
</div>

The numbers marked in diagram:
1. Auth & authz APIs for external calls
2. Render templates via Helm;
3. Persistent jobs and configurations of FATE deployment
4. KubeFATE service is hosted in Kubernetes as one app

There are two parts of KubeFATE:
* The KubeFATE CLI. KubeFATE CLI is a executable binary helps to quickly initial and manage FATE cluster with interactive CLIs. It can be run outside of the Kubernetes, and does not require Kubernetes authz. Eventually, KubeFATE CLI will call KubeFATE Service for detail operations with KubeFATE user token.
* The KubeFATE Service. As KubeFATE provides RESTful APIs for manage FATE clusters. A KubeFATE service will be deployed in Kubernetes, and exposed APIs via [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/). For the auth and authz, KubeFATE service implements [JWT](https://jwt.io/introduction/), and neutral to other security solutions which can be added to Kubernetes ingress.

KubeFATE is designed to seperate the detail FATE cluster configuration including most of the version specified setting. Ideally, KubeFATE CLI and service can work for several FATE releases.

## Use Scenarios
Suppose in a organization, there are two roles:
* System Admin: who responisble for the infrastructure management as well as Kubernetes administration
* ML Infra. Ops.: who responsible for managing the machine learning cluster like FATE

<div align="center">
  <img src="./images/swim.PNG">
</div>

### Initial a new FATE deployment
#### Create role, namespace and other resource in Kubernetes
The sample yaml can be [rbac-config.yaml](./rbac-config.yaml). In the sample yaml, we create a kube-fate namespace for KubeFATE service. More limitation can be applied to kube-fate namespace, refer to [Kubernetes Namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/), [Configure Memory and CPU Quotas for Namespace](https://kubernetes.io/docs/tasks/administer-cluster/manage-resources/quota-memory-cpu-namespace/)
```
kubectl apply -f ./rbac-config.yaml
```
Note that, the default username and password of KubeFATE service can be set in `rbac-config.yaml` Secret->kubefate-secret->stringData as:

```
stringData:
  kubefateUsername: admin
  kubefatePassword: admin
```

#### Prepare domain and deploy KubeFATE in Kubernetes
Because KubeFATE service expose RESTful APIs for external integrated system, system admin have to prepare a domain for KubeFATE service. In our sample config, there is `kubefate.net` . And also, system admin should create a namespace (e.g. fate-9999), limit its quota for FATE deployment, and give the infos to ML Infra. Ops.
```
kubectl apply -f ./kubefate.yaml
kubectl create namespace fate-9999
```
The details of KubeFATE service configuration, please refer to: [KubeFATE service Configuration Guild](../docs/configurations/kubefate_service_configuration.md).

#### Prepare cluster conf. and deploy FATE
When the system admin deployed KubeFATE service and prepared the namespace for FATE. The ML Infra. Ops. is able to start FATE deployment. According to the infomations from SA, there a `config.yaml` for `kubefate` CLI is required. It contains KubeFATE access username and password, the KubeFATE service URL.

```
log:
  level: info
user:
  username: admin
  password: admin

serviceurl: kubefate.net
```

|Name       |Type    |Description                                                       |
|-----------|--------|------------------------------------------------------------------|
|log        |scalars |The log level of command line.                                    |
|user       |mappings|User name and password when logging into KubeFATE service.        |
|serviceurl |scalars |kubeFATE service's ingress domain name, defined in kubefate.yaml. |

And, according to the FATE deploy plan, create a `cluster.yaml` for cluster configuration. The details of Cluster configuration, please refer to: [FATE Cluster Configuration Guild](../docs/configurations/FATE_cluster_configuration.md). Then intall FATE cluster,

```
$ kubefate cluster install -f ./cluster.yaml
create job success, job id=fe846176-0787-4879-9d27-622692ce181c
```
#### Check the status of "Install Cluster" job
A job will be created for installing FATE cluster. Use `kubefate job describe` to check the status of job, util we see the result turns to `install success`

```
$ kubefate job describe fe846176-0787-4879-9d27-622692ce181c
StartTime       2020-05-15 06:34:25
EndTime         2020-05-15 06:35:14
Status          Success
Creator         admin
ClusterId       27e37a60-fffb-4031-a76f-990b2ff43cf2
Result          install success
SubJobs         []
```
#### Decribe the cluster and find FATE access infos
When we see the `install cluster` job success, use `kubefate cluster describe` to check the FATE access infos
```
$ kubefate cluster describe 27e37a60-fffb-4031-a76f-990b2ff43cf2
UUID            27e37a60-fffb-4031-a76f-990b2ff43cf2
Name            fate-9999
NameSpace       fate-9999
ChartName       fate
ChartVersion    v1.4.0
REVISION        1
Status          Running
Values          {"chartName":"fate","chartVersion":"v1.4.0","name":"fate-9999","namespace":"fate-9999","nodemanager":{"count":3,"list":[{"accessMode":"ReadWriteOnce","existingClaim":"","name":"nodemanager","nodeSelector":{},"sessionProcessorsPerNode":2,"size":"1Gi","storageClass":"nodemanager","subPath":"nodemanager"}],"sessionProcessorsPerNode":4},"partyId":9999,"persistence":false, ...... }
ChartName       fate
Info            map[dashboard:[9999.fateboard.kubefate.net] ip:192.168.9.2 ...... ]
```

### Other Common Use Scenarios
#### [Manage FATE and FATE-Serving Version](../docs/Manage_FATE_and_FATE-Serving_Version.md)
#### [Update and Delete a FATE Cluster](../docs/Update_and_Delete_a_FATE_Cluster.md)

## KubeFATE Service RESTful APIs Reference
[API Reference](https://app.swaggerhub.com/apis-docs/vmware-octo/kubefate2/1.0.0-oas3#/cluster/createcluster)
