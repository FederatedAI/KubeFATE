#!/bin/bash

if [ -z "${FATE_VERSION}" ]; then
    FATE_VERSION="latest"
fi

if [ -z "${FATE_SERVING_VERSION}" ]; then
    FATE_SERVING_VERSION="latest"
fi

if [ -z "${KUBEFATE_VERSION}" ]; then
    KUBEFATE_VERSION="latest"
fi

if [ -z "${DOCKER_REGISTRY}" ]; then
    DOCKER_REGISTRY="docker.io"
fi

if [ -z "${BASE_DIR}" ]; then
    BASE_DIR="/tmp"
fi

docker_version="docker-19.03.10"
dist_name=""
DEPLOY_DIR="${BASE_DIR}/cicd-${ANSIBLE_HOST}"
