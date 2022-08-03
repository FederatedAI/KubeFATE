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

Generate the private key for the rollsite party:
```bash
mkdir fate-9999
openssl genrsa -out fate-9999/client_priv.key 2048
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
openssl genrsa -out fate-9999/server_priv.key 2048
openssl pkcs8 -topk8 -inform PEM -in fate-9999/server_priv.key -outform PEM -out fate-9999/server.key -nocrypt
openssl req -config openssl.cnf -key fate-9999/server_priv.key -new -sha256 -out fate-9999/server.csr
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
2. Create the secrets for each party
3. Turn on "enableTLS" for each party's rollsite

## Exchange mode
TODO