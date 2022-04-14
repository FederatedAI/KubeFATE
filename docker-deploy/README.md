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

***If you have deployed other versions of FATE before, please delete and clean up before deploying the new version, [Deleting the cluster](#deleting-the-cluster).***
### Setting up a local registry Harbor (Optional)
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

**NOTE:** For Chinese user who has difficulty to access docker hub, you can set `RegistryURI` to `hub.c.163.com` to use the mirror of the registry within China.


### Configuring multiple parties of FATE
There are usually multiple parties participating a federated training. Each party should install FATE using a set of configuration files and scripts.

The following steps illustrate how to generate necessary configuration files and deploy two parties on different hosts.

Before deploying the FATE system, multiple parties should be defined in the configuration file: `docker-deploy/parties.conf`.

In the following sample of `docker-deploy/parties.conf` , two parities are specified by id as `10000` and `9999`. Their cluster are going to be deployed on hosts with IP addresses of *192.168.7.1* and *192.168.7.2*. By default, to save time for downloading images, KubeFATE will use images without neural network dependencies, set the `enabled_nn` to `true` in "parties.conf" if neural network workflow is required.

```bash
user=fate
dir=/data/projects/fate
partylist=(10000 9999)
partyiplist=(192.168.7.1 192.168.7.2)
servingiplist=(192.168.7.1 192.168.7.2)
exchangeip=
# computing_backend could be eggroll or spark.
computing_backend=eggroll

# true if you need python-nn else false, the default value will be false
enabled_nn=false

fateboard_username=admin                    # Username to access fateboard
fateboard_password=admin                    # Password to access fateboard
```

Spark was introduced in FATE v1.5 as the underlying computing backend, for more details
about FATE on Spark please refer to this [document](../docs/FATE_On_Spark.md).

On the host running FATE, the non-root user needs the owner permission of `/data/projects/fate` folder and Docker permission. No other action is required if the user is root.

```bash
# Create a fate user whose group is docker
[user@localhost]$ sudo useradd -s /bin/bash -g docker -d /home/fate fate
# Set user password
[user@localhost]$ sudo passwd fate
# Create docker-compose deployment directory
[user@localhost]$ sudo mkdir -p /data/projects/fate
# Modify the corresponding users and groups of docker-compose deployment directory
[user@localhost]$ sudo chown -R fate:docker /data/projects/fate
# Select users
[user@localhost]$ sudo su fate
# Check whether you have docker permission
[fate@localhost]$ docker ps
CONTAINER ID  IMAGE   COMMAND   CREATED   STATUS    PORTS   NAMES
# View docker-compose deployment directory
[fate@localhost]$ ls -l /data/projects/
total 0
drwxr-xr-x. 2 fate docker 6 May 27 00:51 fate
```

By default, the exchange service is not deployed. The exchange service runs on port 9371. If an exchange (co-locates on the host of the same party or runs standalone) service is needed, update the value of `exchangeip` to the IP address of the desired host.

After editting the above configuration file, use the following commands to generate configuration of target hosts.  

```bash
$ cd docker-deploy
$ ./generate_config.sh
```

Now, tar files have been generated for each party including the exchange node (party). They are named as ```confs-<party-id>.tar ``` and ```serving-<party-id>.tar```.

### Deploying FATE to target hosts

**Note:** Before running the below commands, all target hosts must

* allow password-less SSH access with SSH key;
* meet the requirements specified in [Prerequisites](#Prerequisites).

To deploy FATE to all configured target hosts, use the below command:
```bash
$ ./docker_deploy.sh all
```

The script copies tar files (e.g. `confs-<party-id>.tar` or `serving-<party-id>.tar`) to corresponding target hosts. It then launches a FATE cluster on each host using `docker-compose` commands.

By default, the script starts the training and serving cluster simultaneously. If you need to start them separately, add the `--training` or `--serving` to the `docker_deploy.sh` as follows.

(Optional) To deploy all parties training cluster, use the below command:
```bash
$ ./docker_deploy.sh all --training
```

(Optional) To deploy all parties serving cluster, use the below command:
```bash
$ ./docker_deploy.sh all --serving
```

(Optional) To deploy FATE to a single target host, use the below command with the party's id (10000 in the below example):
```bash
$ ./docker_deploy.sh 10000
```

(Optional) To deploy the exchange node to a target host, use the below command:
```bash
$ ./docker_deploy.sh exchange
```


Once the commands finish, log in to any host and use `docker ps` to verify the status of the cluster. A sample output is as follows:

```
CONTAINER ID        IMAGE                                        COMMAND                  CREATED              STATUS              PORTS                                 NAMES
69b8b36af395        federatedai/eggroll:<tag>          "bash -c 'java -Dlog…"   2 hours ago         Up 2 hours    
      0.0.0.0:9371->9370/tcp                                                   confs-exchange_exchange_1
71cd792ba088        federatedai/serving-proxy:<tag>    "/bin/sh -c 'java -D…"   2 hours ago         Up 2 hours    
      0.0.0.0:8059->8059/tcp, 0.0.0.0:8869->8869/tcp, 8879/tcp                 serving-10000_serving-proxy_1
2c79047918c6        federatedai/serving-server:<tag>   "/bin/sh -c 'java -c…"   2 hours ago         Up 2 hours    
      0.0.0.0:8000->8000/tcp                                                   serving-10000_serving-server_1
b1a5384a55dc        redis:5                            "docker-entrypoint.s…"   2 hours ago         Up 2 hours    
      6379/tcp                                                                 serving-10000_redis_1
321c4e29313b        federatedai/client:<tag>           "/bin/sh -c 'sleep 5…"   2 hours ago         Up 2 hours    
      0.0.0.0:20000->20000/tcp                                                 confs-10000_client_1
c1b3190126ab        federatedai/fateboard:<tag>        "/bin/sh -c 'java -D…"   2 hours ago         Up 2 hours    
      0.0.0.0:8080->8080/tcp                                                   confs-10000_fateboard_1
cc679996e79f        federatedai/python:<tag>           "/bin/sh -c 'sleep 5…"   2 hours ago         Up 2 hours    
      0.0.0.0:8484->8484/tcp, 0.0.0.0:9360->9360/tcp, 0.0.0.0:9380->9380/tcp   confs-10000_python_1
c79800300000        federatedai/eggroll:<tag>          "bash -c 'java -Dlog…"   2 hours ago         Up 2 hours    
      4671/tcp                                                                 confs-10000_nodemanager_1
ee2f1c3aad99        federatedai/eggroll:<tag>          "bash -c 'java -Dlog…"   2 hours ago         Up 2 hours    
      4670/tcp                                                                 confs-10000_clustermanager_1
a1f784882d20        federatedai/eggroll:<tag>          "bash -c 'java -Dlog…"   2 hours ago         Up 2 hours                  0.0.0.0:9370->9370/tcp                                                   confs-10000_rollsite_1
2b4526e6d534        mysql:8                            "docker-entrypoint.s…"   2 hours ago         Up 2 hours                  3306/tcp, 33060/tcp                                                      confs-10000_mysql_1
```

### Verifying the deployment
On the target node of each party, a container named  `confs-<party_id>_python_1` should have been created and running the `fate-flow` service. For example, on Party 10000's node, run the following commands to verify the deployment:
```bash
$ docker exec -it confs-10000_client_1 bash
$ flow test toy --guest-party-id 10000 --host-party-id 9999
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

### Verifying the serving service
#### Steps on the host
##### Logging in to the python container
```bash
docker exec -it confs-10000_client_1 bash
```

##### Modifying examples/upload_host.json 
```bash
cat > fateflow/examples/upload/upload_host.json <<EOF
{
  "file": "examples/data/breast_hetero_host.csv",
  "id_delimiter": ",",
  "head": 1,
  "partition": 10,
  "namespace": "experiment",
  "table_name": "breast_hetero_host"
}
EOF
```

##### Uploading data
```bash
flow data upload -c fateflow/examples/upload/upload_host.json
```

#### Steps on the guest
##### Getting in to the python container
```bash
docker exec -it confs-9999_client_1 bash
```

##### Modifying examples/upload_guest.json 

```bash
cat > fateflow/examples/upload/upload_guest.json <<EOF
{
  "file": "examples/data/breast_hetero_guest.csv",
  "id_delimiter": ",",
  "head": 1,
  "partition": 4,
  "namespace": "experiment",
  "table_name": "breast_hetero_guest"
}
EOF
```

##### Uploading data

```bash
flow data upload -c fateflow/examples/upload/upload_guest.json
```

##### Modifying examples/test_hetero_lr_job_conf.json

```bash
cat > fateflow/examples/lr/test_hetero_lr_job_conf.json <<EOF
{
  "dsl_version": "2",
  "initiator": {
    "role": "guest",
    "party_id": 9999
  },
  "role": {
    "guest": [
      9999
    ],
    "host": [
      10000
    ],
    "arbiter": [
      10000
    ]
  },
  "job_parameters": {
    "common": {
      "task_parallelism": 2,
      "computing_partitions": 8,
      "task_cores": 4,
      "auto_retries": 1
    }
  },
  "component_parameters": {
    "common": {
      "intersection_0": {
        "intersect_method": "raw",
        "sync_intersect_ids": true,
        "only_output_key": false
      },
      "hetero_lr_0": {
        "penalty": "L2",
        "optimizer": "rmsprop",
        "alpha": 0.01,
        "max_iter": 3,
        "batch_size": 320,
        "learning_rate": 0.15,
        "init_param": {
          "init_method": "random_uniform"
        }
      }
    },
    "role": {
      "guest": {
        "0": {
          "reader_0": {
            "table": {
              "name": "breast_hetero_guest",
              "namespace": "experiment"
            }
          },
          "dataio_0": {
            "with_label": true,
            "label_name": "y",
            "label_type": "int",
            "output_format": "dense"
          }
        }
      },
      "host": {
        "0": {
          "reader_0": {
            "table": {
              "name": "breast_hetero_host",
              "namespace": "experiment"
            }
          },
          "dataio_0": {
            "with_label": false,
            "output_format": "dense"
          },
          "evaluation_0": {
            "need_run": false
          }
        }
      }
    }
  }
}
EOF
```

##### Modifying examples/test_hetero_lr_job_dsl.json

```bash
cat > fateflow/examples/lr/test_hetero_lr_job_dsl.json <<EOF
{
  "components": {
    "reader_0": {
      "module": "Reader",
      "output": {
        "data": [
          "table"
        ]
      }
    },
    "dataio_0": {
      "module": "DataIO",
      "input": {
        "data": {
          "data": [
            "reader_0.table"
          ]
        }
      },
      "output": {
        "data": [
          "train"
        ],
        "model": [
          "dataio"
        ]
      },
      "need_deploy": true
    },
    "intersection_0": {
      "module": "Intersection",
      "input": {
        "data": {
          "data": [
            "dataio_0.train"
          ]
        }
      },
      "output": {
        "data": [
          "train"
        ]
      }
    },
    "hetero_feature_binning_0": {
      "module": "HeteroFeatureBinning",
      "input": {
        "data": {
          "data": [
            "intersection_0.train"
          ]
        }
      },
      "output": {
        "data": [
          "train"
        ],
        "model": [
          "hetero_feature_binning"
        ]
      }
    },
    "hetero_feature_selection_0": {
      "module": "HeteroFeatureSelection",
      "input": {
        "data": {
          "data": [
            "hetero_feature_binning_0.train"
          ]
        },
        "isometric_model": [
          "hetero_feature_binning_0.hetero_feature_binning"
        ]
      },
      "output": {
        "data": [
          "train"
        ],
        "model": [
          "selected"
        ]
      }
    },
    "hetero_lr_0": {
      "module": "HeteroLR",
      "input": {
        "data": {
          "train_data": [
            "hetero_feature_selection_0.train"
          ]
        }
      },
      "output": {
        "data": [
          "train"
        ],
        "model": [
          "hetero_lr"
        ]
      }
    },
    "evaluation_0": {
      "module": "Evaluation",
      "input": {
        "data": {
          "data": [
            "hetero_lr_0.train"
          ]
        }
      },
      "output": {
        "data": [
          "evaluate"
        ]
      }
    }
  }
}
EOF
```

##### Submitting a job
```bash
flow job submit -d fateflow/examples/lr/test_hetero_lr_job_dsl.json -c fateflow/examples/lr/test_hetero_lr_job_conf.json
```

output：
```json
{
    "data": {
        "board_url": "http://fateboard:8080/index.html#/dashboard?job_id=202111230933232084530&role=guest&party_id=9999",
        "code": 0,
        "dsl_path": "/data/projects/fate/fate_flow/jobs/202111230933232084530/job_dsl.json",
        "job_id": "202111230933232084530",
        "logs_directory": "/data/projects/fate/fate_flow/logs/202111230933232084530",
        "message": "success",
        "model_info": {
            "model_id": "arbiter-10000#guest-9999#host-10000#model",
            "model_version": "202111230933232084530"
        },
        "pipeline_dsl_path": "/data/projects/fate/fate_flow/jobs/202111230933232084530/pipeline_dsl.json",
        "runtime_conf_on_party_path": "/data/projects/fate/fate_flow/jobs/202111230933232084530/guest/9999/job_runtime_on_party_conf.json",
        "runtime_conf_path": "/data/projects/fate/fate_flow/jobs/202111230933232084530/job_runtime_conf.json",
        "train_runtime_conf_path": "/data/projects/fate/fate_flow/jobs/202111230933232084530/train_runtime_conf.json"
    },
    "jobId": "202111230933232084530",
    "retcode": 0,
    "retmsg": "success"
}
```

##### Checking status of training jobs
```bash
flow task query -r guest -j 202111230933232084530 | grep -w f_status
```

output:
```json
            "f_status": "success",
            "f_status": "waiting",
            "f_status": "running",
            "f_status": "waiting",
            "f_status": "waiting",
            "f_status": "success",
            "f_status": "success",
```

Wait for all waiting states to change to success.

##### Deploy model

```bash
flow model deploy --model-id arbiter-10000#guest-9999#host-10000#model --model-version 202111230933232084530
```

```json
{
    "data": {
        "arbiter": {
            "10000": 0
        },
        "detail": {
            "arbiter": {
                "10000": {
                    "retcode": 0,
                    "retmsg": "deploy model of role arbiter 10000 success"
                }
            },
            "guest": {
                "9999": {
                    "retcode": 0,
                    "retmsg": "deploy model of role guest 9999 success"
                }
            },
            "host": {
                "10000": {
                    "retcode": 0,
                    "retmsg": "deploy model of role host 10000 success"
                }
            }
        },
        "guest": {
            "9999": 0
        },
        "host": {
            "10000": 0
        },
        "model_id": "arbiter-10000#guest-9999#host-10000#model",
        "model_version": "202111230954255210490"
    },
    "retcode": 0,
    "retmsg": "success"
}
```

*The `model_version` that needs to be used later are all obtained in this step `"model_version": "202111230954255210490"`*

##### Modifying the configuration of loading model
```bash
cat > fateflow/examples/model/publish_load_model.json <<EOF
{
  "initiator": {
    "party_id": "9999",
    "role": "guest"
  },
  "role": {
    "guest": [
      "9999"
    ],
    "host": [
      "10000"
    ],
    "arbiter": [
      "10000"
    ]
  },
  "job_parameters": {
    "model_id": "arbiter-10000#guest-9999#host-10000#model",
    "model_version": "202111230954255210490"
  }
}
EOF
```


##### Loading a model
```bash 
flow model load -c fateflow/examples/model/publish_load_model.json
```

output:
```json
{
    "data": {
        "detail": {
            "guest": {
                "9999": {
                    "retcode": 0,
                    "retmsg": "success"
                }
            },
            "host": {
                "10000": {
                    "retcode": 0,
                    "retmsg": "success"
                }
            }
        },
        "guest": {
            "9999": 0
        },
        "host": {
            "10000": 0
        }
    },
    "jobId": "202111240844337394000",
    "retcode": 0,
    "retmsg": "success"
}
```

##### Modifying the configuration of binding model

```bash
cat > fateflow/examples/model/bind_model_service.json <<EOF
{
    "service_id": "test",
    "initiator": {
        "party_id": "9999",
        "role": "guest"
    },
    "role": {
        "guest": ["9999"],
        "host": ["10000"],
        "arbiter": ["10000"]
    },
    "job_parameters": {
        "work_mode": 1,
        "model_id": "arbiter-10000#guest-9999#host-10000#model",
        "model_version": "202111230954255210490"
    }
}
EOF
```


##### Binding a model
```bash
flow model bind -c fateflow/examples/model/bind_model_service.json
```

output:
```json
{
    "retcode": 0,
    "retmsg": "service id is test"
}
```

##### Testing online serving
Send the following message to serving interface "{SERVING_SERVICE_IP}:8059/federation/v1/inference" of the "GUEST" party.

```bash
$ curl -X POST -H 'Content-Type: application/json' -i 'http://192.168.7.2:8059/federation/v1/inference' --data '{
  "head": {
    "serviceId": "test"
  },
  "body": {
    "featureData": {
        "x0": 1.88669,
        "x1": -1.359293,
        "x2": 2.303601,
        "x3": 2.00137,
        "x4": 1.307686
    },
    "sendToRemoteFeatureData": {
        "phone_num": "122222222"
    }
  }
}'
```

output:
```json
{"retcode":0,"retmsg":"","data":{"score":0.018025086161221948,"modelId":"guest#9999#arbiter-10000#guest-9999#host-10000#model","modelVersion":"202111240318516571130","timestamp":1637743473990},"flag":0}
```

### Deleting the cluster
Use this command to stop all cluster:
```
./docker_deploy.sh --delete all
```

To delete the cluster completely, log in to each host and run the commands as follows:
```bash
$ cd /data/projects/fate/confs-<id>/  # id of party
$ docker-compose down
$ rm -rf ../confs-<id>/               # delete the legacy files
```
