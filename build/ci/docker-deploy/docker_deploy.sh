#!/bin/bash
set -x
dir=$(dirname $0)

CONTAINER_NUM=13

echo "# config prepare"
target_dir=/data/projects/fate
target_party_id=10000
mkdir -p ${target_dir}
rm -f ${target_dir}/confs-${target_party_id}.tar ${target_dir}/serving-${target_party_id}.tar
echo "# config is ok!"
echo "# training cluster and serving cluster copy begin"

cd ${dir}/../../../docker-deploy
mv ./outputs/confs-${target_party_id}.tar ./outputs/serving-${target_party_id}.tar ${target_dir}/
echo "# training cluster and serving cluster copy is ok!"
echo "# training cluster deploy begin"

cd ${target_dir}
tar -xzf confs-${target_party_id}.tar

cd confs-${target_party_id}
docker-compose down
docker volume rm -f confs-${target_party_id}_shared_dir_examples
docker volume rm -f confs-${target_party_id}_shared_dir_federatedml
# exclude client service to save time !
docker-compose up -d

cd ../
rm -f confs-${target_party_id}.tar
echo "# party ${target_party_id} training cluster deploy is ok!"

echo "# serving cluster deploy begin"
tar -xzf serving-${target_party_id}.tar
cd serving-${target_party_id}
docker-compose down
docker-compose up -d

cd ../
rm -f serving-${target_party_id}.tar
echo "# party ${target_party_id} serving cluster deploy is ok!"
echo "# check containers"

MAX_TRY=12
for ((i = 1; i <= MAX_TRY; i++)); do
    result=$(docker ps | wc -l)
    if [ "${result}" -eq ${CONTAINER_NUM} ]; then
        echo "# containers are ok"
        FATE_FLOW_STATUS=$(curl -s -X POST localhost:9380/v1/version/get)
        success='"retmsg":"success"'
        result=$(echo $FATE_FLOW_STATUS | grep "${success}")
        if [[ "$result" != "" ]]
        then
          echo ${FATE_FLOW_STATUS}
          exit 0
        fi
        echo "FATEFLOW STATUS: ${FATE_FLOW_STATUS}, want ${success}"
    fi
    echo "# Currently have containers: ${result} want ${CONTAINER_NUM}"
    sleep 3
done
echo "# containers run overtime"
exit 1
