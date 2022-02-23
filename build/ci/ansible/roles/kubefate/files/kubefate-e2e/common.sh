#!/bin/bash
dir=$(cd $(dirname $0) && pwd)

source ~/.profile
source $dir/color.sh

kubefateWorkDir=$dir/../../../../../../../k8s-deploy

check_kubectl() {
    # check kubectl
    loginfo "check kubectl"
    kubectl version
    if [[ $? -ne 0 ]]; then
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
    if [[ $KubeFATE_IMG != "" ]]; then
        IMG=${KubeFATE_IMG}:${KUBEFATE_VERSION}
    else
        IMG=federatedai/kubefate:latest
    fi
    logdebug "IMG=${IMG}"
    # set kubefate image:tag
    sed -i "s#image: federatedai/kubefate:.*#image: ${IMG}#g" kubefate.yaml
    # deploy kubefate
    kubectl apply -f kubefate.yaml

    # check kubefate deploy success
    kubectl wait --namespace kube-fate \
        --for=condition=ready pod \
        --selector="fate=mariadb" \
        --timeout=180s
    if [[ $? != 0 ]]; then
        kubectl get pod -n kube-fate
        echo "kubefate deploy timeOut, please check"
        return 1
    fi

    # check kubefate deploy success
    kubectl wait --namespace kube-fate \
        --for=condition=ready pod \
        --selector="fate=kubefate" \
        --timeout=180s
    if [[ $? != 0 ]]; then
        kubectl get pod -n kube-fate
        echo "kubefate deploy timeOut, please check"
        return 1
    fi

    sleep 10
    echo "# kubefate deploy ok"
    return 0
}

set_host() {
    cd $kubefateWorkDir
    # get ingress nodeip
    ingressPodName=$(kubectl -n ingress-nginx get pod -l app.kubernetes.io/component=controller -o jsonpath='{.items[0].metadata.name}')
    ingressNodeIp=$(kubectl -n ingress-nginx get pod/$ingressPodName -o jsonpath='{.status.hostIP}')
    # set host
    loginfo "set hosts"
    echo $ingressNodeIp example.com >>/etc/hosts
    loginfo "$ingressNodeIp example.com"
    # set SERVICEURL
    loginfo "check kubefate version"
    # get ingress 80 nodeport
    ingressNodePort=$(kubectl -n ingress-nginx get svc/ingress-nginx-controller -o jsonpath='{.spec.ports[0].nodePort}')
    export FATECLOUD_SERVICEURL=example.com:$ingressNodePort
    echo $FATECLOUD_SERVICEURL
}

check_kubefate_version() {
    cd $kubefateWorkDir

    MAX_TRY=20
    for ((i = 1; i <= $MAX_TRY; i++)); do
        kubefate version
        if [[ $? == 0 ]]; then
            logsuccess "Check kubefate version success"
            return 0
        fi
        echo "[INFO] Kubefate version not Success"
        sleep 5
    done
    logerror "Kubefate command line error, checkout ingress"
    return 1
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

    if [[ $DOCKER_REGISTRY != "docker.io" ]]; then
        # set kubefate image:tag
        sed -i "s#registry: .*#registry: ${DOCKER_REGISTRY}#g" cluster.yaml
        sed -i "s#registry: .*#registry: ${DOCKER_REGISTRY}#g" cluster-spark.yaml
        sed -i "s#registry: .*#registry: ${DOCKER_REGISTRY}#g" cluster-serving.yaml
    fi
    if [[ $FATE_VERSION != "" ]]; then
        sed -i "s#imageTag: .*#imageTag: ${FATE_VERSION}#g" cluster.yaml
        sed -i "s#imageTag: .*#imageTag: ${FATE_VERSION}#g" cluster-spark.yaml
    fi
    if [[ $FATE_VERSION != "" ]]; then
        sed -i "s#imageTag: .*#imageTag: ${FATE_SERVING_VERSION}#g" cluster-serving.yaml
    fi
    logdebug "REGISTRY=${DOCKER_REGISTRY}"
    logdebug "FATE_IMG_TAG=${FATE_VERSION}"
    logdebug "FATE_SERVING_IMG_TAG=${fate_serving_version}"
}

jobUUID=""

cluster_install() {
    # create cluster
    loginfo "Cluster Install"
    rust=$(kubefate cluster install -f cluster.yaml)
    jobUUID=$(echo $rust | sed "s/^.*=//g" | sed "s/\r//g")
    logdebug "jobUUID=$jobUUID"
    if [[ $jobUUID == "" ]]; then
        logerror "$rust"
        return 1
    fi
    MAX_TRY=120
    for ((i = 1; i <= $MAX_TRY; i++)); do
        jobstatus=$(kubefate job describe $jobUUID | grep -w Status | awk '{print $2}')
        if [[ $jobstatus == "Success" ]]; then
            logsuccess "ClusterInstall job success"
            return 0
        fi
        if [[ $jobstatus != "Pending" ]] && [[ $jobstatus != "Running" ]]; then
            logerror "ClusterInstall job status error, status: $jobstatus"
            kubefate job describe $jobUUID
            return 1
        fi
        echo "[INFO] Current kubefate ClusterInstall job status: $jobstatus want Success"
        sleep 5
    done

    logerror "ClusterInstall job timeOut, please check"
    kubefate job describe $jobUUID
    return 1
}

cluster_update() {
    # update cluster
    loginfo "Cluster Update"
    rust=$(kubefate cluster update -f cluster-spark.yaml)
    jobUUID=$(echo $rust | sed "s/^.*=//g" | sed "s/\r//g")
    logdebug "jobUUID=$jobUUID"
    if [[ $jobUUID == "" ]]; then
        logerror "$rust"
        return 1
    fi
    for ((i = 1; i <= $MAX_TRY; i++)); do
        jobstatus=$(kubefate job describe $jobUUID | grep -w Status | awk '{print $2}')
        if [[ $jobstatus == "Success" ]]; then
            logsuccess "ClusterUpdate job success"
            return 0
        fi
        if [[ $jobstatus != "Pending" ]] && [[ $jobstatus != "Running" ]]; then
            logerror "ClusterUpdate job status error, status: $jobstatus"
            kubefate job describe $jobUUID
            return 1
        fi
        echo "[INFO] Current kubefate ClusterUpdate job status: $jobstatus want Success"
        sleep 3
    done

    logerror "ClusterUpdate job timeOut, please check"
    kubefate job describe $jobUUID
    return 1
}
clusterUUID=""
check_cluster_status() {
    # cluster list
    # gotUUID=$(bin/kubefate cluster list |  grep -w  | awk '{print $2}' )
    loginfo "Cluster Describe"
    clusterUUID=$(kubefate job describe $jobUUID | grep -w ClusterId | awk '{print $2}')
    logdebug "clusterUUID=$clusterUUID"
    # cluster describe
    clusterStatus=$(kubefate cluster describe $clusterUUID | grep -w Status | awk '{print $2}')
    if [[ $clusterStatus == "Running" ]]; then
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
    rust=$(kubefate cluster delete $clusterUUID)
    jobUUID=$(echo $rust | sed "s/^.*=//g" | sed "s/\r//g")
    logdebug "jobUUID=$jobUUID"
    if [[ $jobUUID == "" ]]; then
        logerror "$rust"
        return 1
    fi
    for ((i = 1; i <= $MAX_TRY; i++)); do
        jobstatus=$(kubefate job describe $jobUUID | grep -w Status | awk '{print $2}')
        if [[ $jobstatus == "Success" ]]; then
            logsuccess "ClusterDelete job success"
            return 0
        fi
        if [[ $jobstatus != "Pending" ]] && [[ $jobstatus != "Running" ]]; then
            logerror "ClusterDelete job status error, status: $jobstatus"
            kubefate job describe $jobUUID
            return 1
        fi
        echo "[INFO] Current kubefate ClusterDelete job status: $jobstatus want Success"
        sleep 3
    done
    logerror "ClusterDelete job timeOut, please check"
    kubefate job describe $jobUUID
    return 1
}

binary_install() {
    cd $kubefateWorkDir
    chmod +x ./bin/kubefate
    export PATH=$PATH:$kubefateWorkDir/bin
}
