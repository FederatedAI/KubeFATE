# Deploy an ATS exchange to pulsar

#### Create an Exchange Cluster Yaml

```bash
$ cat exchange.yaml 
name: fate-exchange
namespace: fate-exchange
chartName: fate-exchange
chartVersion: v1.6.0
partyId: 1
registry: ""
imageTag: "1.6.0-release"
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



#### Import the secret key to Kubernetes

Before that, you need to finish generating the secret key. https://github.com/FederatedAI/KubeFATE/blob/master/docs/FATE_On_Spark_With_Pulsar.md#generate-cas-certificate

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



#### Configure cluster to connect to exchange

```bash
$ cat cluster.yaml
...
pulsar:
  ...
  exchange:
    ip: 192.168.0.1
    port: 30000
```

Then you can deploy and test.