#!/bin/bash
DIR=$(cd $(dirname $0) && pwd)
source ${DIR}/../const.sh

source ${DIR}/common.sh
source ${DIR}/color.sh

# Get distribution
get_dist_name() {
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
  echo "dist_name: " $dist_name
}

centos() {
  yum remove docker docker-client docker-client-latest docker-common docker-latest docker-latest-logrotate docker-logrotate docker-engine
  yum install -y yum-utils
  yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
  yum erase podman buildah
  yum install -y docker-ce docker-ce-cli containerd.io
  systemctl start docker
}

fedora() {
  dnf remove docker docker-client docker-client-latest docker-common docker-latest docker-latest-logrotate docker-logrotate docker-selinux docker-engine-selinux docker-engine
  dnf -y install dnf-plugins-core
  dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo
  dnf install -y docker-ce docker-ce-cli containerd.io
  systemctl start docker
}

debian() {
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

ubuntu() {
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

install_separately() {
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
      ;;
    esac
  else
    echo "Fatal: Unknown system version"
    exit 1
  fi
}

trap 'onCtrlC' INT
function onCtrlC() {
  echo 'Ctrl+C is captured'
}

generate_cluster_config() {
  ip=$(kubectl get nodes -o wide | sed -n "2p" | awk -F ' ' '{printf $6}')
  cp ./cluster.yaml fate-9999.yaml
  if [[ $DOCKER_REGISTRY != "docker.io" ]]; then
    sed -i 's#registry: .*#registry: "'${DOCKER_REGISTRY}'"#g' fate-9999.yaml
  fi

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

EXPECT_PYTHON_STATUS=' * Running on http://'

check_fate_10000_fateflow_status() {
  MAX_TRY=20
  for ((i = 1; i <= MAX_TRY; i++)); do
    echo "# containers are ok"
    python_status=$(kubectl logs -n fate-10000 svc/fateflow -c python --tail 1 2>&1)
    echo "${python_status}"
    if [[ "${python_status}" =~ "${EXPECT_PYTHON_STATUS}" ]]; then
      return 0
    fi
    echo "fate-10000 fateflow successfully started"
    echo "# fate-9999 fateflow log: ${python_status} want ${EXPECT_PYTHON_STATUS}"
    sleep 3
  done

  return 1
}

check_fate_9999_fateflow_status() {
  MAX_TRY=20
  for ((i = 1; i <= MAX_TRY; i++)); do
    echo "# containers are ok"
    python_status=$(kubectl logs -n fate-9999 svc/fateflow -c python --tail 1 2>&1)
    echo "${python_status}"
    if [[ "${python_status}" =~ "${EXPECT_PYTHON_STATUS}" ]]; then
      echo "fate-9999 fateflow successfully started"
      return 0
    fi
    echo "# fate-9999 fateflow log: ${python_status} want ${EXPECT_PYTHON_STATUS}"
    sleep 3
  done

  return 1
}

main() {
  cd ${BASE_DIR}

  # install kubefate and kubefate CLI
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

  # Install two fate parties: fate-9999 and fate-10000
  kubectl create namespace fate-9999
  kubectl create namespace fate-10000

  # Copy the cluster.yaml sample in the working folder. One for party 9999, the other one for party 10000
  generate_cluster_config

  echo "[debug]"
  cat ./fate-9999.yaml
  cat ./fate-10000.yaml

  kubefate cluster install -f ./fate-9999.yaml
  kubefate cluster install -f ./fate-10000.yaml

  sleep 30s

  kubefate cluster ls

  selector_fate9999="name=fate-9999"
  selector_fate10000="name=fate-10000"
  kubectl wait --namespace fate-9999 \
    --for=condition=ready pod \
    --selector=${selector_fate9999} \
    --timeout=180s

  kubectl wait --namespace fate-10000 \
    --for=condition=ready pod \
    --selector=${selector_fate10000} \
    --timeout=180s

  check_fate_10000_fateflow_status

  check_fate_9999_fateflow_status

  kubectl exec -n fate-9999 -it svc/fateflow -c python -- bash -c "cd /data/projects/fate/examples/toy_example && \
  python run_toy_example.py 9999 9999 1"

  if [[ $? -eq 0 ]]; then
    loginfo "Unilateral test successful"
  else
    exit 1
  fi

  kubectl exec -n fate-9999 -it svc/fateflow -c python -- bash -c "cd /data/projects/fate/examples/toy_example && \
  python run_toy_example.py 9999 10000 1"
  if [[ $? -eq 0 ]]; then
    loginfo "Bilateral test successful"
  else
    exit 1
  fi
  # delete fate
  kubefate cluster ls | grep "fate" | awk '{print $1}' | xargs -n1 kubefate cluster delete

  sleep 30s

  kubectl delete namespace fate-9999
  kubectl delete namespace fate-10000

  kubefate_uninstall

  clean_host

  loginfo "deploy done."

}

main
