#!/bin/bash

source party.config

echo ${chartVersion}
echo ${imageTAG}
echo ${party_9999_IP}
echo ${party_10000_IP}
echo ${party_exchange_IP}

# 9999 config
sed -i "s/chartVersion: .*/chartVersion: ${chartVersion}/g" ./party-9999/cluster.yaml
sed -i "s/chartVersion: .*/chartVersion: ${chartVersion}/g" ./party-9999/cluster-serving.yaml
sed -i "s/chartVersion: .*/chartVersion: ${chartVersion}/g" ./party-9999/cluster-spark.yaml
sed -i "s/chartVersion: .*/chartVersion: ${chartVersion}/g" ./party-9999/cluster-spark-pulsar.yaml

sed -i "s/imageTag: .*/imageTag: ${imageTAG}/g" ./party-9999/cluster.yaml
sed -i "s/imageTag: .*/imageTag: ${imageTAG}/g" ./party-9999/cluster-serving.yaml
sed -i "s/imageTag: .*/imageTag: ${imageTAG}/g" ./party-9999/cluster-spark.yaml
sed -i "s/imageTag: .*/imageTag: ${imageTAG}/g" ./party-9999/cluster-spark-pulsar.yaml


sed -i "s/192.168.9.1/${party_9999_IP}/g" ./party-9999/cluster.yaml
sed -i "s/192.168.9.1/${party_9999_IP}/g" ./party-9999/cluster-serving.yaml
sed -i "s/192.168.9.1/${party_9999_IP}/g" ./party-9999/cluster-spark.yaml
sed -i "s/192.168.9.1/${party_9999_IP}/g" ./party-9999/cluster-spark-pulsar.yaml

sed -i "s/192.168.10.1/${party_10000_IP}/g" ./party-9999/cluster.yaml
sed -i "s/192.168.10.1/${party_10000_IP}/g" ./party-9999/cluster-serving.yaml
sed -i "s/192.168.10.1/${party_10000_IP}/g" ./party-9999/cluster-spark.yaml
sed -i "s/192.168.10.1/${party_10000_IP}/g" ./party-9999/cluster-spark-pulsar.yaml

sed -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-9999/cluster.yaml
sed -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-9999/cluster-serving.yaml
sed -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-9999/cluster-spark.yaml
sed -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-9999/cluster-spark-pulsar.yaml

# 10000 config

sed -i "s/chartVersion: .*/chartVersion: ${chartVersion}/g" ./party-10000/cluster.yaml
sed -i "s/chartVersion: .*/chartVersion: ${chartVersion}/g" ./party-10000/cluster-serving.yaml
sed -i "s/chartVersion: .*/chartVersion: ${chartVersion}/g" ./party-10000/cluster-spark.yaml
sed -i "s/chartVersion: .*/chartVersion: ${chartVersion}/g" ./party-10000/cluster-spark-pulsar.yaml

sed -i "s/imageTag: .*/imageTag: ${imageTAG}/g" ./party-10000/cluster.yaml
sed -i "s/imageTag: .*/imageTag: ${imageTAG}/g" ./party-10000/cluster-serving.yaml
sed -i "s/imageTag: .*/imageTag: ${imageTAG}/g" ./party-10000/cluster-spark.yaml
sed -i "s/imageTag: .*/imageTag: ${imageTAG}/g" ./party-10000/cluster-spark-pulsar.yaml

sed -i "s/192.168.9.1/${party_9999_IP}/g" ./party-10000/cluster.yaml
sed -i "s/192.168.9.1/${party_9999_IP}/g" ./party-10000/cluster-serving.yaml
sed -i "s/192.168.9.1/${party_9999_IP}/g" ./party-10000/cluster-spark.yaml
sed -i "s/192.168.9.1/${party_9999_IP}/g" ./party-10000/cluster-spark-pulsar.yaml

sed -i "s/192.168.10.1/${party_10000_IP}/g" ./party-10000/cluster.yaml
sed -i "s/192.168.10.1/${party_10000_IP}/g" ./party-10000/cluster-serving.yaml
sed -i "s/192.168.10.1/${party_10000_IP}/g" ./party-10000/cluster-spark.yaml
sed -i "s/192.168.10.1/${party_10000_IP}/g" ./party-10000/cluster-spark-pulsar.yaml

sed -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-10000/cluster.yaml
sed -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-10000/cluster-serving.yaml
sed -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-10000/cluster-spark.yaml
sed -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-10000/cluster-spark-pulsar.yaml


# exchange config

sed -i "s/chartVersion: .*/chartVersion: ${chartVersion}/g" ./party-exchange/rollsite.yaml
sed -i "s/chartVersion: .*/chartVersion: ${chartVersion}/g" ./party-exchange/trafficServer.yaml

sed -i "s/imageTag: .*/imageTag: ${imageTAG}/g" ./party-exchange/rollsite.yaml
sed -i "s/imageTag: .*/imageTag: ${imageTAG}/g" ./party-exchange/trafficServer.yaml

sed -i "s/192.168.9.1/${party_9999_IP}/g" ./party-exchange/rollsite.yaml
sed -i "s/192.168.9.1/${party_9999_IP}/g" ./party-exchange/trafficServer.yaml

sed -i "s/192.168.10.1/${party_10000_IP}/g" ./party-exchange/rollsite.yaml
sed -i "s/192.168.10.1/${party_10000_IP}/g" ./party-exchange/trafficServer.yaml

sed -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-exchange/rollsite.yaml
sed -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-exchange/trafficServer.yaml
