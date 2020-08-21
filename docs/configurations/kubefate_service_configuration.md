`kubefate.yaml` is used to deploy KubeFATE service on kubernetes. Consists of two parts: KubeFATE Service and MongoDB.

## kubeFATE Service
kubeFATE Service have three parts: deployment, service and ingress.
### Deployment
Following is the description of the environment variables of KubeFATE deployment.
```
env:
  - name: FATECLOUD_DB_TYPE
    value: "mysql"
  - name: FATECLOUD_DB_HOST
    value: "mysql"
  - name: FATECLOUD_DB_PORT
    value: "3306"
  - name: FATECLOUD_DB_NAME
    value: "kube_fate"
  - name: FATECLOUD_DB_USERNAME
    value: "kubefate"
  - name: FATECLOUD_DB_PASSWORD
    value: "kubeFATE!23"
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
| Name                     | Description                                                  |
| ------------------------ | ------------------------------------------------------------ |
| FATECLOUD_DB_TYPE        | Database types that support kubefate. At present, there is only "mysql". |
| FATECLOUD_DB_HOST        | DB host, for example "localhost".                            |
| FATECLOUD_DB_PORT        | DB port, for example "3306".                                 |
| FATECLOUD_DB_NAME        | Database name.                                               |
| FATECLOUD_DB_USERNAME    | Database user.                                               |
| FATECLOUD_DB_PASSWORD    | Database password.                                           |
| FATECLOUD_REPO_NAME      | Remote helm [chart repository](https://helm.sh/docs/topics/chart_repository/) name. |
| FATECLOUD_REPO_URL       | Remote helm chart repository url.                            |
| FATECLOUD_USER_USERNAME  | Username of KubeFATE service default user.                   |
| FATECLOUD_USER_PASSWORD  | Password of KubeFATE service default user.                   |
| FATECLOUD_SERVER_ADDRESS | KubeFATE service address.                                    |
| FATECLOUD_SERVER_PORT    | KubeFATE service port.                                       |
| FATECLOUD_LOG_LEVEL      | KubeFATE service log level.                                  |

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

## MySQL
### Deployment
```
env:
  - name: MYSQL_DATABASE
    value: "kube_fate"
  - name: MYSQL_ALLOW_EMPTY_PASSWORD
    value: "1"
  - name: MYSQL_USER
    value: "kubefate"
  - name: MYSQL_PASSWORD
    value: "kubeFATE!23"
```
Define MySQL initial username and password.