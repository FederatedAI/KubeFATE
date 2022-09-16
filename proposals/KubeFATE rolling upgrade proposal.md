# Modification History

| Date       | Content                                                 | person    |
|------------|---------------------------------------------------------|-----------|
| 2022/07/01 | initial version                                         | Chen Jing |
| 2022/07/06 | Mainly change the upgrading matrix of using SC scenario | Chen Jing |


# The Goal
The goal of this feature is to make the upgrade smooth without any data loss.

1. Reuse the command "kubefate cluster update -f <cluster.yaml>", make it able to detect the version change and do the rolling upgrade.
2. Make sure there is no data loss, every persistence volume of an older version FATE can be reclaimed by the newer FATE.
3. Make sure when rolling upgrade is failed, everything can be back to the origin status.

# Functional spec
## Upgrade
### If the user doesn't enable persistence
Then definitely there will be data loss, even when a pod restarts.

### If the user is using storage class to provision persistence volumes
We cannot support 100% automatic rolling upgrade for the older versions when we upgrade them to v1.9.0+ versions. Because in the version v1.9.0, we will change some FATE components to statefulSet, and it cannot re-attach the old pvc/pv because they used to be deployments.


| From/To | v1.7.1     | v1.7.2     | v1.8.0     | v1.9.0  | v1.10.0+   |
|---------|------------|------------|------------|---------|------------|
| v1.7.0  | no pv loss | no pv loss | no pv loss | pv loss | pv loss    |
| v1.7.1  |            | no pv loss | no pv loss | pv loss | pv loss    |
| v1.7.2  |            |            | no pv loss | pv loss | pv loss    |
| v1.8.0  |            |            |            | pv loss | pv loss    |
| v1.9.0  |            |            |            |         | no pv loss |

#### Workaround
* Manually relocate the old files: If the storage facility can help to archive the old files, then the user can manually move the old files to the directories of the newer PV.
* Turn to use self-managed PVC: If the storage class has the "Retain" policy, then the old PVs will be preserved. The user can manage a PVC based on the old PV, and do not leverage storage class to automatically create PVC after upgrading.

#### The influenced data of pv loss
* Mysql:
The metadata of the FATE jobs.

* Jupyter Notebook client:
The new created files under the "persistence".

* Node manager:
The uploaded dataset and the intermediate results.

* HDFS:
The uploaded dataset and the intermediate results.

### If the user is using existing PVC

| From/To | v1.7.1     | v1.7.2     | v1.8.0     | v1.9.0     | v1.10.0+   |
|---------|------------|------------|------------|------------|------------|
| v1.7.0  | no pv loss | no pv loss | no pv loss | no pv loss | no pv loss |
| v1.7.1  |            | no pv loss | no pv loss | no pv loss | no pv loss |
| v1.7.2  |            |            | no pv loss | no pv loss | no pv loss |
| v1.8.0  |            |            |            | no pv loss | no pv loss |
| v1.9.0  |            |            |            |            | no pv loss |

If the user is using self-managed PVC, then no data will be loss.

## Rollback
If the upgrade failed, then we make sure to reproduce everything before the upgrade happens.

# User story
Suppose the user's name is Tom.

1. Tom has a cluster.yaml which declares a v1.7.0 FATE cluster, in which each component has been configured an existing PV, the FATE cluster has been up for a long time.
2. Tom changes the chart version and image version in the cluster.yaml file to v1.8.0.
3. Tom executes "kubefate cluster update -f cluster.yaml" in the terminal.
4. Tom sees that all the old pods turn into a "Terminating" status, and the new pods starting to pop out as "Running"
5. The old PVs and PVCs will be re-attached to the newer version pods automatically.
6. The sql update scripts will be executed against the mysql database to update the schema, based on the version diff. For example, we have 3 .sql scripts which are named as: 170-171.sql, 171-172.sql and 172-180.sql, then the three scripts will be executed against the mysql database in sequence when the upgrade is from v1.7.0 to v1.8.0. This step will only be conducted after the cluster turns into the "running" status.
7. If the upgrade fails, which means finally the cluster cannot turn into a "running" status, we need to fall back to the previous version. Kubefate service will help to fetch the previous cluster.yaml from the database and re-install it. As the pv/pvc doesn't change, and the mysql upgrade script hasn't been executed at this stage, nothing more need to be done.

# Implementation plan
Plan to use the golang sdk of helm, and to use the helm upgrade functionality. It has already supported pv/pvc auto-reattach and rollback. We just need to make sure we run the mysql script against the mysql container after helm upgrade has been done.

The mysql script is at: https://github.com/FederatedAI/FATE/tree/develop-1.9.0/deploy/upgrade/sql