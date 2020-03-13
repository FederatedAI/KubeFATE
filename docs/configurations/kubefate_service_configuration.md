`kubefate.yaml` is used to deploy KubeFATE service on kubernetes. Consists of two parts: KubeFATE Service and MongoDB.

## kubeFATE Service
kubeFATE Service have three parts: deployment, service and ingress.
### Deployment
Following is the description of the environment variables of KubeFATE deployment.
```
env:
  - name: FATECLOUD_MONGO_URL
    value: "mongo:27017"
  - name: FATECLOUD_MONGO_USERNAME
    value: "root"
  - name: FATECLOUD_MONGO_PASSWORD
    value: "root"
  - name: FATECLOUD_MONGO_DATABASE
    value: "KubeFate"
  - name: FATECLOUD_REPO_NAME
    value: "kubefate"
  - name: FATECLOUD_REPO_URL
    value: "https://federatedai.github.io/KubeFATE/"
  - name: FATECLOUD_USER_USERNAME
    value: "admin"
  - name: FATECLOUD_USER_PASSWORD
    value: "admin"
  - name: FATECLOUD_SERVER_ADDRESS
    value: "0.0.0.0"
  - name: FATECLOUD_SERVER_PORT
    value: "8080"
  - name: FATECLOUD_LOG_LEVEL
    value: "debug"
```
|Name                     |Description                                |
|-------------------------|-------------------------------------------|
|FATECLOUD_MONGO_URL      |MongoDB url, MongoDB configuration depends on MongoDB service. |
|FATECLOUD_MONGO_USERNAME |MongoDB username.                          |
|FATECLOUD_MONGO_PASSWORD |MongoDB password.                          |
|FATECLOUD_MONGO_DATABASE |Database in MongoDB.                       |
|FATECLOUD_REPO_NAME      |Remote helm [chart repository](https://helm.sh/docs/topics/chart_repository/) name. |
|FATECLOUD_REPO_URL       |Remote helm chart repository url.          |
|FATECLOUD_USER_USERNAME  |Username of KubeFATE service default user. |
|FATECLOUD_USER_PASSWORD  |Password of KubeFATE service default user. |
|FATECLOUD_SERVER_ADDRESS |KubeFATE service address.                  |
|FATECLOUD_SERVER_PORT    |KubeFATE service port.                     |
|FATECLOUD_LOG_LEVEL      |KubeFATE service log level.                |

### Service
```
ports:
  - name: "8080"
    port: 8080
    targetPort: 8080
    protocol: TCP
type: ClusterIP
```
Ports is defined in KubeFATE deployment. KubeFATE service uses ingress and clusterIP.

### Ingress
```
rules:
  - host: kubefate.net
    http:
      paths:
        - path: /
          backend:
            serviceName: kubefate
            servicePort: 8080
```
|Name    |Description                          |
|--------|-------------------------------------|
|host    |Defining the domain name of ingress. |
|backend |Backend depends on KubeFATE service. |

## MongoDB
### Deployment
```
env:
  - name: MONGO_INITDB_ROOT_USERNAME
    value: root
  - name: MONGO_INITDB_ROOT_PASSWORD
    value: root
```
Define MongoDB initial username and password.