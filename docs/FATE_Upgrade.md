# Using KubeFATE to upgrade a FATE cluster

## Overview
Since KubeFATE v1.4.5, KubeFATE can help to upgrade a FATE cluster.
This document is going to illustrate how to do that, and also show the limitations for this new feature due to historical reasons.

## Which versions are supported?

"Supported" means that support using KubeFATE CLI to upgrade a FATE cluster automatically.

The supporting matrices depend on the configuration of the data persistence of the FATE cluster, in specific:

### Case 1: When using existing persistence volume claim for MySQL
Support FATE versions

| From/To         | v1.7.2 | v1.8.0 | v1.9.0 | future versions |
|-----------------|--------|--------|--------|-----------------|
| v1.7.1 or lower | no     | no     | no     | no              |
| v1.7.2          |        | yes    | yes    | yes             |
| v1.8.0          |        |        | yes    | yes             |
| v1.9.0          |        |        |        | yes             |

#### Q&A
1. What does "yes" or "no" means in the form?
   1. "yes" means that the FATE cluster can work properly after upgrading, also, the data produced during the previous version, such as the job info or model info, can still be accessed in the new version's FATE cluster.
2. Why "v1.7.1 or lower" is not supported?
   1. This is because in v1.7.1 and previous versions, the MySQL image version in our helm chart is set to 8, which is like setting to the latest version of the MySQL 8 series. Since v1.7.2, we set the MySQL image version to 8.0.28. So for example if you install a v1.7.1 cluster now, your MySQL version should be larger than 8.0.28, however, downgrade a MySQL database is not supported. [reference](https://dev.mysql.com/doc/refman/8.0/en/downgrading.html).
3. Is there a workaround for "v1.7.1 or lower"?
   1. Yes, please check [workaround](#workarounds).


### Case 2: When using storage class to auto-provision the persistence
Support FATE versions

| From/To         | v1.7.2 | v1.8.0 | v1.9.0 | future versions |
|-----------------|--------|--------|--------|-----------------|
| v1.7.1 or lower | no     | no     | no     | no              |
| v1.7.2          |        | yes    | no     | no              |
| v1.8.0          |        |        | no     | no              |
| v1.9.0          |        |        |        | yes             |

#### Q&A
1. Why "v1.7.1 or lower" is not supported?
   1. Same with above.
2. Why upgrade from v1.7.2 to v1.8.0 is supported, but upgrade from v1.8.0 to v1.9.0 is not supported?
   1. In v1.9.0, we make some adjustments to the app type of several FATE K8s components. Basically, for each "deployment" who will have a PVC/PV when persistence is enabled, we change the app type from "deployment" to "statefulSet". We believe this is the right thing to do, but the side effect is that the data from v1.8.0 cannot be carried forward. The good thing is, in the future versions we will not change the app type once again, so this gap will happen only once.

### Case 3: When disabling persistence

This is similar with destroy an old-versioned FATE cluster and create a new-versioned one. 
The new one would be a brand new one with no data from the previous version.
It can work properly, but it will abandon all the history data.


## Usage guidance

In this section, we will introduce how to use KubeFATE to conduct an upgrade.
We will upgrade from v1.7.2 to v1.8.0 to illustrate the processes.

### Install a FATE cluster of v1.7.2

The cluster.yaml:

```yaml
name: fate-9999
namespace: fate-9999
chartName: fate
chartVersion: v1.7.2
partyId: 9999
registry: ""
imageTag: 1.7.2-release
pullPolicy:
imagePullSecrets:
   - name: myregistrykey
persistence: true
istio:
   enabled: false
podSecurityPolicy:
   enabled: false
ingressClassName: nginx
modules:
   - rollsite
   - clustermanager
   - nodemanager
   - mysql
   - python
   - fateboard
   - client

backend: eggroll

ingress:
   fateboard:
      hosts:
         - name: party9999.fateboard.example.com
   client:
      hosts:
         - name: party9999.notebook.example.com

rollsite:
   type: NodePort
   nodePort: 30091
   partyList:
      - partyId: 10000
        partyIp: 192.168.10.1
        partyPort: 30101

python:
   type: NodePort
   httpNodePort: 30097
   grpcNodePort: 30092
   logLevel: INFO
   storageClass: nfs-client

client:
   storageClass: nfs-client

servingIp: 192.168.9.1
servingPort: 30095

mysql:
   storageClass: nfs-client

nodemanager:
   count: 2
   sessionProcessorsPerNode: 4
   storageClass: nfs-client
   accessMode: ReadWriteOnce
   size: 2Gi
   list:
      - name: nodemanager
        nodeSelector:
        sessionProcessorsPerNode: 4
        subPath: "nodemanager"
        existingClaim: ""
        storageClass: nfs-client
        accessMode: ReadWriteOnce
        size: 1Gi
```

In this example, we have already installed [nfs-subdir-external-provisioner](https://github.com/kubernetes-sigs/nfs-subdir-external-provisioner) as the storage class.
For the background of how to set up a nfs server, you could check [here](https://www.itzgeek.com/how-tos/linux/centos-how-tos/how-to-setup-nfs-server-on-centos-7-rhel-7-fedora-22.html)

After installation, the pods are like:
```
client-75f67f5d47-svfh2           1/1     Running   0          2m30s
clustermanager-7bd6fb46c8-bnx57   1/1     Running   0          2m30s
mysql-77f95d4844-8dfjj            1/1     Running   0          2m30s
nodemanager-0-6dbcc56dd4-ddldx    2/2     Running   0          2m30s
nodemanager-1-646ddbf48c-2lllz    2/2     Running   0          2m30s
python-5bd5d69779-c4t5r           2/2     Running   0          2m30s
rollsite-5d45c85c4d-dbpnt         1/1     Running   0          2m30s
```
The PVC/PC are provisioned by the storage class, they are like:

PVC:
```
NAMESPACE   NAME                 STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
fate-9999   client-data          Bound    pvc-e7ea7c79-c5e3-4824-8269-4eac390dcae8   1Gi        RWO            nfs-client     3m27s
fate-9999   mysql-data           Bound    pvc-0865a743-9bfb-4990-bd26-84b0241f0084   1Gi        RWO            nfs-client     3m27s
fate-9999   nodemanager-0-data   Bound    pvc-eaac241e-4d0f-4d06-96d0-30e85e284adc   2Gi        RWO            nfs-client     3m27s
fate-9999   nodemanager-1-data   Bound    pvc-eed13a83-442e-4c9d-bb00-e6a7ce7c05f2   2Gi        RWO            nfs-client     3m27s
fate-9999   python-data          Bound    pvc-5fc77a70-fe67-4b10-8db5-d3b0ef76d8f7   1Gi        RWO            nfs-client     3m27s
```

PV:
```
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                          STORAGECLASS   REASON   AGE
pvc-0865a743-9bfb-4990-bd26-84b0241f0084   1Gi        RWO            Retain           Bound    fate-9999/mysql-data           nfs-client              3m44s
pvc-5fc77a70-fe67-4b10-8db5-d3b0ef76d8f7   1Gi        RWO            Retain           Bound    fate-9999/python-data          nfs-client              3m43s
pvc-e7ea7c79-c5e3-4824-8269-4eac390dcae8   1Gi        RWO            Retain           Bound    fate-9999/client-data          nfs-client              3m44s
pvc-eaac241e-4d0f-4d06-96d0-30e85e284adc   2Gi        RWO            Retain           Bound    fate-9999/nodemanager-0-data   nfs-client              3m44s
pvc-eed13a83-442e-4c9d-bb00-e6a7ce7c05f2   2Gi        RWO            Retain           Bound    fate-9999/nodemanager-1-data   nfs-client              3m44s
```

### Backup the Mysql database
Before any upgrade, we should back up the database. Currently, KubeFATE doesn't support back up the database automatically. This need to be done by the mysqldump tool manually. In the future this could be another feature of KubeFATE.

There are several ways to do the backup work, in this example, we do that in the mysql pod:

```
$ kubectl exec -it mysql-77f95d4844-8dfjj -n fate-9999 bash                       
kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl exec [POD] -- [COMMAND] instead.
root@mysql-77f95d4844-8dfjj:/# 
```

Note that "mysql-77f95d4844-8dfjj" is the pod id which must be different in your FATE cluster. Check the pod id by ```kubectl get pods -n <your namespace>```.

We need to change the directory to "/var/lib/mysql/", because that is where the PV is mounted to.
If we generate a mysql snapshot there, the file will be persisted in the nfs server even the pod is gone.
```
root@mysql-77f95d4844-8dfjj:/# cd /var/lib/mysql/
root@mysql-77f95d4844-8dfjj:/var/lib/mysql# 
```
Then we can run the mysqldump command to create a snapshot of the current status of the database:
```
root@mysql-77f95d4844-8dfjj:/var/lib/mysql# mysqldump -h localhost -u root eggroll_meta > snapshot.sql
root@mysql-77f95d4844-8dfjj:/var/lib/mysql# ls | grep snapshot
snapshot.sql
```
Now we can quit the container of mysql. login to the nfs server and manage the snapshot file.
```
nfsServer ~ → cd /mnt/myshareddir/
nfsServer myshareddir → ls
fate-9999-client-data-pvc-e7ea7c79-c5e3-4824-8269-4eac390dcae8         fate-9999-nodemanager-1-data-pvc-eed13a83-442e-4c9d-bb00-e6a7ce7c05f2
fate-9999-mysql-data-pvc-0865a743-9bfb-4990-bd26-84b0241f0084          fate-9999-python-data-pvc-5fc77a70-fe67-4b10-8db5-d3b0ef76d8f7
fate-9999-nodemanager-0-data-pvc-eaac241e-4d0f-4d06-96d0-30e85e284adc
nfsServer myshareddir → ls fate-9999-mysql-data-pvc-0865a743-9bfb-4990-bd26-84b0241f0084 | grep snapshot
snapshot.sql
```

### Using KubeFATE CLI to do the upgrade
We just need to change two lines of the cluster.yaml file:
```
chartVersion: v1.7.2	  ->   chartVersion: v1.8.0
imageTag: 1.7.2-release   ->   imageTag: 1.8.0-release
```
Then execute:
```
kubefate cluster update -f cluster.yaml
```
A KubeFATE job will be created. At this time, a K8s one-time job will be launched, it will do 2 things:
1. Shut down the "python" pod, which includes FATE-Flow and FATE-Board.
2. It executes one or more .sql script(s) to help upgrade the schema of the MySQL database.

After that, the K8s job will be deleted by default. For debugging purpose, to keep the finished job on the K8s cluster, we can run:
```
kubefate cluster update -f cluster.yaml --keepUpgradeJob
```
Soon, the new pods will be spawned:
```
fate-9999       client-7dd4fc8cbf-jtz2x                                   1/1     Running     0          5m4s
fate-9999       clustermanager-7987759b84-7bgbf                           1/1     Running     0          5m26s
fate-9999       mysql-fc76db6db-rdvw5                                     1/1     Running     0          5m26s
fate-9999       nodemanager-0-54845ccd98-jjptt                            2/2     Running     0          5m25s
fate-9999       nodemanager-1-cdd975587-cn96q                             2/2     Running     0          5m26s
fate-9999       python-57f8c857d8-ckp2z                                   2/2     Running     0          5m36s
fate-9999       rollsite-6b67d75c78-bngqb                                 1/1     Running     0          5m26s
```
At this time, the cluster has already been upgrade to v1.8.0, which means the upgrade has been done:
```
➜ kubefate cluster list                                                        
UUID                                    NAME            NAMESPACE       REVISION        STATUS  CHART   ChartVERSION    AGE  
4a8372d1-a15e-4653-8420-60627d89f858    fate-9999       fate-9999       2               Running fate    v1.8.0          3h35m
```

## Workarounds

As mentioned, due to a historical reason, upgrading from a FATE cluster whose version is less or equal than v1.7.1 is not supported by KubeFATE. However, there is a manual workaround for that purpose:

Suppose we have a v1.7.1 FATE cluster, how to upgrade to v1.7.2 manually without data loss in the MySQL database?
### Export the data in the MySQL database
Get into the container of MySQL by ```kubectl exec -it <mysql_pod_id> -n <your_name_space> bash```, go to "/var/lib/mysql" and run below command:

```
mysqldump -h localhost -u root --no-create-info \
--ignore-table=eggroll_meta.server_node \
--ignore-table=eggroll_meta.t_component_info \
--ignore-table=eggroll_meta.t_component_provider_info \
--ignore-table=eggroll_meta.t_component_registry \
--ignore-table=eggroll_meta.t_engine_registry \
eggroll_meta > data.sql
```
The reason we ignore the tables is: the information are version-specific. The data in those tables will be re-generated when a new "python" container of the new version is spawned. It means nothing to carry these data from v1.7.1 to v1.7.2.

Also, we need to prepare the snapshot for rollback, in case of upgrade failure:
```
mysqldump -h localhost -u root eggroll_meta > snapshot.sql
```

### Destroy the v1.7.1 cluster and install a fresh v1.7.2 cluster
Use ```kubefate cluster delete <uuid>``` and ```kubefate cluster install -f cluster.yaml``` to make this happen.

### Import the data into the new MySQL database
Get into the container of MySQL by ```kubectl exec -it <mysql_pod_id> -n <your_name_space> bash```. Make sure before that the dumped data file "data.sql" has been put into the corresponding nfs folder.

Change directory to "/var/lib/mysql" and run below command:
```
mysql -u root eggroll_meta < data.sql
```
After that, we have "upgraded" the FATE cluster from v1.7.1 to v1.7.2.

## Other notes
1. We also support upgrade over multiple versions, just simply change the version number in the cluster.yaml file, then KubeFATE can figure out which sql scripts it needs to execute for the MySQL database. Currently, if you try to upgrade from v1.7.2 to v1.9.0, skipping v1.8.0, this will only work when you have configured an existing PVC/PV. However, in the future, if you would like to upgrade from v1.9.0 and skip several in-the-middle versions, it will also work even when you are using a storage class.
2. We support changing the architecture during upgrade, for example, you can upgrade from an Eggroll based FATE cluster of version v1.7.2 to a Spark based FATE cluster of version v1.8.0. The new cluster will work properly, although the temporary data in the Eggroll PV will be lost.
3. Again, at current stage, KubeFATE v1.4.5, you are responsible to do the backup of the MySQL database before any kind of upgrade.