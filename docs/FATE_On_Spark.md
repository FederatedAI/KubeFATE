# Overview

Originally, the FATE use the underlying [EggRoll]("https://github.com/WeBankFinTech/eggroll") as the underlying computing engine, the
following picture illustrates the overall architecture.

<div align="center">
  <img src="./images/arch_eggroll.png">
</div>

As the above figure show, the EggRoll provide both computing and storage resource.

Since FATE v1.5.0 a user can select Spark as the underlying computing engine, however, spark itself is an in-memory computing engine without the data persistence. Thus, HDFS is also needed to be deployed to help on data persistence. For example, a user need to upload their data to HDFS through FATE before doing any training job, and the output data of each component will also be stored in the HDFS module.

**Currently the verifed Spark version is [3.1.2](https://archive.apache.org/dist/spark/spark-3.1.2/spark-3.1.2-bin-hadoop3.2.tgz) and the Hadoop is [3.2](https://archive.apache. org/dist/hadoop/common/hadoop-3.2/hadoop-3.2.tar.gz)**

The following picture shows the architecture of FATE on Spark:
<div align="center">
  <img src="./images/arch_spark_pulsar.png">
</div> 

In current implementation, the `fate_flow` service uses the `spark-submit` binary tool to submit job to the Spark cluster. With the configuration of the fate's job, a user can also specify the configuration for the spark application, here is an example:
```
{
  "initiator": {
    "role": "guest",
    "party_id": 10000
  },
  "job_parameters": {
    "spark_run": {
      "executor-memory": "4G",
      "total-executor-cores": 4
    },
    ...
```

The above configuration limit the maximum memory and the cores that can be used by the executors. For more about the supported "spark_run" parameters please refer to this [page](https://spark.apache.org/docs/latest/submitting-applications.html)