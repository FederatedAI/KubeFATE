#!/bin/bash
dir=$(dirname $0)
source $dir/init.sh

echo "# config prepare"
cd ${dir}/../../docker-deploy
host_ip=$(echo $(hostname -I | awk '{print $1}'))
cp ./parties.conf ./parties.conf.back
sed -i "/mysql_user=fate/ ! s/user=fate/user=root/g" ./parties.conf
sed -i "s/partylist=(10000 9999)/partylist=(10000)/g" ./parties.conf
sed -i "s/partyiplist=(192.168.1.1 192.168.1.2)/partyiplist=(${host_ip})/g" ./parties.conf
sed -i "s/servingiplist=(192.168.1.1 192.168.1.2)/servingiplist=(${host_ip})/g" ./parties.conf
cat ./parties.conf
echo "# config prepare is ok"

echo "# generate config"
bash generate_config.sh
echo "# check ./docker-deploy/outputs"
ls ./outputs
cd $WD
echo "# ok"
