#!/bin/bash
set -e
dir=$(dirname $0)

echo "# config prepare"
cd ${dir}/../../../docker-deploy
host_ip=$(echo $(hostname -I | awk '{print $1}'))

sed -i "/mysql_user=fate/ ! s/user=fate/user=root/g" ./parties.conf
sed -i "s/party_list=(10000 9999)/party_list=(10000)/g" ./parties.conf
sed -i "s/party_ip_list=(192.168.1.1 192.168.1.2)/party_ip_list=(${host_ip})/g" ./parties.conf
sed -i "s/serving_ip_list=(192.168.1.1 192.168.1.2)/serving_ip_list=(${host_ip})/g" ./parties.conf

# Replace tag to latest
# TODO should replace the serving as well
sed -i "s/^TAG=.*/TAG=latest/g" .env
echo "# config prepare is ok"

echo "# generate config"
bash generate_config.sh
