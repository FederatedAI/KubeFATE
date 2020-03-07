########################################################
# Copyright 2019-2020 program was created VMware, Inc. #
# SPDX-License-Identifier: Apache-2.0                  #
########################################################

#!/bin/bash

BASEDIR=$(dirname "$0")
cd $BASEDIR
WORKINGDIR=`pwd`

# fetch fate-python image
source ${WORKINGDIR}/.env
source ${WORKINGDIR}/parties.conf

cd ${WORKINGDIR}

Deploy() {
  if [ "$1" = "" ];then
    echo "No party id was provided, please check your arguments "
    exit 1
  fi

  while [ "$1" != "" ]; do
    case $1 in
      splitting_proxy)
        shift
        DeployPartyInternal $@
        break
        ;;
      all)
        for party in ${partylist[*]}
        do
          if [ "$2" != "" ]; then
            case $2 in
              --training)
              DeployPartyInternal $party
                if [ "${exchangeip}" != "" ]; then
                  DeployPartyInternal exchange
                fi
              ;;
              --serving)
              DeployPartyServing $party
              ;;
            esac
          else
          DeployPartyInternal $party
          DeployPartyServing $party
            if [ "${exchangeip}" != "" ]; then
              DeployPartyInternal exchange
            fi
          fi
        done
        break
        ;;
      *)
        if [ "$2" != "" ]; then
          case $2 in
            --training)
            DeployPartyInternal $1
            break
            ;;
            --serving)
            DeployPartyServing $1
            break
            ;;
          esac
        else
          DeployPartyInternal $1
          DeployPartyServing $1
        fi
        ;;
    esac
    shift

  done
}

Delete() {
  if [ "$1" = "" ];then
    echo "No party id was provided, please check your arguments "
    exit 1
  fi

  while [ "$1" != "" ]; do
    case $1 in
      all)
        for party in ${partylist[*]}
        do
          if [ "$2" != "" ]; then
            DeleteCluster $party $2
          else
            DeleteCluster $party
          fi
        done
        if [ "${exchangeip}" != "" ]; then
          DeleteCluster exchange
        fi
        break
        ;;
      *)
      DeleteCluster $@
      break
      ;;
    esac
  done
}

DeployPartyInternal() {
  target_party_id=$1
  # should not use localhost at any case
  target_party_ip="127.0.0.1"
  
  # check configuration files
  if [ ! -d ${WORKINGDIR}/outputs ];then
    echo "Unable to find outputs dir, please generate config files first."
    exit 1
  fi
  if [ ! -f ${WORKINGDIR}/outputs/confs-${target_party_id}.tar ];then
    echo "Unable to find deployment file for party $target_party_id, please generate it first."
    exit 1
  fi
  # extract the ip address of the target party
  if [ "$target_party_id" = "exchange" ];then
    target_party_ip=${exchangeip}
  elif [ "$2" != "" ]; then
    target_party_ip="$2"
  else
    for ((i=0;i<${#partylist[*]};i++))
    do
      if [ "${partylist[$i]}" = "$target_party_id" ];then
        target_party_ip=${partyiplist[$i]}
      fi
    done
  fi
  # verify the target_party_ip
  if [ "$target_party_ip" = "127.0.0.1" ]; then
    echo "Unable to find Party: $target_party_id, please check you input."
    exit 1
  fi

  if [ "$3" != "" ]; then
    user=$3
  fi

  scp ${WORKINGDIR}/outputs/confs-$target_party_id.tar $user@$target_party_ip:~/
  #rm -f ${WORKINGDIR}/outputs/confs-$target_party_id.tar
  echo "$target_party_ip training cluster copy is ok!"
  ssh -tt $user@$target_party_ip<< eeooff
mkdir -p $dir
mv ~/confs-$target_party_id.tar $dir
cd $dir
tar -xzf confs-$target_party_id.tar
cd confs-$target_party_id
docker-compose down
docker volume rm confs-${target_party_id}_shared_dir_examples
docker volume rm confs-${target_party_id}_shared_dir_federatedml
docker-compose up -d
cd ../
rm -f confs-${target_party_id}.tar
exit
eeooff
  echo "party ${target_party_id} deploy is ok!"
}

DeployPartyServing() {
  target_party_id=$1
  # should not use localhost at any case
  target_party_serving_ip="127.0.0.1"

  # check configuration files
  if [ ! -d ${WORKINGDIR}/outputs ];then
    echo "Unable to find outputs dir, please generate config files first."
    exit 1
  fi
  if [ ! -f ${WORKINGDIR}/outputs/serving-${target_party_id}.tar ];then
    echo "Unable to find deployment file for party $target_party_id, please generate it first."
    exit 1
  fi
  # extract the ip address of the target party
  for ((i=0;i<${#partylist[*]};i++))
  do
    if [ "${partylist[$i]}" = "$target_party_id" ];then
      target_party_serving_ip=${servingiplist[$i]}
    fi
  done
  # verify the target_party_ip
  if [ "$target_party_serving_ip" = "127.0.0.1" ]; then
    echo "Unable to find Party : $target_party_id serving address, please check you input."
    exit 1
  fi

  scp ${WORKINGDIR}/outputs/serving-$target_party_id.tar $user@$target_party_serving_ip:~/
  echo "party $target_party_id serving cluster copy is ok!"
  ssh -tt $user@$target_party_serving_ip<< eeooff
mkdir -p $dir
mv ~/serving-$target_party_id.tar $dir
cd $dir
tar -xzf serving-$target_party_id.tar
cd serving-$target_party_id
docker-compose down
docker-compose up -d
cd ../
rm -f serving-$target_party_id.tar
exit
eeooff
  echo "party $target_party_id serving cluster deploy is ok!"
}

DeleteCluster() {
  target_party_id=$1
  cluster_type=$2
  target_party_serving_ip="127.0.0.1"
  target_party_ip="127.0.0.1"

  # extract the ip address of the target party
  if [ "$target_party_id" == "exchange" ]; then
    target_party_ip=${exchangeip}
  else
    for ((i=0;i<${#partylist[*]};i++))
    do
      if [ "${partylist[$i]}" = "$target_party_id" ];then
        target_party_ip=${partyiplist[$i]}
      fi
    done
  fi

  for ((i=0;i<${#partylist[*]};i++))
  do
    if [ "${partylist[$i]}" = "$target_party_id" ];then
      target_party_serving_ip=${servingiplist[$i]}
    fi
  done

  # delete training cluster
  if [ "$cluster_type" == "--training" ]; then
    ssh -tt $user@$target_party_ip<< eeooff
cd $dir/confs-$target_party_id
docker-compose down
exit
eeooff
    echo "party $target_party_id training cluster is deleted!"
  # delete serving cluster
  elif [ "$cluster_type" == "--serving" ]; then
    ssh -tt $user@$target_party_serving_ip<< eeooff
cd $dir/serving-$target_party_id
docker-compose down
exit
eeooff
    echo "party $target_party_id serving cluster is deleted!"
  # delete training cluster and serving cluster
  else
    # if party is exchange then delete exchange cluster
    if [ "$target_party_id" == "exchange" ]; then
    ssh -tt $user@$target_party_ip<< eeooff
cd $dir/confs-$target_party_id
docker-compose down
exit
eeooff
    else
      ssh -tt $user@$target_party_ip<< eeooff
cd $dir/confs-$target_party_id
docker-compose down
exit
eeooff
      echo "party $target_party_id training cluster is deleted!"
      ssh -tt $user@$target_party_serving_ip<< eeooff
cd $dir/serving-$target_party_id
docker-compose down
exit
eeooff
      echo "party $target_party_id serving cluster is deleted!"
    fi
  fi
}

ShowUsage() {
  echo "Usage: "
  echo "Deploy all parties or specified partie(s): bash docker_deploy.sh partyid1[partyid2...] | all"
}

main() {
  if [ "$1" = "" ] || [ "$" = "--help" ]; then
    ShowUsage
    exit 1
  elif [ "$1" = "--delete" ] || [ "$1" = "--del" ]; then
    shift
    Delete $@
  else
    Deploy "$@"
  fi

  exit 0
}

main $@
