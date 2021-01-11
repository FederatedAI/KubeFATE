#!/bin/bash
dir=$(cd $(dirname $0) && pwd)

source $dir/color.sh

kubefateWorkDir=$dir/../../k8s-deploy

check_kubectl() {
    # check kubectl
    loginfo "check kubectl"
    kubectl version
    if [ $? -ne 0 ]; then
        logerror "K8s environment abnormal"
        exit 1
    fi
}

kubefate_install() {
    # deploy kubefate
    cd $kubefateWorkDir

    loginfo "apply rbac"
    # namespace and rbac
    kubectl apply -f rbac-config.yaml

    loginfo "apply kubefate"
    # Is mirror specified
    if [$KubeFATE_IMG == ""]; then
        IMG=federatedai/kubefate:latest
    else
        IMG=$KubeFATE_IMG
    fi
    logdebug "IMG=${IMG}"
    # set kubefate image:tag
    sed -i "s#image: federatedai/kubefate:.*#image: ${IMG}#g" kubefate.yaml
    # deploy kubefate
    kubectl apply -f kubefate.yaml

    # check kubefate deploy success
    loginfo "check kubefate deploy success"
    MAX_TRY=60
    for ((i = 1; i <= $MAX_TRY; i++)); do
        status=$(kubectl get pod -l fate=kubefate -n kube-fate -o jsonpath='{.items[0].status.phase}')
        if [ $status == "Running" ]; then
            sleep 5
            echo "# kubefate are ok"
            return 0
        fi
        echo "# Current kubefate pod status: $status want Running"
        sleep 3
    done
    echo "kubefate deploy timeOut, please check"
    return 1
}

set_host() {
    cd $kubefateWorkDir
    # get ingress nodeip
    ingressPodName=$(kubectl -n ingress-nginx get pod -l app.kubernetes.io/component=controller -o jsonpath='{.items[0].metadata.name}')
    ingressNodeIp=$(kubectl -n ingress-nginx get pod/$ingressPodName -o jsonpath='{.status.hostIP}')
    # set host
    loginfo "set hosts"
    echo $ingressNodeIp kubefate.net >>/etc/hosts
    loginfo "$ingressNodeIp kubefate.net"
    # set SERVICEURL
    loginfo "check kubefate version"
    # get ingress 80 nodeport
    ingressNodePort=$(kubectl -n ingress-nginx get svc/ingress-nginx-controller -o jsonpath='{.spec.ports[0].nodePort}')
    export FATECLOUD_SERVICEURL=kubefate.net:$ingressNodePort
    echo $FATECLOUD_SERVICEURL
}

check_kubefate_version() {
    cd $kubefateWorkDir
    bin/kubefate version
    if [ $? -ne 0 ]; then
        logerror "kubefate command line error, checkout ingress"
        return 1
    fi
    return 0
}

kubefate_uninstall() {
    # delete
    loginfo "clean kubefate"
    kubectl delete -f kubefate.yaml
    kubectl delete -f rbac-config.yaml
}

clean_host() {
    # clean host
    sed -i '$d' /etc/hosts
}

build_chart() {
    cd $kubefateWorkDir
}

upload_chart() {
    cd $kubefateWorkDir
}

set_cluster_image() {
    # Is mirror specified
    if [$FATE_IMG_REGISTRY == ""]; then
        REGISTRY=""
    else
        REGISTRY=$FATE_IMG_REGISTRY
    fi
    if [$FATE_IMG_TAG == ""]; then
        FATE_IMG_TAG="latest"
    fi
    if [$FATE_SERVING_IMG_TAG == ""]; then
        FATE_SERVING_IMG_TAG="latest"
    fi
    # set kubefate image:tag
    sed -i "s#registry: ""#image: ${REGISTRY}#g" cluster.yaml
    sed -i "s#registry: ""#image: ${REGISTRY}#g" cluster-spark.yaml
    sed -i "s#registry: ""#image: ${REGISTRY}#g" cluster-serving.yaml

    sed -i "s#imageTag: ""#imageTag: ${FATE_IMG_TAG}#g" cluster.yaml
    sed -i "s#imageTag: ""#imageTag: ${FATE_IMG_TAG}#g" cluster-spark.yaml
    sed -i "s#imageTag: ""#imageTag: ${FATE_SERVING_IMG_TAG}#g" cluster-serving.yaml

    logdebug "REGISTRY=${IMG}"
    logdebug "FATE_IMG_TAG=${FATE_IMG_TAG}"
    logdebug "FATE_SERVING_IMG_TAG=${FATE_SERVING_IMG_TAG}"
}

jobUUID=""

cluster_install() {
    # create cluster
    loginfo "Cluster Install"
    rust=$(bin/kubefate cluster install -f cluster.yaml)
    jobUUID=$(echo $rust | sed "s/^.*=//g" | sed "s/\r//g")
    logdebug "jobUUID=$jobUUID"
    if [[ $jobUUID == "" ]]; then
        logerror "$rust"
        return 1
    fi
    MAX_TRY=120
    for ((i = 1; i <= $MAX_TRY; i++)); do
        jobstatus=$(bin/kubefate job describe $jobUUID | grep -w Status | awk '{print $2}')
        if [[ $jobstatus == "Success" ]]; then
            logsuccess "ClusterInstall job success"
            return 0
        fi
        if [[ $jobstatus != "Pending" ]] && [[ $jobstatus != "Running" ]]; then
            logerror "ClusterInstall job status error, status: $jobstatus"
            bin/kubefate job describe $jobUUID
            exit 1
        fi
        echo "[INFO] Current kubefate ClusterInstall job status: $jobstatus want Success"
        sleep 5
    done

    logerror "ClusterInstall job timeOut, please check"
    bin/kubefate job describe $jobUUID
    return 1
}

cluster_update() {
    # update cluster
    loginfo "Cluster Update"
    rust=$(bin/kubefate cluster update -f cluster-spark.yaml)
    jobUUID=$(echo $rust | sed "s/^.*=//g" | sed "s/\r//g")
    logdebug "jobUUID=$jobUUID"
    if [[ $jobUUID == "" ]]; then
        logerror "$rust"
        return 1
    fi
    for ((i = 1; i <= $MAX_TRY; i++)); do
        jobstatus=$(bin/kubefate job describe $jobUUID | grep -w Status | awk '{print $2}')
        if [[ $jobstatus == "Success" ]]; then
            logsuccess "ClusterUpdate job success"
            return 0
        fi
        if [[ $jobstatus != "Pending" ]] && [[ $jobstatus != "Running" ]]; then
            logerror "ClusterUpdate job status error, status: $jobstatus"
            bin/kubefate job describe $jobUUID
            return 1
        fi
        echo "[INFO] Current kubefate ClusterUpdate job status: $jobstatus want Success"
        sleep 3
    done

    logerror "ClusterUpdate job timeOut, please check"
    bin/kubefate job describe $jobUUID
    return 1
}
clusterUUID=""
check_cluster_status() {
    # cluster list
    # gotUUID=$(bin/kubefate cluster list |  grep -w  | awk '{print $2}' )
    loginfo "Cluster Describe"
    clusterUUID=$(bin/kubefate job describe $jobUUID | grep -w ClusterId | awk '{print $2}')
    logdebug "clusterUUID=$clusterUUID"
    # cluster describe
    clusterStatus=$(bin/kubefate cluster describe $clusterUUID | grep -w Status | awk '{print $2}')
    if [ $clusterStatus == "Running" ]; then
        logsuccess "Cluster Status is Running"
    else
        logerror "Cluster Status is $clusterStatus"
        return 1
    fi
    return 0
}

cluster_delete() {
    # delete cluster
    loginfo "Cluster Delete"
    rust=$(bin/kubefate cluster delete $clusterUUID)
    jobUUID=$(echo $rust | sed "s/^.*=//g" | sed "s/\r//g")
    logdebug "jobUUID=$jobUUID"
    if [[ $jobUUID == "" ]]; then
        logerror "$rust"
        return 1
    fi
    for ((i = 1; i <= $MAX_TRY; i++)); do
        jobstatus=$(bin/kubefate job describe $jobUUID | grep -w Status | awk '{print $2}')
        if [[ $jobstatus == "Success" ]]; then
            logsuccess "ClusterDelete job success"
            return 0
        fi
        if [[ $jobstatus != "Pending" ]] && [[ $jobstatus != "Running" ]]; then
            logerror "ClusterDelete job status error, status: $jobstatus"
            bin/kubefate job describe $jobUUID
            return 1
        fi
        echo "[INFO] Current kubefate ClusterDelete job status: $jobstatus want Success"
        sleep 3
    done
    logerror "ClusterDelete job timeOut, please check"
    bin/kubefate job describe $jobUUID
    return 1
}
