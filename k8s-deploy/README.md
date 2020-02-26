# fate-cloud-agent



```bash
git clone https://gitlab.eng.vmware.com/fate/fate-cloud-agent.git

cd ./fate-cloud-agent
```

deploy 
```bash
$ kubectl apply -f ./rbac-config.yaml
$ kubectl apply -f ./kubefate.yaml

$ kubectl get all,ingress -n kube-fate
```

*Service pod must run successfully*

cluster deploy
```bash
$ kubefate cluster install -n <namespaces> -f ./cluster.yaml <clusterName>
```

cluster upgrade 
```bash
$ kubefate cluster upgrade -n <namespaces> -f ./cluster.yaml <clusterName>
```

cluster delete 
```bash
$ kubefate cluster delete <clusterId>
```

cluster list 
```bash
$ kubefate cluster list
```

cluster info 
```bash
$ kubefate cluster describe <clusterId>
```

job info
```bash
$ kubefate job describe <jobUUID>
```

job list
```bash
$ kubefate job list
```

job delete
```bash
$ kubefate job delete <jobUUID>
```
