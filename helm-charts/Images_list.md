# All images required to deploy the chart are recorded here.

## FATE

- federatedai/python:${version-tag}
- mysql:8
- federatedai/fateboard:${version-tag}
- federatedai/client:${version-tag}
- federatedai/eggroll:${version-tag}
- fluent/fluentd:v1.12
- federatedai/python-spark:${version-tag}
- federatedai/spark-master:${version-tag}
- federatedai/spark-worker:${version-tag}
- federatedai/hadoop-datanode:2.0.0-hadoop2.7.4-java8
- federatedai/hadoop-namenode:2.0.0-hadoop2.7.4-java8
- nginx:1.17
- federatedai/nginx:${version-tag}
- federatedai/rabbitmq:3.8.3-management
- federatedai/pulsar:2.7.0

## FATE-Serving

- federatedai/serving-server:${version-tag}
- federatedai/serving-proxy:${version-tag}
- redis:5
- federatedai/serving-admin:${version-tag}
- bitnami/zookeeper:3.7.0

## FATE-Exchange

- federatedai/eggroll:${version-tag}
- federatedai/trafficserver
- federatedai/nginx:${version-tag}
