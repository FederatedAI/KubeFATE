#!/bin/bash

source ~/.profile

DIR=$(cd $(dirname $0) && pwd)
source ${DIR}/../const.sh

create_cluster_with_kind() {
    cat <<EOF | kind create cluster --config=-
    kind: Cluster
    apiVersion: kind.x-k8s.io/v1alpha4
    networking:
      apiServerAddress: "127.0.0.1"
      apiServerPort: 6443
    nodes:
    - role: control-plane
      kubeadmConfigPatches:
      - |
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "ingress-ready=true"
      extraPortMappings:
      - containerPort: 80
        hostPort: 80
        protocol: TCP
      - containerPort: 443
        hostPort: 443
        protocol: TCP
EOF
}

main() {
    cd ${BASE_DIR}

    create_cluster_with_kind

    docker pull ${DOCKER_REGISTRY}/jettech/kube-webhook-certgen:v1.5.0
    kind load docker-image ${DOCKER_REGISTRY}/jettech/kube-webhook-certgen:v1.5.0

    docker pull k8s.gcr.io/ingress-nginx/controller:v0.43.0
    kind load docker-image k8s.gcr.io/ingress-nginx/controller:v0.43.0

    curl_status=$(curl --version)
    if [[ $? -ne 0 ]]; then
        echo "Fatal: Curl does not installed correctly"
        clean
        exit 1
    fi

    # Check if kubectl is installed successfully
    kubectl_status=$(kubectl version --client)
    if [[ $? -eq 0 ]]; then
        echo "Kubectl is installed on this host, no need to install"
    else
        # Install the latest version of kubectl
        curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" && chmod +x ./kubectl && mv ./kubectl /usr/bin/
        kubectl_status=$(kubectl version --client)
        if [[ $? -ne 0 ]]; then
            echo "Fatal: Kubectl does not installed correctly"
            clean
            exit 1
        fi
    fi

    # Check if docker is installed already
    docker_status=$(docker ps)
    if [[ $? -eq 0 ]]; then
        echo "Docker is installed on this host, no need to install"
    else
        # Install Docker with different linux distibutions
        install_separately

        # check if docker is installed correctly
        docker=$(docker ps)
        if [[ $? -ne 0 ]]; then
            echo "Fatal: Docker does not installed correctly"
            clean
            exit 1
        fi
    fi
    INGRESS_FILE=${DIR}/../ingress.yml
    # Enable Ingress step 2.
    # sed -i "s#- --publish-status-address=localhost#- --publish-status-address=${ip}#g" ${INGRESS_FILE}
    kubectl apply -f ${INGRESS_FILE}

    # Config Ingress
    i=0
    cluster_ip=$(kubectl get service -o wide -A | grep ingress-nginx-controller-admission | awk -F ' ' '{print $4}')
    while [[ "$cluster_ip" == "" ]]; do
        if [[ $i == ${INGRESS_TIMEOUT} ]]; then
            echo "Can't install Ingress, Please check you environment"
            exit 1
        fi

        echo "Kind Ingress is not ready, Waiting for Ingress to get ready..."
        cluster_ip=$(kubectl get service -o wide -A | grep ingress-nginx-controller-admission | awk -F ' ' '{print $4}')
        sleep 1
        let i+=1
    done
    echo "Got Ingress Cluster IP: " $cluster_ip
    echo "Waiting for ${INGRESS_TIMEOUT} seconds util Ingress webhook get ready..."
    sleep ${INGRESS_TIMEOUT}
    selector="app.kubernetes.io/component=controller"
    kubectl wait --namespace ingress-nginx \
        --for=condition=ready pod \
        --selector=${selector} \
        --timeout=${INGRESS_KUBEFATE_CLUSTER}s

    # Reinstall Ingress
    kubectl apply -f ${INGRESS_FILE}

    ip=$(kubectl get nodes -o wide | sed -n "2p" | awk -F ' ' '{printf $6}')
    kubefate_domain=$(cat /etc/hosts | grep "example.com")
    if [[ "$kubefate_domain" == "" ]]; then
        echo "${ip}    example.com" >>/etc/hosts
    else
        sed -i "/example.com/d" /etc/hosts
        echo "${ip}    example.com" >>/etc/hosts
    fi

    ingress_nginx_controller_admission=$(cat /etc/hosts | grep "ingress-nginx-controller-admission")
    if [[ "$ingress_nginx_controller_admission" == "" ]]; then
        echo "${cluster_ip}    ingress-nginx-controller-admission" >>/etc/hosts
    else
        sed -i "/ingress-nginx-controller-admission/d" /etc/hosts
        echo "${cluster_ip}    ingress-nginx-controller-admission" >>/etc/hosts
    fi
}

main $1
