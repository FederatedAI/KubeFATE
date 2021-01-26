#!/bin/bash

source ~/.profile
dir=$(cd $(dirname $0) && pwd)

kubefateWorkDir=$dir/../../../../../../../k8s-deploy

docker_build() {
    cd $kubefateWorkDir
    make docker-build IMG=federatedai/kubefate:latest
}

binary_build() {
    cd $kubefateWorkDir
    make
}

binary_install() {
    cd $kubefateWorkDir
    chmod +x ./bin/kubefate
    export PATH=$PATH:$kubefateWorkDir/bin
}

binary_uninstall() {
    cd $kubefateWorkDir
    # exit
}

kubefate_image_load_to_kind_cluster() {
    kind load docker-image federatedai/kubefate:latest
    docker rmi federatedai/kubefate:latest
}

mariadb_image_load_to_kind_cluster() {
    docker pull ${DOCKER_REGISTRY}/mariadb:10
    kind load docker-image mariadb:10
}

mariadb_image_load_to_kind_cluster

docker_build
kubefate_image_load_to_kind_cluster
binary_build
binary_install
