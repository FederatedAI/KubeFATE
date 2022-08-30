# `partys.conf` file introduction

The partys.conf file is the main configuration file for docker-compose to deploy FATE. The meaning of each configuration will be described in detail here.

| Name | Description | default |
| --- | ---- | --- |
| user | Deploy user | fate |
| dir | Deploy PATH | /data/projects/fate |
| party_list | Deploy party list | (10000 9999) |
| party_ip_list | FATE Partys IP | (192.168.1.1 192.168.1.2) |
| serving_ip_list | FATE-Serving Partys IP | (192.168.1.1 192.168.1.2) |
| computing | Computing engine | Eggroll |
| federation | Federation engine | Eggroll |
| storage | Storage engine | Eggroll |
| algorithm | Algorithm | Basic |
| device | Device | CPU |
| compute_core | Cluster compute_core number, it is recommended to be less than the number of cpu cores | 4 |
| exchangeip | Deploy exchange cluster host IP | NULL |
| mysql_ip | External mysql IP | mysql |
| mysql_user | External mysql user | fate |
| mysql_password | External mysql password | fate_dev |
| mysql_db | External mysql database | fate_flow |
| name_node | External hdfs namenode | hdfs://namenode:9000 |
| fateboard_username | Define fateboard login information username | admin |
| fateboard_password | Define fateboard login information password | admin |
| serving_admin_username | Define serving_admin login information username | admin |
| serving_admin_password | Define serving_admin login information password | admin |
