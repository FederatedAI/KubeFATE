# Eggroll with TLS

Since KubeFATE release v1.9.0, we can leverage KubeFATE to deploy eggroll-based FATE clusters who communicate with each other by TLS.

## Generate certificates
Preparations:

In a place where you can access your K8s cluster, run
```bash
mkdir my-ca
cd my-ca
wget https://raw.githubusercontent.com/apache/pulsar/master/tests/certificate-authority/openssl.cnf
export CA_HOME=$(pwd)
mkdir certs crl newcerts private
chmod 700 private/
touch index.txt
echo 1000 > serial
```
Generate the private key of the root cert
```bash
openssl genrsa -aes256 -out private/ca.key.pem 4096
chmod 400 private/ca.key.pem
```
private/ca.key.pem is a key you should not share with anyone.

Then generate the root certification using the private key:
```bash
openssl req -config openssl.cnf -key private/ca.key.pem \
    -new -x509 -days 7300 -sha256 -extensions v3_ca \
    -out certs/ca.cert.pem
chmod 444 certs/ca.cert.pem
```
When it prompts the requirement of the common name, you need to come up with one. In this example we use:
```
example.com
```
This root certificate is also the CA certificate, once it has been created, you can create certificate requests and sign them with this CA, with above common name as the suffix.

For one rollsite, we need 2 pairs of certifications and private keys, one pair for acting as a client and another one pair for acting as a server.

Suppose a site is called party-9999.

For client:

Generate the private key for the rollsite party, because rollsite doesn't support RSA formatted private key, we need to change to format to pkcs8.
```bash
mkdir fate-9999
openssl genrsa -out fate-9999/client_rsa.key 2048
openssl pkcs8 -topk8 -inform PEM -in fate-9999/client_rsa.key -outform PEM -out fate-9999/client.key -nocrypt
```
Generate the certificate request with the private key:
```
openssl req -config openssl.cnf -key fate-9999/client.key -new -sha256 -out fate-9999/client.csr
```
When it prompts the request common name, you can type in something like ```party-9999-client.example.com```.

Generate the client certification for rollsite party-9999:
```
openssl ca -config openssl.cnf -days 10000 -notext -md sha256 -in fate-9999/client.csr -out fate-9999/client.crt
```

For server:

The steps are similar, example:
```bash
openssl genrsa -out fate-9999/server_rsa.key 2048
openssl pkcs8 -topk8 -inform PEM -in fate-9999/server_rsa.key -outform PEM -out fate-9999/server.key -nocrypt
openssl req -config openssl.cnf -key fate-9999/server.key -new -sha256 -out fate-9999/server.csr
```
Type ```party-9999-server.example.com``` as the common name.
```bash
openssl ca -config openssl.cnf -days 10000 -extensions server_cert -notext -md sha256 -in fate-9999/server.csr -out fate-9999/server.crt
```
The last step is to create a K8s secret based on the generated files:
```bash
kubectl -n fate-9999 create secret generic eggroll-certs \
  --from-file=ca.pem=certs/ca.cert.pem \
  --from-file=client.key=fate-9999/client.key \
  --from-file=client.crt=fate-9999/client.crt \
  --from-file=server.key=fate-9999/server.key \
  --from-file=server.crt=fate-9999/server.crt 
```

## Enable TLS for rollsite in cluster.yaml file

Enable this switch under the rollsite module.
```yaml
rollsite: 
  enableTLS: true
```

## P2P mode

In this mode, a FATE cluster's rollsite will communicate with the rollsite of another FATE cluster.

So we need to:
1. Generate the certs files for each party
2. Create the K8s secrets for each party, in the corresponding K8s namespace
3. Turn on "enableTLS" for each party's rollsite

## Exchange mode
In this mode, every FATE cluster's rollsite will talk to the rollsite of the FATE-Exchange.

So we need to:
1. Generate the certs files for each party, and the FATE exchange
2. Create the K8s secrets for each party, in the corresponding K8s namespace
3. Turn on "enableTLS" for each party's rollsite, and also the exchange's rollsite

An example of the steps for exchange:

For client cert:
```bash
mkdir fate-exchange

openssl genrsa -out fate-exchange/client_rsa.key 2048
openssl pkcs8 -topk8 -inform PEM -in fate-exchange/client_rsa.key -outform PEM -out fate-exchange/client.key -nocrypt
openssl req -config openssl.cnf -key fate-exchange/client.key -new -sha256 -out fate-exchange/client.csr
```
Type in for example, ```exchange-client.example.com``` as the common name.
```bash
openssl ca -config openssl.cnf -days 10000 -notext -md sha256 -in fate-exchange/client.csr -out fate-exchange/client.crt
```

For server cert:
```bash
openssl genrsa -out fate-exchange/server_rsa.key 2048
openssl pkcs8 -topk8 -inform PEM -in fate-exchange/server_rsa.key -outform PEM -out fate-exchange/server.key -nocrypt
openssl req -config openssl.cnf -key fate-exchange/server.key -new -sha256 -out fate-exchange/server.csr
```
Type in for example, ```exchange-server.example.com``` as the common name.
```bash
openssl ca -config openssl.cnf -days 10000 -extensions server_cert -notext -md sha256 -in fate-exchange/server.csr -out fate-exchange/server.crt
```

Create K8s secret:
```bash
kubectl -n fate-exchange create secret generic eggroll-certs \
  --from-file=ca.pem=certs/ca.cert.pem \
  --from-file=client.key=fate-exchange/client.key \
  --from-file=client.crt=fate-exchange/client.crt \
  --from-file=server.key=fate-exchange/server.key \
  --from-file=server.crt=fate-exchange/server.crt 
```

Then in the cluster.yaml file of FATE-Exchange, turn on the ```enableTLS``` switch under the rollsite module.

## Docker-Compose mode

In KubeFATE release v1.9.1, we will not provide a switch for enabling TLS for rollsite. This can be done in below manual steps:

1. Generate the certs, as above documents shows, for every FATE cluster and for the FATE Exchange if needed.
2. Run `docker ps` to get the container id of the rollsite.
3. Run `docker exec -it <rollsite-container-id> bash` to get into the rollsite container.
4. In dir `/data/projects/fate/eggroll/conf`, run `mkdir cert`.
5. Edit 5 files: `ca.pem client.crt client.key server.crt server.key`, input the contents of the certifications and the private keys you have generated. `ca.pem` is the CA's certification, we assume you use this certification to sign both the client certification and the server certification of rollsite.
6. In dir `/data/projects/fate/eggroll/conf`, edit `eggroll.properties`, tail below contents:
    ```
    eggroll.core.security.secure.cluster.enabled=true
    eggroll.core.security.secure.client.auth.enabled=true
    eggroll.core.security.ca.crt.path=conf/cert/ca.pem
    eggroll.core.security.crt.path=conf/cert/server.crt
    eggroll.core.security.key.path=conf/cert/server.key
    eggroll.core.security.client.ca.crt.path=conf/cert/ca.pem
    eggroll.core.security.client.crt.path=conf/cert/client.crt
    eggroll.core.security.client.key.path=conf/cert/client.key
    ```
7. Run `docker restart <rollsite-container-id>`

After above steps, your eggroll deployed on Docker Compose should start to communicate with each other by TLS.

One way to verify that is to check the log of the rollsite container by `docker logs <rollsite-container-id>`, if you see logs like:

```
[INFO ][2046][2022-08-03 06:41:48,543][main,pid:1,tid:1][c.w.e.c.t.GrpcServerUtils:107] - gRPC server at port=9380 starting in secure mode. server private key path: /data/projects/fate/eggroll/conf/cert/server.key, key crt path: /data/projects/fate/eggroll/conf/cert/server.crt, ca crt path: /data/projects/fate/eggroll/conf/cert/ca.pem
[INFO ][2051][2022-08-03 06:41:48,548][main,pid:1,tid:1][c.w.e.r.EggSiteBootstrap:107] - secure server started at 9380
```

Then it indicates that you have enabled TLS on rollsite successfully on Docker Compose.
