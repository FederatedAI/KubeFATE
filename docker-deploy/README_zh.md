## 使用Docker Compose 部署 FATE

### 前言

[FATE](https://www.fedai.org/ )是一个联邦学习框架，能有效帮助多个机构在满足用户隐私保护、数据安全和政府法规的要求下，进行数据使用和建模。项目地址：（https://github.com/FederatedAI/FATE/） 本文档介绍使用Docker Compose部署FATE集群的方法。

### Docker Compose 简介

Compose是用于定义和运行多容器Docker应用程序的工具。通过Compose，您可以使用YAML文件来配置应用程序的服务。然后，使用一个命令，就可以从配置中创建并启动所有服务。要了解有关Compose的所有功能的更多信息，请参阅[相关文档](https://docs.docker.com/compose/#features)。

使用Docker compose 可以方便的部署FATE，下面是使用步骤。

### 目标

两个可以互通的FATE实例，每个实例均包括FATE所有组件。

### 准备工作

1. 两个主机（物理机或者虚拟机，都是Centos7系统）；
2. 所有主机安装Docker 版本 : 18+；
3. 所有主机安装Docker-Compose 版本: 1.24+；
4. 部署机可以联网，所以主机相互之间可以网络互通；
5. 运行机已经下载FATE 的各组件镜像（离线构建镜像参考文档[构建镜像](https://github.com/FederatedAI/FATE/tree/master/docker-build)）。

### 下载部署脚本

在任意机器上下载合适的KubeFATE版本，可参考 [releases pages](https://github.com/FederatedAI/KubeFATE/releases)，然后解压。

### 修改镜像配置文件（可选）

在默认情况下，脚本在部署期间会从 [Docker Hub](https://hub.docker.com/search?q=federatedai&type=image)中下载镜像。

对于中国的用户可以用使用国内镜像源：
具体方法是通过编辑docker-deploy目录下的.env文件,给`RegistryURI`参数填入以下字段

```bash
RegistryURI=hub.c.163.com
```

如果在运行机器上已经下载或导入了所需镜像，部署将会变得非常容易。

### 手动下载镜像（可选）

如果运行机没有FATE组件的镜像，可以通过以下命令从Docker Hub获取镜像。FATE镜像的版本`<version>`可在[release页面](https://github.com/FederatedAI/FATE/releases)上查看，其中serving镜像的版本信息在[这个页面](https://github.com/FederatedAI/FATE-Serving/releases)：

```bash
$ docker pull federatedai/eggroll:<version>-release
$ docker pull federatedai/fateboard:<version>-release
$ docker pull federatedai/python:<version>-release
$ docker pull federatedai/serving-server:<version>-release
$ docker pull federatedai/serving-proxy:<version>-release
$ docker pull redis:5
$ docker pull mysql:8
```

检查所有镜像是否下载成功。

```bash
$ docker images
REPOSITORY                         TAG 
federatedai/eggroll                <version>-release
federatedai/fateboard              <version>-release
federatedai/python                 <version>-release
federatedai/client                 <version>-release
federatedai/serving-server         <version>-release
federatedai/serving-proxy          <version>-release
redis                              5
mysql                              8
```

### 离线部署（可选）

当我们的运行机器处于无法连接外部网络的时候，就无法从Docker Hub下载镜像，建议使用[Harbor](https://goharbor.io/)作为本地镜像仓库。安装Harbor请参考[文档](https://github.com/FederatedAI/KubeFATE/blob/master/registry/install_harbor.md)。在`.env`文件中，将`RegistryURI`变量更改为Harbor的IP。如下面 192.168.10.1是Harbor IP的示例。

```bash
$ cd KubeFATE/
$ vi .env

...
RegistryURI=192.168.10.1/federatedai
...
```

### 用Docker Compose部署FATE

  ***如果在之前你已经部署过其他版本的FATE，请删除清理之后再部署新的版本，[删除部署](#删除部署).***

#### 配置需要部署的实例数目

部署脚本提供了部署多个FATE实例的功能，下面的例子我们部署在两个机器上，每个机器运行一个FATE实例，这里两台机器的IP分别为*192.168.7.1*和*192.168.7.2*

根据需求修改配置文件`kubeFATE\docker-deploy\parties.conf`。

下面是修改好的文件，`party 10000`的集群将部署在*192.168.7.1*上，而`party 9999`的集群将部署在*192.168.7.2*上。为了减少所需拉取镜像的大小，KubeFATE在默认情况下，会使用不带神经网络的“python”容器，若需要跑神经网络的算法则需把“parties.conf”中的`enabled_nn`设置成`true`。

```
user=fate                                   # 运行FATE容器的用户
dir=/data/projects/fate                     # docker-compose部署目录
partylist=(10000 9999)                      # 组织id
partyiplist=(192.168.7.1 192.168.7.2)       # id对应训练集群ip
servingiplist=(192.168.7.1 192.168.7.2)     # id对应在线预测集群ip
# computing_backend could be eggroll or spark.
computing_backend=eggroll

# true if you need python-nn else false, the default value will be false
enabled_nn=false

fateboard_username=admin                    # 访问fateboard的用户名
fateboard_password=admin                    # 访问fateboard的密码
```

FATE 1.5 支持使用Spark作为底层的分布式计算引擎，Spark集群默认会通过容器的方式部署。相关的简介可以参考这个[链接](../docs/FATE_On_Spark.md).

**注意**: 默认情况下不会部署exchange组件。如需部署，用户可以把服务器IP填入上述配置文件的`exchangeip`中，该组件的默认监听端口为9371

在运行部署脚本之前，需要确保部署机器可以ssh免密登录到两个运行节点主机上。user代表免密的用户。

在运行FATE的主机上，user是非root用户的，需要有`/data/projects/fate`文件夹权限和docker权限。如果是root用户则不需要任何其他操作。

```bash
# 创建一个组为docker的fate用户
[user@localhost]$ sudo useradd -s /bin/bash -g docker -d /home/fate fate
# 设置用户密码
[user@localhost]$ sudo passwd fate
# 创建docker-compose部署目录
[user@localhost]$ sudo mkdir -p /data/projects/fate
# 修改docker-compose部署目录对应用户和组
[user@localhost]$ sudo chown -R fate:docker /data/projects/fate
# 选择用户
[user@localhost]$ sudo su fate
# 查看是否拥有docker权限
[fate@localhost]$ docker ps
CONTAINER ID  IMAGE   COMMAND   CREATED   STATUS    PORTS   NAMES
# 查看docker-compose部署目录
[fate@localhost]$ ls -l /data/projects/
total 0
drwxr-xr-x. 2 fate docker 6 May 27 00:51 fate
```

#### 执行部署脚本
以下修改可在任意机器执行。

进入目录`kubeFATE\docker-deploy`，然后运行：

```bash
$ ./generate_config.sh          # 生成部署文件
$ ./docker_deploy.sh all        # 在各个party上部署FATE
```
脚本将会生成10000、9999两个组织(Party)的部署文件，然后打包成tar文件。接着把tar文件`confs-<party-id>.tar`、`serving-<party-id>.tar`分别复制到party对应的主机上并解包，解包后的文件默认在`/data/projects/fate`目录下。然后脚本将远程登录到这些主机并使用docker compose命令启动FATE实例。

命令成功执行返回后，登录其中任意一个主机：

```bash
$ ssh root@192.168.7.1
```

使用以下命令验证实例状态，

```bash
$ docker ps
````
输出显示如下，若各个组件都是运行（up）状态，说明部署成功。

```
CONTAINER ID        IMAGE                                     COMMAND                  CREATED             STATUS              PORTS                                 NAMES
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

####  验证部署

docker-compose上的FATE启动成功之后需要验证各个服务是否都正常运行，我们可以通过验证toy_example示例来检测。

选择192.168.7.1这个节点验证，使用以下命令验证：

```bash
#在192.168.7.1上执行下列命令
$ docker exec -it confs-10000_client_1 bash                        #进入python组件容器内部
$ flow test toy --guest-party-id 10000 --host-party-id 9999        #验证
```

如果测试通过，屏幕将显示类似如下消息：

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

有关测试结果的更多详细信息，请参阅"python/examples/toy_example/README.md"这个文件 。

#### 验证Serving-Service功能
##### Host方操作
###### 进入python容器
```bash
docker exec -it confs-10000_client_1 bash
```

###### 修改examples/upload_host.json 
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

###### 上传数据
```bash
flow data upload -c fateflow/examples/upload/upload_host.json
```

##### Guest方操作
###### 进入python容器
```bash
docker exec -it confs-9999_client_1 bash
```

###### 修改examples/upload_guest.json 
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

###### 上传数据
```bash
flow data upload -c fateflow/examples/upload/upload_guest.json
```

###### 修改examples/test_hetero_lr_job_conf.json

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

##### 修改examples/test_hetero_lr_job_dsl.json

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

###### 提交任务
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

###### 查看训练任务状态
```bash
flow task query -r guest -j 202111230933232084530 | grep -w f_status
```

output:
```bash
            "f_status": "success",
            "f_status": "waiting",
            "f_status": "running",
            "f_status": "waiting",
            "f_status": "waiting",
            "f_status": "success",
            "f_status": "success",
```

等到所有的`waiting`状态变为`success`.

###### 部署模型

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

*后面需要用到的`model_version`都是这一步得到的`"model_version": "202111230954255210490"`*

###### 修改加载模型的配置
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

###### 加载模型
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

###### 修改绑定模型的配置
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


###### 绑定模型
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

###### 在线测试
发送以下信息到"GUEST"方的推理服务"{SERVING_SERVICE_IP}:8059/federation/v1/inference"

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
### 删除部署
在部署机器上运行以下命令可以停止所有FATE集群：
```bash
./docker_deploy.sh --delete all
```

如果想要彻底删除在运行机器上部署的FATE，可以分别登录节点，然后运行命令：

```bash
$ cd /data/projects/fate/confs-<id>/  # <id> 组织的id，本例中代表10000或者9999
$ docker-compose down
$ rm -rf ../confs-<id>/               # 删除docker-compose部署文件
```

### 可能遇到的问题

#### python容器退出

```bash
$ docker exec -it confs-10000_python_1 bash
```

进入docker容器后马上又弹出来了。

解决办法：稍等一会再尝试。

因为python服务依赖其他所有服务的正常运行，然而第一次启动的时候MySQL需要初始化数据库，python服务的容器会出现几次重启，当MySQL等其他服务都运行正常之后，就可以正常执行了。

#### 采用docker hub下载镜像速度可能较慢。

解决办法：可以自己构建镜像，自己构建镜像参考[这里](https://github.com/FederatedAI/FATE/tree/master/docker-build)。

#### 运行脚本`./docker_deploy.sh all`的时候提示需要输入密码

解决办法：检查免密登陆是否正常。ps:直接输入对应主机的用户密码也可以继续运行。

#### CPU指令集问题

解决办法：查看[wiki](https://github.com/FederatedAI/KubeFATE/wiki/KubeFATE)页面的storage-service部分
