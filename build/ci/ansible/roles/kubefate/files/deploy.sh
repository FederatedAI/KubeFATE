#! /bin/bash
DIR=$(dirname $0)
source ${DIR}/const.sh

# Get distribution
get_dist_name()
{
  if grep -Eqii "CentOS" /etc/issue || grep -Eq "CentOS" /etc/*-release; then
        dist_name='CentOS'
  elif grep -Eqi "Fedora" /etc/issue || grep -Eq "Fedora" /etc/*-release; then
        dist_name='Fedora'
  elif grep -Eqi "Debian" /etc/issue || grep -Eq "Debian" /etc/*-release; then
        dist_name='Debian'
  elif grep -Eqi "Ubuntu" /etc/issue || grep -Eq "Ubuntu" /etc/*-release; then
        dist_name='Ubuntu'
  else
        dist_name='Unknown'
  fi
  echo "dist_name: " $dist_name;
}

centos()
{
  yum remove docker docker-client docker-client-latest docker-common docker-latest docker-latest-logrotate docker-logrotate docker-engine
  yum install -y yum-utils
  yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
  yum erase podman buildah
  yum install -y docker-ce docker-ce-cli containerd.io
  systemctl start docker
}

fedora()
{
  dnf remove docker docker-client docker-client-latest docker-common docker-latest docker-latest-logrotate docker-logrotate docker-selinux docker-engine-selinux docker-engine
  dnf -y install dnf-plugins-core
  dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo
  dnf install -y docker-ce docker-ce-cli containerd.io
  systemctl start docker
}

debian()
{
  apt-get remove docker docker-engine docker.io containerd runc
  apt-get purge -y docker-ce docker-ce-cli containerd.io
  rm -rf /var/lib/docker
  apt-get update
  apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common
  curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
  add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"
  apt-get update
  apt-get install -y docker-ce docker-ce-cli containerd.io
  apt-get install -y docker-ce docker-ce-cli containerd.io
  systemctl start docker.service
}

ubuntu()
{
  apt-get remove docker docker-engine docker.io containerd runc
  apt-get purge -y docker-ce docker-ce-cli containerd.io
  rm -rf /var/lib/docker
  apt-get update
  apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common

  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

  add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
  apt-get update
  apt-get install -y docker-ce docker-ce-cli containerd.io
  systemctl start docker.service
}

install_separately()
{
  # Install Docker with different linux distibutions
  get_dist_name
  if [ $dist_name != "Unknown" ]; then
    case $dist_name in
      CentOS)
        centos
        ;;
      Fedora)
        fedora
        ;;
      Debian)
        debian
        ;;
      Ubuntu)
        ubuntu
        ;;
      *)
        echo "Unsupported distribution name"
    esac
  else
    echo "Fatal: Unknown system version"
    exit 1
  fi
}

trap 'onCtrlC' INT
function onCtrlC () {
  echo 'Ctrl+C is captured'
}

generate_cluster_config()
{
  ip=$(kubectl get nodes -o wide | sed -n "2p" | awk -F ' ' '{printf $6}')
  cp ./cluster.yaml fate-9999.yaml
  sed -i 's#registry: ""#registry: "'${DOCKER_REGISTRY}'/federatedai"#g' fate-9999.yaml

  sed -i 's$# rollsite:$rollsite:$g' fate-9999.yaml
  sed -i '0,/# type:/s//type:/' fate-9999.yaml
  sed -i '0,/# nodePort:/s//nodePort:/' fate-9999.yaml
  sed -i '0,/# partyList:/s//partyList:/' fate-9999.yaml
  sed -i '0,/# - partyId:/s//- partyId:/' fate-9999.yaml
  sed -i '0,/# partyIp: 192.168.0.1/s//partyIp: '${ip}'/' fate-9999.yaml
  sed -i '0,/# partyPort:/s//partyPort:/' fate-9999.yaml

  sed -i 's$# python:$python:$g' fate-9999.yaml
  sed -i '0,/# type:/s//type:/' fate-9999.yaml
  sed -i '0,/# httpNodePort:/s//httpNodePort:/' fate-9999.yaml
  sed -i '0,/# grpcNodePort:/s//grpcNodePort:/' fate-9999.yaml
  # delete rows commented
  sed -i '/#/'d fate-9999.yaml
  # delete continuous idle row
  sed -i 'N;/^\n/D' fate-9999.yaml

  cp ./fate-9999.yaml fate-10000.yaml
  sed -i 's/9999/10000/g' fate-10000.yaml
  sed -i '0,/nodePort: 30091/s//nodePort: 30101/' fate-10000.yaml
  sed -i '0,/- partyId: 10000/s//- partyId: 9999/' fate-10000.yaml
  sed -i '0,/partyPort: 30101/s//partyPort: 30091/' fate-10000.yaml
  sed -i '0,/httpNodePort: 30097/s//httpNodePort: 30107/' fate-10000.yaml
  sed -i '0,/grpcNodePort: 30092/s//grpcNodePort: 30102/' fate-10000.yaml
}

main()
{
  cd ${BASE_DIR}
  # Download KubeFATE Release Pack, KubeFATE Server Image v1.2.0 and Install KubeFATE Command Lines
  curl -LO https://github.com/FederatedAI/KubeFATE/releases/download/${KUBEFATE_CLI_VERSION}/kubefate-k8s-${KUBEFATE_CLI_VERSION}.tar.gz && tar -xzf ./kubefate-k8s-${KUBEFATE_CLI_VERSION}.tar.gz

  # Move the kubefate executable binary to path,
  chmod +x ./kubefate && mv ./kubefate /usr/bin

  # Download the KubeFATE Server Image
  curl -LO https://github.com/FederatedAI/KubeFATE/releases/download/${KUBEFATE_CLI_VERSION}/kubefate-${KUBEFATE_VERSION}.docker

  # Load into local Docker
  docker load < ./kubefate-${KUBEFATE_VERSION}.docker

  # Create kube-fate namespace and account for KubeFATE service
  kubectl apply -f ./rbac-config.yaml
  # kubectl apply -f ./kubefate.yaml

  # Replace the docker registry if it is not "docker.io"
  if [ "${DOCKER_REGISTRY}" != "docker.io" ]; then
    sed -i "s/mariadb:10/${DOCKER_REGISTRY}\/federatedai\/mariadb:10/g" kubefate.yaml
    sed -i "s/registry: \"\"/registry: \"${DOCKER_REGISTRY}\/federatedai\"/g" cluster.yaml
  fi
  kubectl apply -f ./kubefate.yaml

  # Check if the commands above have been executed correctly
  state=`kubefate version`
  if [ $? -ne 0 ]; then
    echo "Fatal: There is something wrong with the installation of kubefate, please check"
    exit 1
  fi

  # Install two fate parties: fate-9999 and fate-10000
  kubectl create namespace fate-9999
  kubectl create namespace fate-10000

  # Copy the cluster.yaml sample in the working folder. One for party 9999, the other one for party 10000
  generate_cluster_config

  # Start to install these two FATE cluster via KubeFATE with the following command
  echo "Waiting for kubefate service start to create container..."
  sleep ${KUBEFATE_SERVICE_TIMEOUT}

  selector_kubefate="fate=kubefate"
  kubectl wait --namespace kube-fate \
  --for=condition=ready pod \
  --selector=${selector_kubefate} \
  --timeout=${INGRESS_KUBEFATE_CLUSTER}s

  selector_mariadb="fate=mariadb"
  kubectl wait --namespace kube-fate \
  --for=condition=ready pod \
  --selector=${selector_mariadb} \
  --timeout=${INGRESS_KUBEFATE_CLUSTER}s

  echo "Waiting for kubefate service get ready..."
  i=0
  kubefate_status=`kubefate version`
  while [ $? -ne 0 ]
  do
    if [ $i == ${KUBEFATE_CLUSTER_TIMEOUT} ]; then
        echo "Can't install Ingress, Please check you environment"
        exit 1
    fi
    echo "Kubefate  Service Temporarily Unavailable, please wait..."
    sleep 1
    let i+=1
    kubefate_status=`kubefate version`
  done
  kubefate cluster install -f ./fate-9999.yaml
  kubefate cluster install -f ./fate-10000.yaml
  kubefate cluster ls

  sleep ${FATE_SERVICE_TIMEOUT}
  selector_fate9999="name=fate-9999"
  selector_fate10000="name=fate-10000"
  kubectl wait --namespace fate-9999 \
  --for=condition=ready pod \
  --selector=${selector_fate9999} \
  --timeout=${INGRESS_KUBEFATE_CLUSTER}s

  kubectl wait --namespace fate-10000 \
  --for=condition=ready pod \
  --selector=${selector_fate10000} \
  --timeout=${INGRESS_KUBEFATE_CLUSTER}s
}

main
