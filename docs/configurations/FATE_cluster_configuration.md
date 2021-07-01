# FATE cluster configuration
`cluster.yaml` declares information about the FATE cluster to be deployed, which KubeFATE CLI uses to deploy the FATE cluster.

## cluster.yaml
| Name                      | Type               | Description                                                  |
| ------------------------- | ------------------ | ------------------------------------------------------------ |
| * name                    | scalars            | FATE cluster name.                                           |
| * namespace               | scalars            | Kubernetes namespace for FATE cluster.                       |
| * chartName               | scalars            | FATE chart name. (fate/fate-serving)                         |
| * chartVersion            | scalars            | FATE chart corresponding version.                            |
| * partyId                 | scalars            | FATE cluster party id.                                       |
| registry                  | scalars            | Other fate images sources.                                   |
| pullPolicy                | scalars            | kubernetes images pull policy                                |
| * persistence             | bool               | mysql and nodemanager data persistence.                      |
| podSecurityPolicy.enabled | bool               | if `true`, create & use Pod Security Policy resources        |
| * modules                 | sequences          | Modules to be deployed in the FATE cluster.                  |
| backend                   | set(eggroll,spark) | Configure cluster computing engine( eggroll or spark)        |
| host                      | mappings           | Custom domain of FATE UI component                           |
| rollsite                  | mappings           | Configuration of FATE cluster `rollsite` module.             |
| nodemanager               | mappings           | Configuration of FATE cluster `nodemanager` module.          |
| python                    | mappings           | Configuration of FATE cluster `python` module.               |
| mysql                     | mappings           | Configuration of FATE cluster `mysql` module.<br />If you use your own redis, please delete this item. |
| externalMysqlIp           | scalars            | Access your own MySQL.                                       |
| externalMysqlPort         | scalars            | Access your own MySQL.                                       |
| externalMysqlDatabase     | scalars            | Access your own MySQL.                                       |
| externalMysqlUser         | scalars            | Access your own MySQL.                                       |
| externalMysqlPassword     | scalars            | Access your own MySQL.                                       |
| servingIp                 | scalars            | Serving cluster connected to fate.                           |
| servingPort               | scalars            | Serving cluster connected to fate.                           |
| spark                     | mappings           | Configuration of FATE cluster `spark` module.                |
| hdfs                      | mappings           | Configuration of FATE cluster `hdfs` module.                 |
| nginx                     | mappings           | Configuration of FATE cluster `nginx` module.                |
| rabbitmq                  | mappings           | Configuration of FATE cluster `rabbitmq` module.             |



### list of modules

- rollsite

- clustermanager
- nodemanager
- mysql
- python
- fateboard
- client
- spark
- hdfs
- nginx
- rabbitmq



### host mappings

| Name       | Type    | Description                          |
| ---------- | ------- | ------------------------------------ |
| fateboard  | scalars | Configuration of Fateboard UI domain |
| client     | scalars | Configuration of Notebook UI domain  |
| sparkUI    | scalars | Configuration of Spark UI domain     |
| rabbitmqUI | scalars | Configuration of Rabbitmq UI domain  |



### rollsite mappings
It is used to declare the `rollsite ` module in the FATE cluster to be deployed.

| Name         | subitem     | Type      | Description                                                  |
| ------------ | ----------- | --------- | ------------------------------------------------------------ |
| type         |             | scalars   | Kubernetes ServiceTypes, default is NodePort.                |
| nodePort     |             | scalars   | The port used by `proxy` module's kubernetes service, default range: 30000-32767. |
| partyList    |             | sequences | If this FATE cluster is exchange cluster, partyList is all party's sequences of all parties proxy address. If this FATE cluster is one of participants, delete this configuration item. |
| partyList    | partyId     | scalars   | Participant FATE cluster party ID.                           |
| partyList    | partyIp     | scalars   | Participant FATE cluster IP.                                 |
| partyList    | partyPort   | scalars   | Participant FATE cluster port.                               |
| exchange     |             | mappings  | FATE cluster `exchange` module's ip and port.                |
| exchange     | ip          | mappings  | FATE cluster `exchange` module's ip. .                       |
| exchange     | port        | mappings  | FATE cluster `exchange` module's port.                       |
| nodeSelector |             | mappings  | kubernetes nodeSelector.                                     |
| polling      |             |           | rollsite support polling                                     |
| polling      | enabled     |           | enable polling                                               |
| polling      | type        |           | polling type (server/client)                                 |
| polling      | server      |           | if type choose client, you need a polling server.            |
| polling      | clientList  |           | if type choose server, this rollsite serve for clientList.   |
| polling      | concurrency |           | if type choose server, polling client concurrency.           |

FATE cluster has two deployment modes: with exchange and without exchange.
#### Exchange mode:
Every party connected to the `exchange`, which has the proxy addresses of all parties.
- For `exchange` cluster, only need to deploy `proxy` modules. In `proxy` configuration, no need `exchange` item, need to has the proxy addresses of all parties in `partylist`.
- For other FATE clusters, need to fill in the exchange ip and port. Can delete `partylist` configuration item.

#### Direct connection mode:
The parties are directly connected.
- No need to fill in the `exchange` ip and port.
- `partyList` needs the addresses of all other FATE clusters proxy.

### nodemanager mappings

| Name                       | SubItem                    | Type      | Description                                                  |
| -------------------------- | -------------------------- | --------- | ------------------------------------------------------------ |
| count                      |                            | scalars   | Number of nodes deployed nodemanager.                        |
| session-Processors-PerNode |                            | scalars   | Configuration of FATE cluster `nodemanager` module.          |
| list                       |                            | sequences | List of nodemanager nodes.                                   |
| list                       | name                       | scalars   | nodemanager node name.                                       |
| list                       | nodeSelector               | mappings  | kubernetes nodeSelector.                                     |
| list                       | session-Processors-PerNode | scalars   | Configuration of FATE cluster `nodemanager` module.          |
| list                       | subPath                    | scalars   | Path of data persistence, specify the "subPath" if the PVC is shared with other components. |
| list                       | existingClaim              | scalars   | Use the existing PVC which must be created manually before bound. |
| list                       | storageClass               | scalars   | Specify the "storageClass" used to provision the volume. Or the default. StorageClass will be used(the default). Set it to "-" to disable dynamic provisioning. |
| list                       | accessMode                 | scalars   | Kubernetes Persistent Volume Access Modes: <br />ReadWriteOnce<br />ReadOnlyMany <br />ReadWriteMany. |
| list                       | size                       | scalars   | Match the volume size of PVC.                                |

### python mappings

| Name                  | Type     | Description                                                  |
| --------------------- | -------- | ------------------------------------------------------------ |
| type                  | scalars  | Kubernetes ServiceTypes, default is NodePort.<br />Other modules can connect to the fateflow |
| nodePort              | scalars  | The port used by `proxy` module's kubernetes service, default range: 30000-32767. |
| nodeSelector          | mappings | kubernetes nodeSelector.                                     |
| enabledNN             | bool     | If or not neural network workflow is required                |
| spark                 | mappings | If you use your own spark, modify the configuration          |
| spark.cores_per_node  | scalars  | configuration of FATE fateflow module                        |
| spark.nodes           | scalars  | configuration of FATE fateflow module                        |
| spark.existingSpark   | scalars  | If you need to use the existing spark , you can set this configuration |
| spark.driverHost      | scalars  | call back IP of spark executor                               |
| spark.driverHostType  | scalars  | service type of spark driver                                 |
| spark.portMaxRetries  | scalars  | spark driver's configuration                                 |
| driverStartPort       | scalars  | spark driver start port                                      |
| blockManagerStartPort | scalars  | spark driver blockManager start port                         |
| pysparkPython         | scalars  | spark worker node python PATH                                |
| hdfs                  | mappings | If you do not need to use the spark configuration, you can use the spark configuration |
| rabbitmq              | mappings | If you do not need to use the spark configuration, you can use the spark configuration |
| nginx                 | mappings | If you do not need to use the spark configuration, you can use the spark configuration |



### Mysql mappings

Configuration of kubernetes deployment mysql.

| Name          | Type     | Description                                                  |
| ------------- | -------- | ------------------------------------------------------------ |
| nodeSelector  | mappings | kubernetes nodeSelector.                                     |
| ip            | scalars  | Allow other modules to connect to MySQL.                     |
| port          | scalars  | Mysql port.                                                  |
| database      | scalars  | Database name of MySQL.                                      |
| user          | scalars  | User of MySQL.                                               |
| password      | scalars  | User password of MySQL.                                      |
| subPath       | scalars  | Path of data persistence, specify the "subPath" if the PVC is shared with other components. |
| existingClaim | scalars  | Use the existing PVC which must be created manually before bound. |
| storageClass  | scalars  | Specify the "storageClass" used to provision the volume. Or the default. StorageClass will be used(the default). Set it to "-" to disable dynamic provisioning. |
| accessMode    | scalars  | Kubernetes Persistent Volume Access Modes: <br />ReadWriteOnce<br />ReadOnlyMany <br />ReadWriteMany. |
| size          | scalars  | Match the volume size of PVC.                                |



### spark mappings

Configuration of kubernetes deployment spark.

| Name              | SubItem      | Type     | Description                  |
| ----------------- | ------------ | -------- | ---------------------------- |
| master/<br>worker | Image        | scalars  | Image of spark components    |
|                   | ImageTag     | scalars  | ImageTag of spark components |
|                   | replicas     | scalars  | Number of copies of pod      |
|                   | resources    | mappings | resources of Kubernetes      |
|                   | nodeSelector | mappings | kubernetes nodeSelector.     |
|                   | type         | scalars  | Kubernetes ServiceTypes.     |

### hdfs mappings

Configuration of kubernetes deployment hdfs.

| Name                  | SubItem      | Type     | Description                                      |
| --------------------- | ------------ | -------- | ------------------------------------------------ |
| namenode/<br>datanode | nodeSelector | mappings | kubernetes nodeSelector.                         |
|                       | type         | scalars  | Kubernetes ServiceTypes, default is `ClusterIp`. |



### nginx mappings

Configuration of kubernetes deployment hdfs.

| Name         | Type     | Description                  |
| ------------ | -------- | ---------------------------- |
| nodeSelector | mappings | kubernetes nodeSelector.     |
| type         | scalars  | Kubernetes ServiceTypes.     |
| nodePort     | scalars  | Kubernetes Service NodePort. |
| route_table  | mappings | route table of FATE          |

*example of route_table*:

```bash
10000: 
  proxy: 
  - host: 192.168.0.1 
    port: 30103
  fateflow: 
  - host: 192.168.0.1
    port: 30102
9999: 
  proxy: 
  - host: 192.168.0.2 
    port: 30093
  fateflow: 
  - host: 192.168.0.2
    port: 30092
8888: 
  proxy: 
  - host: 192.168.0.3 
    port: 30083
  fateflow: 
  - host: 192.168.0.3
    port: 30082 
```



### rabbitmq mappings

Configuration of kubernetes deployment rabbitmq .

| Name         | Type     | Description                                      |
| ------------ | -------- | ------------------------------------------------ |
| nodeSelector | mappings | kubernetes nodeSelector.                         |
| type         | scalars  | Kubernetes ServiceTypes, default is `ClusterIp`. |
| nodePort     | scalars  | Kubernetes Service NodePort.                     |
| default_user | scalars  | configuration of rabbitmq.                       |
| default_pass | scalars  | configuration of rabbitmq.                       |
| user         | scalars  | configuration of rabbitmq.                       |
| password     | scalars  | configuration of rabbitmq.                       |
| route_table  | mappings | route table of rabbitmq.                         |

*example of route_table*:

```bash
10000:
  host: 192.168.0.1
  port: 30104
9999:
  host: 192.168.0.2
  port: 30094
8888:
  host: 192.168.0.3
  port: 30084
```

