#!/bin/bash
source ~/.profile
DIR=$(dirname $0)
source ${DIR}/const.sh

main() {
  # federatedai/kubefate should build from source code
  docker pull ${DOCKER_REGISTRY}/federatedai/kubefate:${KUBEFATE_VERSION}
  kind load docker-image ${DOCKER_REGISTRY}/federatedai/kubefate:${KUBEFATE_VERSION}

  docker pull ${DOCKER_REGISTRY}/mariadb:10
  kind load docker-image ${DOCKER_REGISTRY}/mariadb:10

  docker pull ${DOCKER_REGISTRY}/mysql:8
  kind load docker-image ${DOCKER_REGISTRY}/mysql:8

  docker pull ${DOCKER_REGISTRY}/fluent/fluentd:v1.11-debian-1
  kind load docker-image ${DOCKER_REGISTRY}/fluent/fluentd:v1.11-debian-1

  docker pull ${DOCKER_REGISTRY}/federatedai/python-spark:${FATE_VERSION}
  kind load docker-image ${DOCKER_REGISTRY}/federatedai/python-spark:${FATE_VERSION}

  docker pull ${DOCKER_REGISTRY}/federatedai/hadoop-namenode:2.0.0-hadoop2.7.4-java8
  kind load docker-image ${DOCKER_REGISTRY}/federatedai/hadoop-namenode:2.0.0-hadoop2.7.4-java8

  docker pull ${DOCKER_REGISTRY}/federatedai/hadoop-datanode:2.0.0-hadoop2.7.4-java8
  kind load docker-image ${DOCKER_REGISTRY}/federatedai/hadoop-datanode:2.0.0-hadoop2.7.4-java8

  docker pull ${DOCKER_REGISTRY}/federatedai/nginx:1.17
  kind load docker-image ${DOCKER_REGISTRY}/federatedai/nginx:1.17

  docker pull ${DOCKER_REGISTRY}/federatedai/spark-master:${FATE_VERSION}
  kind load docker-image ${DOCKER_REGISTRY}/federatedai/spark-master:${FATE_VERSION}

  docker pull ${DOCKER_REGISTRY}/federatedai/spark-worker:${FATE_VERSION}
  kind load docker-image ${DOCKER_REGISTRY}/federatedai/spark-worker:${FATE_VERSION}

  docker pull ${DOCKER_REGISTRY}/federatedai/rabbitmq:3.8.3-management
  kind load docker-image ${DOCKER_REGISTRY}/federatedai/rabbitmq:3.8.3-management

  for image in "fateboard" "python" "eggroll" "client"; do
    docker pull ${DOCKER_REGISTRY}/federatedai/${image}:${FATE_VERSION}
    kind load docker-image ${DOCKER_REGISTRY}/federatedai/${image}:${FATE_VERSION}
  done
}

main
