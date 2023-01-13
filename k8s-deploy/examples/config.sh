#!/bin/bash

source party.config

echo "FATE chartVersion: "${fate_chartVersion}
echo "FATE-Serving chartVersion: "${fate_serving_chartVersion}
echo "Party 9999 IP: "${party_9999_IP}
echo "Party 10000 IP: "${party_10000_IP}
echo "Party exchange IP: "${party_exchange_IP}

# Compatible with Mac
SED=sed
unamestr=`uname`
if [[ "$unamestr" == "Darwin" ]] ; then
    SED=gsed
    type $SED >/dev/null 2>&1 || {
        echo >&2 "$SED it's not installed. Try: brew install gnu-sed" ;
        exit 1;
    }
fi
# 9999 config
$SED -i "s/chartVersion: .*/chartVersion: ${fate_chartVersion}/g" party-9999/cluster.yaml
$SED -i "s/chartVersion: .*/chartVersion: ${fate_serving_chartVersion}/g" ./party-9999/cluster-serving.yaml
$SED -i "s/chartVersion: .*/chartVersion: ${fate_chartVersion}/g" ./party-9999/cluster-spark-rabbitmq.yaml
$SED -i "s/chartVersion: .*/chartVersion: ${fate_chartVersion}/g" ./party-9999/cluster-spark-pulsar.yaml
$SED -i "s/chartVersion: .*/chartVersion: ${fate_chartVersion}/g" ./party-9999/cluster-spark-local-pulsar.yaml

$SED -i "s/192.168.9.1/${party_9999_IP}/g" ./party-9999/cluster.yaml
$SED -i "s/192.168.9.1/${party_9999_IP}/g" ./party-9999/cluster-serving.yaml
$SED -i "s/192.168.9.1/${party_9999_IP}/g" ./party-9999/cluster-spark-rabbitmq.yaml
$SED -i "s/192.168.9.1/${party_9999_IP}/g" ./party-9999/cluster-spark-pulsar.yaml
$SED -i "s/192.168.9.1/${party_9999_IP}/g" ./party-9999/cluster-spark-local-pulsar.yaml

$SED -i "s/192.168.10.1/${party_10000_IP}/g" ./party-9999/cluster.yaml
$SED -i "s/192.168.10.1/${party_10000_IP}/g" ./party-9999/cluster-serving.yaml
$SED -i "s/192.168.10.1/${party_10000_IP}/g" ./party-9999/cluster-spark-rabbitmq.yaml
$SED -i "s/192.168.10.1/${party_10000_IP}/g" ./party-9999/cluster-spark-pulsar.yaml
$SED -i "s/192.168.10.1/${party_10000_IP}/g" ./party-9999/cluster-spark-local-pulsar.yaml

$SED -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-9999/cluster.yaml
$SED -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-9999/cluster-serving.yaml
$SED -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-9999/cluster-spark-rabbitmq.yaml
$SED -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-9999/cluster-spark-pulsar.yaml
$SED -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-9999/cluster-spark-local-pulsar.yaml


for item in `ls party-9999`
do
    $SED -i "s/algorithm: .*/algorithm: ${algorithm}/g" ./party-9999/$item
    $SED -i "s/device: .*/device: ${device}/g" ./party-9999/$item
done

# 10000 config

$SED -i "s/chartVersion: .*/chartVersion: ${fate_chartVersion}/g" ./party-10000/cluster.yaml
$SED -i "s/chartVersion: .*/chartVersion: ${fate_serving_chartVersion}/g" ./party-10000/cluster-serving.yaml
$SED -i "s/chartVersion: .*/chartVersion: ${fate_chartVersion}/g" ./party-10000/cluster-spark-rabbitmq.yaml
$SED -i "s/chartVersion: .*/chartVersion: ${fate_chartVersion}/g" ./party-10000/cluster-spark-pulsar.yaml
$SED -i "s/chartVersion: .*/chartVersion: ${fate_chartVersion}/g" ./party-10000/cluster-spark-local-pulsar.yaml

$SED -i "s/192.168.9.1/${party_9999_IP}/g" ./party-10000/cluster.yaml
$SED -i "s/192.168.9.1/${party_9999_IP}/g" ./party-10000/cluster-serving.yaml
$SED -i "s/192.168.9.1/${party_9999_IP}/g" ./party-10000/cluster-spark-rabbitmq.yaml
$SED -i "s/192.168.9.1/${party_9999_IP}/g" ./party-10000/cluster-spark-pulsar.yaml
$SED -i "s/192.168.9.1/${party_9999_IP}/g" ./party-10000/cluster-spark-local-pulsar.yaml

$SED -i "s/192.168.10.1/${party_10000_IP}/g" ./party-10000/cluster.yaml
$SED -i "s/192.168.10.1/${party_10000_IP}/g" ./party-10000/cluster-serving.yaml
$SED -i "s/192.168.10.1/${party_10000_IP}/g" ./party-10000/cluster-spark-rabbitmq.yaml
$SED -i "s/192.168.10.1/${party_10000_IP}/g" ./party-10000/cluster-spark-pulsar.yaml
$SED -i "s/192.168.10.1/${party_10000_IP}/g" ./party-10000/cluster-spark-local-pulsar.yaml

$SED -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-10000/cluster.yaml
$SED -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-10000/cluster-serving.yaml
$SED -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-10000/cluster-spark-rabbitmq.yaml
$SED -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-10000/cluster-spark-pulsar.yaml
$SED -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-10000/cluster-spark-local-pulsar.yaml

for item in `ls party-10000`
do
    $SED -i "s/algorithm: .*/algorithm: ${algorithm}/g" ./party-10000/$item
    $SED -i "s/device: .*/device: ${device}/g" ./party-10000/$item
done

# exchange config

$SED -i "s/chartVersion: .*/chartVersion: ${fate_chartVersion}/g" ./party-exchange/rollsite.yaml
$SED -i "s/chartVersion: .*/chartVersion: ${fate_chartVersion}/g" ./party-exchange/trafficServer.yaml

$SED -i "s/192.168.9.1/${party_9999_IP}/g" ./party-exchange/rollsite.yaml
$SED -i "s/192.168.9.1/${party_9999_IP}/g" ./party-exchange/trafficServer.yaml

$SED -i "s/192.168.10.1/${party_10000_IP}/g" ./party-exchange/rollsite.yaml
$SED -i "s/192.168.10.1/${party_10000_IP}/g" ./party-exchange/trafficServer.yaml

$SED -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-exchange/rollsite.yaml
$SED -i "s/192.168.0.1/${party_exchange_IP}/g" ./party-exchange/trafficServer.yaml
