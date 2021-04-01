```bash
kubectl create secret generic pulsar-cert-9999 \
	--from-file=broker.cert.pem=9999.fate.org/broker.cert.pem \
	--from-file=broker.key-pk8.pem=9999.fate.org/broker.key-pk8.pem \
	--from-file=ca.cert.pem=certs/ca.cert.pem
```

```bash
kubectl create secret generic pulsar-cert-10000 \
	--from-file=broker.cert.pem=10000.fate.org/broker.cert.pem \
	--from-file=broker.key-pk8.pem=10000.fate.org/broker.key-pk8.pem \
	--from-file=ca.cert.pem=certs/ca.cert.pem
```

```bash
kubectl create secret generic traffic-server-cert \
	--from-file=proxy.cert.pem=proxy.fate.org/broker.cert.pem \
	--from-file=proxy.key.pem=proxy.fate.org/broker.key.pem 
```

