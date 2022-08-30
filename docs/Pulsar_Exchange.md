# Deploy an ATS exchange for pulsar

## Create an Exchange Cluster Yaml

```bash
$ cat exchange.yaml 
name: fate-exchange
namespace: fate-exchange
chartName: fate-exchange
chartVersion: v1.6.0
partyId: 1
registry: ""
pullPolicy: 
imagePullSecrets: 
- name: myregistrykey
persistence: false
istio:
  enabled: false
modules:
  - trafficServer

trafficServer:
  type: NodePort
  nodePort: 30000
  route_table: 
    sni:
    - fqdn: 10000.fate.org
      tunnelRoute: 192.168.10.1:30109
    - fqdn: 9999.fate.org
      tunnelRoute: 192.168.9.1:30099

```

## Import the secret keys to Kubernetes

Before that, you need to finish generating the secret key by folling this [doc](https://github.com/FederatedAI/KubeFATE/blob/master/docs/FATE_On_Spark_With_Pulsar.md#generate-cas-certificate).

Execute the corresponding command in the corresponding cluster.

```bash
kubectl -n fate-9999 create secret generic pulsar-cert \
	--from-file=broker.cert.pem=9999.fate.org/broker.cert.pem \
	--from-file=broker.key-pk8.pem=9999.fate.org/broker.key-pk8.pem \
	--from-file=ca.cert.pem=certs/ca.cert.pem
```

```bash
kubectl -n fate-10000  create secret generic pulsar-cert \
	--from-file=broker.cert.pem=10000.fate.org/broker.cert.pem \
	--from-file=broker.key-pk8.pem=10000.fate.org/broker.key-pk8.pem \
	--from-file=ca.cert.pem=certs/ca.cert.pem
```

```bash
kubectl -n fate-exchange  create secret generic traffic-server-cert \
	--from-file=proxy.cert.pem=proxy.fate.org/broker.cert.pem \
	--from-file=proxy.key.pem=proxy.fate.org/broker.key.pem \
	--from-file=ca.cert.pem=certs/ca.cert.pem
```

## Configure the clusters to connect to the exchange service

```bash
...
pulsar:
  ...
  exchange:
    ip: 192.168.0.1
    port: 30000
```

Then you can deploy and start to use the FATE cluster.