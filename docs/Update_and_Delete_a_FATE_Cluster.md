# Update and Delete a FATE Cluster
Besides install a new cluster, KubeFATE provides command to update, delete and describe a cluster. And the cluster being managed, not only FATE, but also FATE-Serving, even other clusters to add supports with chart(Refer to [Manage FATE and FATE-Serving Version](./Manage_FATE_and_FATE-Serving_Version.md)).

## Show the detail of a cluster
With the command `kubefate cluster describe ${cluster_id}`, the details of cluster deployed can be shown, including information of how to access the cluster.

```
$ kubefate cluster ls
UUID                                    NAME            NAMESPACE       REVISION        STATUS  CHART   ChartVERSION
246ba46d-a748-42c8-a499-de44b7d1fa4e    fate-10000      fate-10000      1               Running fate    v1.5.0
47a82461-47a4-4608-b208-4d843288bc6b    fate-9999       fate-9999       1               Running fate    v1.5.0

$ kubefate cluster describe 47a82461-47a4-4608-b208-4d843288bc6b
UUID            246ba46d-a748-42c8-a499-de44b7d1fa4e             Name            fate-9999
NameSpace       fate-9999
ChartName       fate
ChartVersion    v1.5.0
REVISION        1
Age             92s
Status          Running
Spec            name: fate-9999
                namespace: fate-9999
                chartName: fate
                chartVersion: v1.5.0
                partyId: 9999
                ......
                
Info            dashboard:
                - party9999.notebook.example.com
                - party9999.fateboard.example.com
                ip: 192.168.0.1
                pod:
                - clustermanager-78f98b85bf-ph2hv
                ......
                status:
                  modules:
                    client: Running
                    clustermanager: Running
                    fateboard: Running
                    mysql: Running
                    nodemanager-0: Running
                    python: Running
                    rollsite: Running
```

## Update a cluster's config (and re-deploy)
KubeFATE provides command to redeploy the cluster, just change the cluster.yaml (The detail refers to [FATE Cluster Configuration](./configurations/FATE_cluster_configuration.md)). It is a useful feature for (but not only): 
1. Scaling out one or several components in cluster, such as eggs;
2. Change the configurations of one or several components;
3. Update FATE/FATE-Serving version. (**Note: If use NFS setting, it may cause FATE internal error, because the data schema changed. Online upgrade feature is developing and depends on FATE support.**)

Comparing to the yaml shown in our tutorials, [Build Two Parties FATE Cluster in One Linux Machine with MiniKube](./tutorials/Build_Two_Parties_FATE_Cluster_in_One_Linux_Machine_with_MiniKube.md), the setting of egg has been change from 1 to 2 for scaling out nodemanager.
```
name: fate-9999
namespace: fate-9999
chartName: fate
chartVersion: v1.5.0
partyId: 9999
registry: ""
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
    partyIp: 192.168.10.1
    partyPort: 30101

nodemanager:
  count: 2
```

then, we can just run,
```
$ kubefate cluster update -f ./fate-9999.yaml
create job success, job id=e8f5440c-a245-44f8-856d-3721efd9c9cf

$ kubectl get pod -n fate-9999 | grep nodemanager
nodemanager-0-84c98f98cb-smt62          1/1     Running       0          31s
nodemanager-1-7c4956c466-fxhpz          1/1     Running       0          31s

$ kubefate cluster ls
UUID                                    NAME            NAMESPACE       REVISION        STATUS  CHART   ChartVERSION
246ba46d-a748-42c8-a499-de44b7d1fa4e    fate-10000      fate-10000      1               Running fate    v1.5.0
47a82461-47a4-4608-b208-4d843288bc6b    fate-9999       fate-9999       2               Running fate    v1.5.0
```

and find the nodemanager has been scaled out to 2 instances. And the revision have been updated to 2. 

## Delete a cluster
To delete a cluster can be use command `kubefate cluster delete ${cluster_id}`
```
$ kubefate cluster delete 47a82461-47a4-4608-b208-4d843288bc6b
create job success, job id=3883a4d1-665c-480f-9cfa-383a05592cc3

$ kubefate cluster ls
UUID                                    NAME            NAMESPACE       REVISION        STATUS  CHART   ChartVERSION
246ba46d-a748-42c8-a499-de44b7d1fa4e    fate-10000      fate-10000      1               Running fate    v1.5.0
```