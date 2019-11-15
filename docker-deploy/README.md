# Deployment by Docker Compose

## Prerequisites
1. Docker: 18
2. Docker-Compose: 1.24
3. [The FATE Images](https://github.com/FederatedAI/FATE/tree/docker_1.1_contribution/docker-build) have been built and downloaded by nodes.

## Deploying FATE
Use the following command to clone repo if you did not clone before:
```bash
$ git clone git@github.com:FederatedAI/KubeFATE.git
```
By default, the script pulls the images from [Docker Hub](https://hub.docker.com/search?q=federatedai&type=image) during the deployment.

### Use Third Party Registry (Optional)
It is recommended that non-Internet clusters use Harbor as a third-party registry. Please refer to [this guide](https://github.com/FederatedAI/KubeFATE/blob/master/registry/install_harbor.md) to install Harbor. Change the `RegistryURI` to [Harbor](https://goharbor.io/) hostname in the `.env` file. `192.168.10.1` is an example of Harbor IP.
```bash
$ cd KubeFATE/
$ vi .env

RegistryURI=192.168.10.1/federatedai
```

### Configure Parties
The following steps will illustrate how to deploy two parties on different hosts.

### Generate startup files
Before starting the FATE system, the user needs to define their parties in configuration file `./parties.conf`. 

The following sample of `parties.conf` defines two parities, they are party `10000` hosted on a machine *192.168.7.1* and `9999` hosted on a machine *192.168.7.2*.

```bash
user=root
dir=/data/projects/fate
partylist=(10000 9999)
partyiplist=(192.168.7.1 192.168.7.2)
exchangeip=192.168.7.1
```

**NOTE**: By default, the machine of the first party will also host the exchange on the 9371 port. A user can change the exchange IP if needed.

Use the following command to deploy each party. Before running the command, ***please make sure host 192.168.7.1 and 192.168.7.2 allow password-less SSH access with SSH key***:

```bash
$ bash generate_config.sh  # generate the config file
$ bash docker_deploy.sh    # launch the deployment
```

The script will copy "10000-confs.tar" and "9999-confs.tar" to host 192.168.7.1 and 192.168.7.2.

Afterward the script will log in to these hosts and use docker-compose command to start the FATE cluster.

Once the command returns, log in to any host and use `docker ps` to verify the status of cluster, an example output is as follows:

```
CONTAINER ID        IMAGE                                 COMMAND                  CREATED              STATUS              PORTS                                 NAMES
d4686d616965        federatedai/python:1.1-release           "/bin/bash -c 'sourc…"   About a minute ago   Up 52 seconds       9360/tcp, 9380/tcp                    confs-10000_python_1
4086ef0dc2de        federatedai/fateboard:1.1-release        "/bin/sh -c 'cd /dat…"   About a minute ago   Up About a minute   0.0.0.0:8080->8080/tcp                confs-10000_fateboard_1
5cf3e1f1731a        federatedai/roll:1.1-release             "/bin/sh -c 'cd roll…"   About a minute ago   Up About a minute   8011/tcp                              confs-10000_roll_1
11c01143540b        federatedai/meta-service:1.1-release     "/bin/sh -c 'java -c…"   About a minute ago   Up About a minute   8590/tcp                              confs-10000_meta-service_1
f0976f48f0f7        federatedai/proxy:1.1-release            "/bin/sh -c 'cd /dat…"   About a minute ago   Up About a minute   0.0.0.0:9370->9370/tcp                confs-10000_proxy_1
7354af787036        redis                                 "docker-entrypoint.s…"   About a minute ago   Up About a minute   6379/tcp                              confs-10000_redis_1
ed11ce8eb20d        federatedai/egg:1.1-release              "/bin/sh -c 'cd /dat…"   About a minute ago   Up About a minute   7778/tcp, 7888/tcp, 50001-50004/tcp   confs-10000_egg_1
6802d1e2bd21        mysql                                 "docker-entrypoint.s…"   About a minute ago   Up About a minute   3306/tcp, 33060/tcp                   confs-10000_mysql_1
5386bcb7565f        federatedai/federation:1.1-release       "/bin/sh -c 'cd /dat…"   About a minute ago   Up About a minute   9394/tcp                              confs-10000_federation_1
```

### Verify the Deployment
Since the `confs-10000_python_1` container hosts the `fate-flow` service, so we need to perform the test within that container. Use the following commands to launch:
```bash
$ docker exec -it confs-10000_python_1 bash
$ source /data/projects/python/venv/bin/activate
$ cd /data/projects/fate/python/examples/toy_example
$ python run_toy_example.py 10000 9999 1
```
If the test passed, the screen will print some messages like the follows:
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
For more details about the testing result, please refer to "/data/projects/fate/python/examples/toy_example/README.md" 
