#!/bin/bash
dir=$(dirname $0)
source $dir/init.sh

EXPECT_PYTHON_STATUS=' * Running on http://0.0.0.0:9380/ (Press CTRL+C to quit)'

echo "# config prepare"
target_dir=/data/projects/fate
target_party_id=10000
target_party_ip=$(hostname -I | awk '{print $1}')
sudo mkdir -p $target_dir
sudo rm -f $target_dir/confs-$target_party_id.tar $target_dir/serving-$target_party_id.tar
echo "# config is ok!"
echo "# training cluster and serving cluster copy begin"
cd ${dir}/../../docker-deploy
sudo cp ./outputs/confs-$target_party_id.tar ./outputs/serving-$target_party_id.tar $target_dir/
sudo rm -f ./outputs/confs-$target_party_id.tar ./outputs/serving-$target_party_id.tar
echo "# training cluster and serving cluster copy is ok!"
echo "# training cluster deploy begin"
cd $target_dir
sudo tar -xzf confs-$target_party_id.tar
cd confs-$target_party_id
sudo docker-compose down
sudo docker volume rm -f confs-${target_party_id}_shared_dir_examples
sudo docker volume rm -f confs-${target_party_id}_shared_dir_federatedml
sudo docker-compose up -d
cd ../
sudo rm -f confs-${target_party_id}.tar
echo "# party ${target_party_id} training cluster deploy is ok!"

echo "# serving cluster deploy begin"
sudo tar -xzf serving-$target_party_id.tar
cd serving-$target_party_id
sudo docker-compose down
sudo docker-compose up -d
cd $WD
sudo rm -f serving-$target_party_id.tar
echo "# party $target_party_id serving cluster deploy is ok!"
echo "# check containers"
MAX_TRY=10
for (( i=1; i<=$MAX_TRY; i++ ))
do
    result=$(sudo docker ps | wc -l)
    if [ $result -eq 11 ]
    then
            echo "# containers are ok"
	    python_status=$(docker logs confs-10000_python_1 --tail 1 2>&1)
	    echo "$python_status"
	    if [ "$python_status" = "$EXPECT_PYTHON_STATUS" ]
	    then
		    exit 0
	    fi
    fi
    echo "# Currently have containers: $result want 11"
    sleep 3
done
echo "# containers run overtime"
exit 1
