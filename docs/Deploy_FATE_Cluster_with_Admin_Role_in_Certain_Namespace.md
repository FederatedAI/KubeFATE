# Deploy FATE Cluster with Admin Role in Certain Namespace

## Background

When deploying KubeFATE and FATE cluster with Kubernetes, user may not have full control over every resource in the cluster. Kubernetes provides [Role-based access control(RBAC) authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) to restrict the actions a user can take. In default case, we ask user to create a new role with a wide set of permissions and create clusterrolebinding to grant access within the whole cluster. But it may not work when user only has access to some limited resources in particular namespace. We studied all the [user-facing roles](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles)(`cluster-admin`, `admin`, `edit`, `view`) in Kubernetes and found that to deploy KubeFATE and run FATE jobs, `admin` is necessary because of the privilege to create roles and role bindings within the namespace. And in this case, some configuration files need to be modified.

## Deploy Steps

### Environment

Start with a K8s cluster which does not have KubeFATE preinstalled, but have 3 namespaces (`fate-exchange`，`fate-9999`，`fate-10000`) and 2 users:

- User `9999` which bound with `admin` role in namespace `fate-9999` and `fate-exchange`
- User `10000` which bound with `admin` role in namespace `fate-10000`

Similar to [Deploy an exchange central multi parties federated learning network with KubeFATE](https://github.com/FederatedAI/KubeFATE/wiki/Deploy-an-exchange-central-multi-parties-federated-learning-network-with-KubeFATE), in this document we will deploy two Parties with an exchange and run FATE test job.

| party       | party ID | owner      | namespace     | K8s version | KubeFATE version | FATE version |
|-------------|----------|------------|---------------|-------------|------------------|--------------|
| exchange    | 1        | user-9999  | fate-exchange | v1.24.3     | v1.4.4           | v1.8.0       |
| party-9999  | 9999     | user-9999  | fate-9999     | v1.24.3     | v1.4.4           | v1.8.0       |
| party-10000 | 10000    | user-10000 | fate-10000    | v1.24.3     | v1.4.4           | v1.8.0       |

### Install KubeFATE

Because there is no KubeFATE in the cluster, each user needs to [install KubeFATE](https://github.com/FederatedAI/KubeFATE/tree/master/k8s-deploy#readme) separately. Default configuration files needs to be modified beforehand.

User `9999` will use the `rbac-config.yaml` file below:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kubefate-admin
  namespace: fate-9999
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: sc-9999-edit-binding
  namespace: fate-exchange
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: edit
subjects:
  - kind: ServiceAccount
    name: kubefate-admin
    namespace: fate-9999
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: sc-edit-binding
  namespace: fate-9999
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: edit
subjects:
  - kind: ServiceAccount
    name: kubefate-admin
    namespace: fate-9999
---
apiVersion: v1
kind: Secret
metadata:
  name: kubefate-secret
  namespace: fate-9999
type: Opaque
stringData:
  kubefateUsername: admin
  kubefatePassword: admin
  mariadbUsername: kubefate
  mariadbPassword: kubefate
```

Compare to default `rbac-config.yaml`, all the namespaces should be modified because KubeFATE should be deployed to user's own namespace. In general, we need to remove any clusterrole or clusterrolebinding and create only role and rolebinding in the corresponding namespaces. PodSecurityPolicy was also removed because `admin` role cannot access the resource at cluster scope.

User `10000` use the `rbac-config.yaml` similar to user `9999` but doesn't create rolebinding in `fate-exchange` namespace.

Default namespace in `kubefate.yaml` should be changed to the namespace with privileges correspondingly.

Execute the commands below after all files are ready:

```bash
$ kubectl apply -f ./rbac-config.yaml
$ kubectl apply -f ./kubefate.yaml
```

### Deploy Exchange

In this document we use Nginx and ATS to deploy Exchange, so the certificate configuration between each Party and Exchange should be solved first. Refer to [pulsar and certificate generation of ATS](https://github.com/FederatedAI/FATE/blob/v1.6.0/cluster-deploy/doc/fate_on_spark/FATE_on_spark_with_pulsar_zh.md) for more information.

The template configuration file of Exchange is [trafficServer.yaml](https://github.com/FederatedAI/KubeFATE/blob/master/k8s-deploy/examples/party-exchange/trafficServer.yaml). And `route_table` of nginx and trafficServer will be configured according to the services address and port of the Party after the FATE clusters are up and running.

When `cluster-exchange.yaml` is configured, user 9999 will deploy it to `fate-exchange` namespace with KubeFATE:

```sh
kubefate cluster install -f ./cluster-exchange.yaml
```

Check the status of Exchange cluster to make sure it is up and running:

```sh
kubefate cluster ls
```

### Add Parties

We use Spark + Pulsar as the backend of FATE. So the certificates for Pulsar should be installed beforehand like exchange.

Refer to [cluster-spark-pulsar.yaml](https://github.com/FederatedAI/KubeFATE/blob/master/k8s-deploy/examples/party-9999/cluster-spark-pulsar.yaml) to get the configuration template of FATE cluster. Don't forget to change the ip and ports of Exchange service.

`cluster-exchange.yaml` should be modified after Party 9999 and Party 1000 are all set. Add hosts and ports to route tables to connect each Party. And update the Exchange cluster:

```sh
kubefate cluster update -f ./cluster-exchange.yaml
```

### Test

When all deployments are success, there is a federated learning network in the K8s cluster which contains two Parties and an Exchange. Create a new terminal in the notebook of Party 9999 and run command below to check if it works smoothly:

```sh
flow test toy -gid 9999 -hid 10000
```

## Tips

When installing multiple KubeFATE services in one K8s cluster, the domains should be unique, or they may conflict with each other. Change the `serviceurl` in `config.yaml` under `./kubefate` directory to access different KubeFATE service.

## Limits

There will be some limits when user only has admin privilege in certain namespace:

1. Some commands of KubeFATE won't work:
   - `kubefate namespace ls`
   - `kubefate cluster describe`
2. PodSecurityPolicy can not be enabled.
