# FATE cluster configuration
`cluster.yaml` declares information about the FATE cluster to be deployed, which KubeFATE CLI uses to deploy the FATE cluster.

## cluster.yaml
| Name                  | Type      | Description                                                  |
| --------------------- | --------- | ------------------------------------------------------------ |
| name                  | scalars   | FATE cluster name.                                           |
| namespace             | scalars   | Kubernetes namespace for FATE cluster.                       |
| chartName             | scalars   | FATE chart name. (fate/fate-serving)                         |
| chartVersion          | scalars   | FATE chart corresponding version.                            |
| partyId               | scalars   | FATE cluster party id.                                       |
| registry              | scalars   | Other fate images sources.                                   |
| pullPolicy            | scalars   | kubernetes images pull policy                                |
| persistence           | bool      | mysql and nodemanager data persistence.                      |
| modules               | sequences | Modules to be deployed in the FATE cluster.                  |
| rollsite              | mappings  | Configuration of FATE cluster `rollsite` module.             |
| nodemanager           | mappings  | Configuration of FATE cluster `nodemanager` module.          |
| python                | mappings  | Configuration of FATE cluster `python` module.               |
| mysql                 | mappings  | Configuration of FATE cluster `mysql` module.<br />If you use your own redis, please delete this item. |
| externalMysqlIp       | scalars   | Access your own MySQL.                                       |
| externalMysqlPort     | scalars   | Access your own MySQL.                                       |
| externalMysqlDatabase | scalars   | Access your own MySQL.                                       |
| externalMysqlUser     | scalars   | Access your own MySQL.                                       |
| externalMysqlPassword | scalars   | Access your own MySQL.                                       |
| servingIp             | scalars   | Serving cluster connected to fate.                           |
| servingPort           | scalars   | Serving cluster connected to fate.                           |

### rollsite mappings
It is used to declare the `rollsite ` module in the FATE cluster to be deployed.

| Name         | subitem   | Type      | Description                                                  |
| ------------ | --------- | --------- | ------------------------------------------------------------ |
| type         |           | scalars   | Kubernetes ServiceTypes, default is NodePort.                |
| nodePort     |           | scalars   | The port used by `proxy` module's kubernetes service, default range: 30000-32767. |
| partyList    |           | sequences | If this FATE cluster is exchange cluster, partyList is all party's sequences of all parties proxy address. If this FATE cluster is one of participants, delete this configuration item. |
| partyList    | partyId   | scalars   | Participant FATE cluster party ID.                           |
| partyList    | partyIp   | scalars   | Participant FATE cluster IP.                                 |
| partyList    | partyPort | scalars   | Participant FATE cluster port.                               |
| exchange     |           | mappings  | FATE cluster `exchange` module's ip and port.                |
| exchange     | ip        | mappings  | FATE cluster `exchange` module's ip. .                       |
| exchange     | port      | mappings  | FATE cluster `exchange` module's port.                       |
| nodeSelector |           | mappings  | kubernetes nodeSelector.                                     |

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

| Name             | Type     | Description                                                  |
| ---------------- | -------- | ------------------------------------------------------------ |
| fateflowType     | scalars  | Kubernetes ServiceTypes, default is NodePort.<br />Other modules can connect to the fateflow |
| fateflowNodePort | scalars  | The port used by `proxy` module's kubernetes service, default range: 30000-32767. |
| nodeSelector     | mappings | kubernetes nodeSelector.                                     |

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

