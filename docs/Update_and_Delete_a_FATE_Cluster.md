# Update and Delete a FATE Cluster
Besides intall a new cluster, KubeFATE provides command to update, delete and describe a cluster. And the cluster being managed, not only FATE, but also FATE-Serving, even other clusters to add supports with chart(Refer to [Manage FATE and FATE-Serving Version](./Manage_FATE_and_FATE-Serving_Version.md)).

## Show the detail of a cluster
With the command `kubefate cluster describe ${cluster_id}`, the details of cluster deployed can be shown, including information of how to access the cluster.

```
$ kubefate cluster ls
UUID                                    NAME            NAMESPACE       REVISION        STATUS  CHART   ChartVERSION
246ba46d-a748-42c8-a499-de44b7d1fa4e    fate-10000      fate-10000      1               Running fate    v1.3.0-a
47a82461-47a4-4608-b208-4d843288bc6b    fate-9999       fate-9999       1               Running fate    v1.3.0-a

$ kubefate cluster describe 246ba46d-a748-42c8-a499-de44b7d1fa4e
UUID            246ba46d-a748-42c8-a499-de44b7d1fa4e                                                                                                                                                                                                        
Name            fate-10000                                                                                                                                                                                                                                  
NameSpace       fate-10000                                                                                                                                                                                                                                  
ChartName       fate                                                                                                                                                                                                                                        
ChartVersion    v1.3.0-a                                                                                                                                                                                                                                    
Revision        1                                                                                                                                                                                                                                           
Status          Running                                                                                                                                                                                                                                     
Values          {"chartVersion":"v1.3.0-a","egg":{"count":1},"modules":["proxy","egg","federation","metaService","mysql","redis","roll","python"],"name":"fate-10000","namespace":"fate-10000","partyId":10000,"proxy":{"nodePort":30010,"partyList":[{"partyId":9999,"partyIp":"192.168.100.123","partyPort":30009}],"type":"NodePort"}}
Config          map[chartVersion:v1.3.0-a egg:map[count:1] modules:[proxy egg federation metaService mysql redis roll python] name:fate-10000 namespace:fate-10000 partyId:10000 proxy:map[nodePort:30010 partyList:[map[partyId:9999 partyIp:192.168.100.123 partyPort:30009]] type:NodePort]]
Info            map[dashboard:10000.fateboard.kubefate.net ip:192.168.100.123 modules:[egg0-5b44548fbd-rvmd8 federation-6d799b5cfd-d92wq meta-service-54db9f8fbc-2n9w2 mysql-6bc77fc46c-5dq2w proxy-8d758c997-sgpdr python-77bb96fd78-rfq5s redis-9546f56b-fw5cf roll-77dfbb54dc-g897x] port:30010]

```

## Update a cluster's config (and re-deploy)
KubeFATE provides command to redeploy the cluster, just change the cluster.yaml (The detail refers to [FATE Cluster Configuration](./configurations/FATE_cluster_configuration.md)). It is a useful feature for (but not only): 
1. Scaling out one or several components in cluster, such as eggs;
2. Change the configurations of one or several components;
3. Update FATE/FATE-Serving version. (**Note: If use NFS setting, it may cause FATE internal error, because the data schema changed. Online upgrade feature is developing and depends on FATE support.**)

Commparing to the yaml shown in our tutorials, [Build Two Parties FATE Cluster in One Linux Machine with MiniKube](./tutorials/Build_Two_Parties_FATE_Cluster_in_One_Linux_Machine_with_MiniKube.md), the setting of egg has been change from 1 to 2 for scaling out egg.
```
name: fate-9999
namespace: fate-9999
chartVersion: v1.3.0-a
partyId: 9999
modules:
  - proxy
  - egg
  - federation
  - metaService
  - mysql
  - redis
  - roll
  - python

proxy:
  type: NodePort
  nodePort: 30009
  partyList:
    - partyId: 10000
      partyIp: 192.168.100.123
      partyPort: 30010
egg:
  count: 2
```

then, we can just run,
```
$ kubefate cluster update -f ./fate-9999.yaml
create job success, job id=e8f5440c-a245-44f8-856d-3721efd9c9cf

$ kubectl get pod -n fate-9999 | grep egg
egg0-79768fbffb-2v5qv          1/1     Running       0          31s
egg1-6bd6b965cf-g64b6          1/1     Running       0          31s

$ kubefate cluster ls
UUID                                    NAME            NAMESPACE       REVISION        STATUS  CHART   ChartVERSION
246ba46d-a748-42c8-a499-de44b7d1fa4e    fate-10000      fate-10000      1               Running fate    v1.3.0-a
47a82461-47a4-4608-b208-4d843288bc6b    fate-9999       fate-9999       2               Running fate    v1.3.0-a
```

and find the egg has been scaled out to 2 instances. And the revision have been updated to 2. 

## Delete a cluster
To delete a cluster can be use command `kubefate cluster delete ${cluster_id}`
```
$ kubefate cluster delete 47a82461-47a4-4608-b208-4d843288bc6b
create job success, job id=3883a4d1-665c-480f-9cfa-383a05592cc3

$ kubefate cluster ls
UUID                                    NAME            NAMESPACE       REVISION        STATUS  CHART   ChartVERSION
246ba46d-a748-42c8-a499-de44b7d1fa4e    fate-10000      fate-10000      1               Running fate    v1.3.0-a
```