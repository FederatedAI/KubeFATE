# FATE deployment using Docker Compose

This guide describes the process of deploying FATE using Docker Compose.

## Prerequisites
The nodes (target nodes) to install FATE must meet the following requirements:

1. A Linux host
2. Docker: 18+
3. Docker-Compose: 1.24+
4. Network connection to Internet to pull container images from Docker Hub. If network connection to Internet is not available, consider to set up [Harbor as a local registry](../registry/README.md) or use [offline images](https://github.com/FederatedAI/FATE/tree/master/docker-build). 

## Deploying FATE
A Linux host can be used as a deployment machine to run installation scripts to deploy FATE onto target hosts.

First, on a Linux host, download KubeFATE from [releases pages](https://github.com/FederatedAI/KubeFATE/releases), unzip it into folder KubeFATE.

By default, the installation script pulls the images from Docker Hub during the deployment. If the target node is not connected to Internet, refer to the below section to set up a local registry such as Harbor and use the offline images.

### Set up a local registry Harbor (Optional)
Please refer to [this guide](../registry/README.md) to install Harbor as a local registry. 

After setting up a Harbor registry, update the setting in the `.env` file. Change `RegistryURI` to the hostname or IP address of the Harbor instance. This setting lets the installation script use a local registry instead of Docker Hub.

In the below example, `192.168.10.1` is the IP address of Harbor.

```bash
$ cd KubeFATE/
$ vi .env

...

RegistryURI=192.168.10.1/federatedai

...
```

### Configuring multiple parties of FATE
There are usually multiple parties participating a federated training. Each party should install FATE using a set of configuration files and scripts. 

The following steps illustrate how to generate necessary configuration files and deploy two parties on different hosts.

Before deploying the FATE system, multiple parties should be defined in the configuration file: `docker-deploy/parties.conf`. 

In the following sample of `docker-deploy/parties.conf` , two parities are specified by id as `10000` and `9999`. They are going to be deployed on hosts with IP addresses of *192.168.7.1* and *192.168.7.2*, respectively. 

```bash
user=root
dir=/data/projects/fate
partylist=(10000 9999)
partyiplist=(192.168.7.1 192.168.7.2)
exchangeip=192.168.7.1
```
By default, the exchange node co-locates on the same host of the first party. The exchange service runs on port 9371. For this reason, the IP address of the exchange node should be the same as that of the first party. If a standalone exchange node is needed, update the value of `exchangeip` to the IP address of the desired host.

After completing the above configuration file, use the following commands to generate configuration of target hosts.  
```bash
$ cd docker-deploy
$ bash generate_config.sh
```

Now, tar files have been generated for each party including the exchange node (party). They are named as ```<party-id>-confs.tar ```.

### Deploying FATE to target hosts

**Note:** Before running the below commands, all target hosts must

* allow password-less SSH access with SSH key;
* meet the requirements specified in [Prerequisites](#Prerequisites).

To deploy FATE to all configured target hosts, use the below command:
```bash
$ bash docker_deploy.sh all
```

The script copies tar files (e.g. `10000-confs.tar` or `9999-confs.tar`) to corresponding target hosts. It then launches a FATE cluster on each host using `docker-compose` commands.


To deploy FATE to a single target host, use the below command with the party's id (10000 in the below example):
```bash
$ bash docker_deploy.sh 10000
```
To deploy the exchange node to a target host, use the below command:
```bash
$ bash docker_deploy.sh exchange
```


Once the commands finish, log in to any host and use `docker ps` to verify the status of the cluster. A sample output is as follows:

```
CONTAINER ID        IMAGE                                    COMMAND                  CREATED              STATUS              PORTS                                 NAMES
d4686d616965        federatedai/python:1.2.0-release         "/bin/bash -c 'sourc…"   About a minute ago   Up 52 seconds       9360/tcp, 9380/tcp                    confs-10000_python_1
4086ef0dc2de        federatedai/fateboard:1.2.0-release      "/bin/sh -c 'cd /dat…"   About a minute ago   Up About a minute   0.0.0.0:8080->8080/tcp                confs-10000_fateboard_1
5cf3e1f1731a        federatedai/roll:1.2.0-release           "/bin/sh -c 'cd roll…"   About a minute ago   Up About a minute   8011/tcp                              confs-10000_roll_1
11c01143540b        federatedai/meta-service:1.2.0-release   "/bin/sh -c 'java -c…"   About a minute ago   Up About a minute   8590/tcp                              confs-10000_meta-service_1
f0976f48f0f7        federatedai/proxy:1.2.0-release          "/bin/sh -c 'cd /dat…"   About a minute ago   Up About a minute   0.0.0.0:9370->9370/tcp                confs-10000_proxy_1
7354af787036        redis:5                                  "docker-entrypoint.s…"   About a minute ago   Up About a minute   6379/tcp                              confs-10000_redis_1
ed11ce8eb20d        federatedai/egg:1.2.0-release            "/bin/sh -c 'cd /dat…"   About a minute ago   Up About a minute   7778/tcp, 7888/tcp, 50001-50004/tcp   confs-10000_egg_1
6802d1e2bd21        mysql:8                                  "docker-entrypoint.s…"   About a minute ago   Up About a minute   3306/tcp, 33060/tcp                   confs-10000_mysql_1
5386bcb7565f        federatedai/federation:1.2.0-release     "/bin/sh -c 'cd /dat…"   About a minute ago   Up About a minute   9394/tcp                              confs-10000_federation_1
```

### Verifying the deployment
On the target node of each party, a container named  `confs-<party_id>_python_1` should have been created and running the `fate-flow` service. For example, on Party 10000's node, run the following commands to verify the deployment:
```bash
$ docker exec -it confs-10000_python_1 bash
$ cd /data/projects/python/examples/toy_example/
$ python run_toy_example.py 10000 9999 1
```
If the test passed, the output may look like the following:
```
"2019-08-29 07:21:25,353 - secure_add_guest.py[line:96] - INFO: begin to init parameters of secure add example guest"
"2019-08-29 07:21:25,354 - secure_add_guest.py[line:99] - INFO: begin to make guest data"
"2019-08-29 07:21:26,225 - secure_add_guest.py[line:102] - INFO: split data into two random parts"
"2019-08-29 07:21:29,140 - secure_add_guest.py[line:105] - INFO: share one random part data to host"
"2019-08-29 07:21:29,237 - secure_add_guest.py[line:108] - INFO: get share of one random part data from host"
"2019-08-29 07:21:33,073 - secure_add_guest.py[line:111] - INFO: begin to get sum of guest and host"
"2019-08-29 07:21:33,920 - secure_add_guest.py[line:114] - INFO: receive host sum from guest"
"2019-08-29 07:21:34,118 - secure_add_guest.py[line:121] - INFO: success to calculate secure_sum, it is 2000.0000000000002"
```
For more details about the testing result, please refer to `python/examples/toy_example/README.md` .

