# FATE Exchange Cluster Configuration
`cluster-serving.yaml` declares information about the FATE-Serving cluster to be deployed, which KubeFATE CLI uses to deploy the FATE-Serving cluster.

## cluster-exchange.yaml
| Name                      | Type      | Description                                                  |
| ------------------------- | --------- | ------------------------------------------------------------ |
| name                      | scalars   | FATE-Exchange cluster name.                                  |
| namespace                 | scalars   | Kubernetes namespace for FATE-Exchange cluster.              |
| chartName                 | scalars   | FATE chart name. (fate/fate-serving/fate-exchange)           |
| chartVersion              | scalars   | FATE chart corresponding version.                            |
| partyId                   | scalars   | FATE-Exchange cluster party id.                              |
| registry                  | scalars   | Other fate images sources.                                   |
| imageTag                  | scalars   | Image TAG                                                    |
| pullPolicy                | scalars   | Kubernetes images pull policy.                               |
| persistence               | bool      | Redis and servingServer data persistence.                    |
| podSecurityPolicy.enabled | bool      | if `true`, create & use Pod Security Policy resources        |
| modules                   | sequences | Modules to be deployed in the FATE-Serving cluster.          |
| rollsite                  | mappings  | Configuration of FATE cluster `rollsite` module.             |
| trafficServer             | mappings  | Configuration of FATE cluster `nodemanager` module.          |
| nginx                     | mappings  | Configuration of FATE cluster `python` module.<br />If you use your own redis, please delete this item. |
| externalRedisIp           | mappings  | Access your own Redis.                                       |
| externalRedisPort         | scalars   | Access your own Redis.                                       |
| externalRedisDatabase     | scalars   | Access your own Redis.                                       |
| externalRedisPassword     | scalars   | Access your own Redis.                                       |

### rollsite mappings
It is used to declare the `servingProxy` module in the FATE cluster to be deployed.

| Name         | subitem   | Type      | Description                                                  |
| ------------ | --------- | --------- | ------------------------------------------------------------ |
| nodePort     |           | scalars   | The port used by `proxy` module's kubernetes service, default range: 30000-32767. |
| ingerssHost  |           | scalars   | The entrance of FATE-Service api.                            |
| partyList    |           | sequences | If this FATE cluster is exchange cluster, partyList is all party's sequences of all parties proxy address. If this FATE cluster is one of participants, delete this configuration item. |
| partyList    | partyId   | scalars   | Participant FATE cluster party ID.                           |
| partyList    | partyIp   | scalars   | Participant FATE cluster IP.                                 |
| partyList    | partyPort | scalars   | Participant FATE cluster port.                               |
| exchange     |           | mappings  | FATE cluster `exchange` module's ip and port.                |
| exchange     | ip        | mappings  | FATE cluster `exchange` module's ip. .                       |
| exchange     | port      | mappings  | FATE cluster `exchange` module's port.                       |
| nodeSelector |           | mappings  | kubernetes nodeSelector.                                     |

### servingServer mappings

| Name          | SubItem | Type     | Description                                                  |
| ------------- | ------- | -------- | ------------------------------------------------------------ |
| type          |         | scalars  | Kubernetes ServiceTypes, default is NodePort.                |
| nodePort      |         | scalars  | The port used by `proxy` module's kubernetes service, default range: 30000-32767. |
| fateflow      |         | mappings | FATE cluster `python` module's fateflowIp and fateflowPort.  |
| fateflow      | ip      | scalars  | FATE cluster `python` module's fateflowIp.                   |
| fateflow      | port    | scalars  | FATE cluster `python` module's fateflowPort.                 |
| subPath       |         | scalars  | Path of data persistence, specify the "subPath" if the PVC is shared with other components. |
| existingClaim |         | scalars  | Use the existing PVC which must be created manually before bound. |
| storageClass  |         | scalars  | Specify the "storageClass" used to provision the volume. Or the default. StorageClass will be used(the default). Set it to "-" to disable dynamic provisioning. |
| accessMode    |         | scalars  | Kubernetes Persistent Volume Access Modes: <br />ReadWriteOnce<br />ReadOnlyMany <br />ReadWriteMany. |
| size          |         | scalars  | Match the volume size of PVC.                                |

### servingRedis mappings

Configuration of kubernetes deployment redis.

| Name         | Type     | Description                                                  |
| ------------ | -------- | ------------------------------------------------------------ |
| password     | scalars  | Kubernetes ServiceTypes, default is NodePort.<br />Other modules can connect to the fateflow. |
| nodeSelector | mappings | kubernetes nodeSelector.                                     |
| subPath       | scalars  | Path of data persistence, specify the "subPath" if the PVC is shared with other components. |
| existingClaim | scalars  | Use the existing PVC which must be created manually before bound. |
| storageClass  | scalars  | Specify the "storageClass" used to provision the volume. Or the default. StorageClass will be used(the default). Set it to "-" to disable dynamic provisioning. |
| accessMode    | scalars  | Kubernetes Persistent Volume Access Modes: <br />ReadWriteOnce<br />ReadOnlyMany <br />ReadWriteMany. |
| size          | scalars  | Match the volume size of PVC.                                |
