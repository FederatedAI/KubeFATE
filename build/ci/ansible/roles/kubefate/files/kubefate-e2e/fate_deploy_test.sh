#!/bin/bash

dir=$(cd $(dirname $0) && pwd)
source $dir/color.sh

source $dir/common.sh

binary_install

if kubefate_install; then
    loginfo "kubefate install success"
else
    exit 1
fi

set_host

if check_kubefate_version; then
    loginfo "kubefate CLI ready"
else
    exit 1
fi

set_cluster_image

kubectl create namespace fate-9999

if cluster_install; then
    loginfo "cluster install success"
else
    exit 1
fi

if cluster_update; then
    loginfo "cluster update success"
else
    exit 1
fi

if check_cluster_status; then
    loginfo "check cluster status success"
else
    exit 1
fi

if cluster_delete; then
    loginfo "cluster delete success"
else
    exit 1
fi
loginfo "Cluster CURD test Success!"

kubectl delete namespace fate-9999

kubefate_uninstall

clean_host

loginfo "fate_deploy_test done."
exit 0
