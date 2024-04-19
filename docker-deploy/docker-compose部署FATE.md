###  docker-compose  离线部署FATE

 **请确保机器已经安装好docker, docker-compose**

#### 1.  获取fate&mysql 离线镜像， 获取部署脚本

   ~~~
   # 获取fate&mysql离线镜像
   third-party.images.tar.gz
   fate.images.tar.gz
   
   # 获取docker-deploy部署脚本
   wget https://codeload.github.com/FederatedAI/KubeFATE/zip/refs/heads/develop-2.1.0
   ~~~

#### 2. 加载镜像

   ~~~
   # 加载mysql镜像
   docker load -i third-party.images.tar.gz
   # 加载fate镜像
   docker load -i fate.images.tar.gz
   ~~~

#### 3. 部署脚本配置

   按需配置

   ~~~
   # 部署单边/双边
   vi docker_deploy/parties.conf
   
   user=fate
   dir=/data/projects/fate
   party_list=(9999)               # 部署双边增加一个party_id  party_list=(9999 10000) 
   party_ip_list=(192.168.1.1)      # 部署双边增加一个ip   party_ip_list=(192.168.1.1 192.168.1.2)
   
   # Engines:
   # Computing : Eggroll, Spark, Spark_local
   computing=Eggroll
   # Federation: OSX(computing: Eggroll/Spark/Spark_local), Pulsar/RabbitMQ(computing: Spark/Spark_local)
   federation=OSX
   # Storage: Eggroll(computing: Eggroll), HDFS(computing: Spark), LocalFS(computing: Spark_local)
   storage=Eggroll
   # Algorithm: Basic, NN, ALL
   algorithm=Basic
   # Device: CPU, IPCL, GPU
   device=CPU
   
   # spark and eggroll 
   compute_core=16
   
   # You only need to configure this parameter when you want to use the GPU, the default value is 1
   gpu_count=0
   
   # modify if you are going to use an external db
   mysql_ip=mysql
   mysql_user=fate
   mysql_password=fate_dev
   mysql_db=fate_flow
   serverTimezone=UTC
   
   name_node=hdfs://namenode:9000
   
   # Define fateboard login information
   fateboard_username=admin
   fateboard_password=admin
   ~~~

   配置  .env文件

   ~~~
   # .env 文件在docker_deploy目录下
   # 要求配置的tag名与加载到容器的tag名一致
   vi .env
   RegistryURI=""
   #SERVING_TAG=2.1.6-release
   SSH_PORT=22
   
   # PREFIX: namespace on the registry's server.
   # RegistryURI: address of the local registry
   # TAG: tag of module images.
   # SSH_PORT: port of SSH, default 22
   
   
   KubeFATE_Version=v2.1.0-release
   
   # components version
   
   FATEFlow_IMAGE="federatedai/fateflow"
   FATEFlow_IMAGE_TAG="2.1.0-release"
   FATEBoard_IMAGE="federatedai/fateboard"
   FATEBoard_IMAGE_TAG="2.1.0-release"
   MySQL_IMAGE="mysql"
   MySQL_IMAGE_TAG="8.0.28"
   Client_IMAGE="federatedai/client"
   Client_IMAGE_TAG="2.1.0-release"
   
   EGGRoll_IMAGE="federatedai/eggroll"
   EGGRoll_IMAGE_TAG="2.1.0-release"
   OSX_IMAGE="federatedai/osx"
   OSX_IMAGE_TAG="2.1.0-release"
   
   #Nginx_IMAGE="federatedai/nginx"
   #Nginx_IMAGE_TAG="2.1.0-release"
   #RabbitMQ_IMAGE="federatedai/rabbitmq"
   #RabbitMQ_IMAGE_TAG="3.8.3-management"
   #Pulsar_IMAGE="federatedai/pulsar"
   #Pulsar_IMAGE_TAG="2.10.2"
   #Hadoop_NameNode_IMAGE="federatedai/hadoop-namenode"
   #Hadoop_NameNode_IMAGE_TAG="2.1.0-hadoop3.2.1-java8"
   #Hadoop_DataNode_IMAGE="federatedai/hadoop-datanode"
   #Hadoop_DataNode_IMAGE_TAG="2.1.0-hadoop3.2.1-java8"
   #Spark_Master_IMAGE="federatedai/spark-master"
   #Spark_Master_IMAGE_TAG="2.1.0-release"
   #Spark_Worker_IMAGE="federatedai/spark-worker"
   #Spark_Worker_IMAGE_TAG="2.1.0-release"
   ~~~

   **注意：**在运行以下命令之前，所有目标主机必须允许使用 SSH 密钥进行无密码 SSH 访问（否则我们将需要为每个主机多次输入密码）。
   要将 FATE 部署到所有已配置的目标主机。

#### 4. 生成配置文件

进入目录docker-deploy，然后运行命令，生成部署文件

   ~~~
   cd docker_deploy
   bash ./generate_config.sh     # 生成部署文件
   ~~~

   ​         脚本会更具配置将会生成单边或者双边组织部署文件，然后打包成tar文件。接着把tar文件confs-<party-id>.tar分别复制到party对应的主机上并解包，解包后的文件默认在/data/projects/fate目录下。然后脚本将远程登录到这些主机并使用docker-compose命令启动FATE实例。

#### 5.  执行部署

      ~~~
      # 部署单边
      bash ./docker_deploy.sh 9999
      
      # 部署双边
      bash ./docker_deploy.sh all --training
      ~~~

#### 6.  查看实例状态

      ~~~
      cd /data/projects/fate/confs-10000
      docker-compoes ps
      ~~~

#### 7. 验证部署

      ~~~
      # 进入client组件容器内部
      docker-compose exec client bash
      # 单边toy 验证
      flow test toy -gid 9999 -hid 9999 
      # 双边toy验证
      flow test toy -gid 9999 -hid 10000
      
      # 如果成功，屏幕将显示下方类似结果
      toy test job xxxxx is success
      ~~~

#### 8.  上传数据，发起任务

   ~~~
# 上传数据（单边的， 双边需要在另一方再次执行）
from fate_client.pipeline import FateFlowPipeline
   
guest_data_path="/data/projects/fate/examples/data/breast_hetero_guest.csv"
host_data_path="/data/projects/fate/examples/data/breast_hetero_host.csv"
   
data_pipeline = FateFlowPipeline().set_parties(local="0")
guest_meta = {
       "delimiter": ",", "dtype": "float64", "label_type": "int64","label_name": "y", "match_id_name": "id"
   }
host_meta = {
       "delimiter": ",", "input_format": "dense", "match_id_name": "id"
   }
data_pipeline.transform_local_file_to_dataframe(file=guest_data_path, namespace="experiment", name="breast_hetero_guest",
                                                   meta=guest_meta, head=True, extend_sid=True)
data_pipeline.transform_local_file_to_dataframe(file=host_data_path, namespace="experiment", name="breast_hetero_host",
                                                   meta=host_meta, head=True, extend_sid=True)
   ~~~

   ~~~
   # 发起任务
from fate_client.pipeline.components.fate import (
       HeteroSecureBoost,
       Reader,
       PSI,
       Evaluation
   )
from fate_client.pipeline import FateFlowPipeline
   
   
# create pipeline for training
pipeline = FateFlowPipeline().set_parties(guest="9999", host="10000")
   
# create reader task_desc
reader_0 = Reader("reader_0")
reader_0.guest.task_parameters(namespace="experiment", name="breast_hetero_guest")
reader_0.hosts[0].task_parameters(namespace="experiment", name="breast_hetero_host")
   
# create psi component_desc
psi_0 = PSI("psi_0", input_data=reader_0.outputs["output_data"])
   
# create hetero secure_boost component_desc
hetero_secureboost_0 = HeteroSecureBoost(
       "hetero_secureboost_0", num_trees=1, max_depth=5,
       train_data=psi_0.outputs["output_data"],
       validate_data=psi_0.outputs["output_data"]
   )
   
# create evaluation component_desc
evaluation_0 = Evaluation(
       'evaluation_0', runtime_parties=dict(guest="9999"), metrics=["auc"], input_datas=[hetero_secureboost_0.outputs["train_output_data"]]
   )
   
# add training task
pipeline.add_tasks([reader_0, psi_0, hetero_secureboost_0, evaluation_0])
   
# compile and train
pipeline.compile()
pipeline.fit()
   
# print metric and model info
print (pipeline.get_task_info("hetero_secureboost_0").get_output_model())
print (pipeline.get_task_info("evaluation_0").get_output_metric())
   
# deploy task for inference
pipeline.deploy([psi_0, hetero_secureboost_0])
   
# create pipeline for predicting
predict_pipeline = FateFlowPipeline()
   
# add input to deployed_pipeline
deployed_pipeline = pipeline.get_deployed_pipeline()
reader_1 = Reader("reader_1")
reader_1.guest.task_parameters(namespace="experiment", name="breast_hetero_guest")
reader_1.hosts[0].task_parameters(namespace="experiment", name="breast_hetero_host")
deployed_pipeline.psi_0.input_data = reader_1.outputs["output_data"]
   
# add task to predict pipeline
predict_pipeline.add_tasks([reader_1, deployed_pipeline])
   
# compile and predict
predict_pipeline.compile()
predict_pipeline.predict()
   ~~~

任务成功后，屏幕将显示下方类似结果
output:

```
Job is success!!! Job id is 202404031636558952240, response_data={'apply_resource_time': 1712133417129, 'cores': 4, 'create_time': 1712133415928, 'dag': {'dag': {'conf': {'auto_retries': 0, 'computing_partitions': 8, 'cores': None, 'extra': None, 'inheritance': None, 'initiator_party_id': '9999', 'model_id': '202404031636558952240', 'model_version': '0', 'model_warehouse': {'model_id': '202404031635272687860', 'model_version': '0'}, 'priority': None, 'scheduler_party_id': '9999', 'sync_type': 'callback', 'task': None}, 'parties': [{'party_id': ['9999'], 'role': 'guest'}, {'party_id': ['10000'], 'role': 'host'}], 'party_tasks': {'guest_9999': {'conf': {}, 'parties': [{'party_id': ['9999'], 'role': 'guest'}], 'tasks': {'reader_1': {'conf': None, 'parameters': {'name': 'breast_hetero_guest', 'namespace': 'experiment'}}}}, 'host_10000': {'conf': {}, 'parties': [{'party_id': ['10000'], 'role': 'host'}], 'tasks': {'reader_1': {'conf': None, 'parameters': {'name': 'breast_hetero_host', 'namespace': 'experiment'}}}}}, 'stage': 'predict', 'tasks': {'hetero_secureboost_0': {'component_ref': 'hetero_secureboost', 'conf': None, 'dependent_tasks': ['psi_0'], 'inputs': {'data': {'test_data': {'task_output_artifact': [{'output_artifact_key': 'output_data', 'output_artifact_type_alias': None, 'parties': [{'party_id': ['9999'], 'role': 'guest'}, {'party_id': ['10000'], 'role': 'host'}], 'producer_task': 'psi_0'}]}}, 'model': {'input_model': {'model_warehouse': {'output_artifact_key': 'output_model', 'output_artifact_type_alias': None, 'parties': [{'party_id': ['9999'], 'role': 'guest'}, {'party_id': ['10000'], 'role': 'host'}], 'producer_task': 'hetero_secureboost_0'}}}}, 'outputs': None, 'parameters': {'gh_pack': True, 'goss': False, 'goss_start_iter': 0, 'hist_sub': True, 'l1': 0, 'l2': 0.1, 'learning_rate': 0.3, 'max_bin': 32, 'max_depth': 5, 'min_child_weight': 1, 'min_impurity_split': 0.01, 'min_leaf_node': 1, 'min_sample_split': 2, 'num_class': 2, 'num_trees': 1, 'objective': 'binary:bce', 'other_rate': 0.1, 'split_info_pack': True, 'top_rate': 0.2}, 'parties': None, 'stage': None}, 'psi_0': {'component_ref': 'psi', 'conf': None, 'dependent_tasks': ['reader_1'], 'inputs': {'data': {'input_data': {'task_output_artifact': {'output_artifact_key': 'output_data', 'output_artifact_type_alias': None, 'parties': [{'party_id': ['9999'], 'role': 'guest'}, {'party_id': ['10000'], 'role': 'host'}], 'producer_task': 'reader_1'}}}, 'model': None}, 'outputs': None, 'parameters': {}, 'parties': None, 'stage': 'default'}, 'reader_1': {'component_ref': 'reader', 'conf': None, 'dependent_tasks': None, 'inputs': None, 'outputs': None, 'parameters': {}, 'parties': None, 'stage': 'default'}}}, 'kind': 'fate', 'schema_version': '2.1.0'}, 'description': '', 'elapsed': 62958, 'end_time': 1712133480145, 'engine_name': 'eggroll', 'flow_id': '', 'inheritance': {}, 'initiator_party_id': '9999', 'job_id': '202404031636558952240', 'memory': 0, 'model_id': '202404031636558952240', 'model_version': '0', 'parties': [{'party_id': ['9999'], 'role': 'guest'}, {'party_id': ['10000'], 'role': 'host'}], 'party_id': '9999', 'progress': 100, 'protocol': 'fate', 'remaining_cores': 4, 'remaining_memory': 0, 'resource_in_use': False, 'return_resource_time': 1712133480016, 'role': 'guest', 'scheduler_party_id': '9999', 'start_time': 1712133417187, 'status': 'success', 'status_code': None, 'tag': 'job_end', 'update_time': 1712133480145, 'user_name': ''}
Total time: 0:01:04
```



#### 9. 删除部署



在部署机器上运行以下命令可以停止所有FATE集群：

```
bash ./docker_deploy.sh --delete all
```



如果想要彻底删除在运行机器上部署的FATE，可以分别登录节点，然后运行命令：

```
cd /data/projects/fate/confs-<id>/  # <id> 组织的id，本例中代表10000或者9999
docker-compose down
rm -rf ../confs-<id>/               # 删除docker-compose部署文件
```



#### 10. 可能遇到的问题



##### 采用docker hub下载镜像速度可能较慢

解决办法：可以自己构建镜像，自己构建镜像参考[这里](https://github.com/FederatedAI/FATE/tree/master/docker-build)。



##### 运行脚本`./docker_deploy.sh all`的时候提示需要输入密码

解决办法：检查免密登陆是否正常。ps:直接输入对应主机的用户密码也可以继续运行。



##### CPU指令集问题

解决办法：查看[wiki](https://github.com/FederatedAI/KubeFATE/wiki/KubeFATE)页面的storage-service部分