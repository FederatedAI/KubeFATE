# FATE cluster configuration
`cluster.yaml` declares information about the FATE cluster to be deployed, which KubeFATE CLI uses to deploy the FATE cluster.

## cluster.yaml
|Name      |Type     |Description                                        |
|----------|---------|---------------------------------------------------|
|name      |scalars  |FATE cluster name.                                 |
|namespace |scalars  |Kubernetes namespace for FATE cluster.             |
|version   |scalars  |FATE cluster version.                              |
|partyId   |scalars  |FATE cluster party id.                             |
|modules   |sequences|Modules to be deployed in the FATE cluster.        |
|proxy     |mappings |Configuration of FATE cluster `proxy` module.      |
|egg       |mappings |Configuration of FATE cluster `egg` module.        |

### egg mappings
|Name      |Type     |Description                                        |
|----------|---------|---------------------------------------------------|
|count     |scalars  |Number of FATE cluster `egg` modules.              |

### proxy mappings
It is used to declare the `proxy` module in the FATE cluster to be deployed.

|Name      |subitem   |Type     |Description                                        |
|----------|----------|---------|---------------------------------------------------|
|type      |          |scalars  |Kubernetes ServiceTypes, default is NodePort.      |
|nodePort  |          |scalars  |The port used by `proxy` module's kubernetes service, default range: 30000-32767. |
|partyList |          |sequences|If this FATE cluster is exchange cluster, partyList is all party's sequences of all parties proxy address. If this FATE cluster is one of participants, delete this configuration item. |
|partyList |partyId   |scalars  |Participant FATE cluster party ID.                 |
|partyList |partyIp   |scalars  |Participant FATE cluster IP.                       |
|partyList |partyPort |scalars  |Participant FATE cluster port.                     |
|exchange  |          |mappings |FATE cluster `exchange` module's ip and port.      |
|exchange  |ip        |mappings |FATE cluster `exchange` module's ip. .             |
|exchange  |port      |mappings |FATE cluster `exchange` module's port.             |

FATE cluster has two deployment modes: with exchange and without exchange.
#### Exchange mode:
Every party connected to the `exchange`, which has the proxy addresses of all parties.
- For `exchange` cluster, only need to deploy `proxy` modules. In `proxy` configuration, no need `exchange` item, need to has the proxy addresses of all parties in `partylist`.
- For other FATE clusters, need to fill in the exchange ip and port. Can delete `partylist` configuration item.

#### Direct connection mode:
The parties are directly connected.
- No need to fill in the `exchange` ip and port.
- `partyList` needs the addresses of all other FATE clusters proxy.
