#/bin/bash
source ~/.profile
DIR=$(dirname $0)
source ${DIR}/const.sh

clean() {
  # rm -rf ${BASE_DIR}/*

  # Delete kubernetes cluster
  kind_status=$(kind version)
  if [ $? -eq 0 ]; then
    echo "Deleting kubernetes cluster..."
    kind delete cluster
  fi

  # delete docker containers
  # docker stop $(docker ps -a -q)
  # docker rm $(docker ps -a -q)

  # delete docker images
  # docker rmi -f $(docker images --format "{{.Repository}}\t{{.Tag}}\t{{.ID}}" | grep -E '(${DOCKER_REGISTRY}/)?federatedai' | awk -F ' ' '{print $3}')
}

main() {
  if [ "$1" != "" ]; then
    if [ "$1" == "failed" ]; then
      clean
      echo "exit with errors"
      exit 1
    fi
  else
    clean
  fi
}

main $1
