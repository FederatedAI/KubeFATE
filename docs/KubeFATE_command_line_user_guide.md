# KubeFATE Command Line User Guide

### What is KubeFATE command line

KubeFATE command line is a command line tool that connects KubeFATE  service to deploy FATE on Kubernetes.

### Install KubeFATE command line

Before using `kubefate` command line, you need to deploy KubeFATE  service.

#### Deploy KubeFATE service on Kubernetes

Get the code

```bash
git clone https://github.com/FederatedAI/KubeFATE.git
cd KubeFATE/k8s-deploy
```

Deploy KubeFATE service on Kubernetes

```
kubectl apply -f ./rbac-config.yaml
kubectl apply -f ./kubefate.yaml
```

*A more detailed deployment process is here([deploy KubeFATE in Kubernetes](https://github.com/FederatedAI/KubeFATE/tree/master/k8s-deploy#initial-a-new-fate-deployment)).*

#### Install command line

`kubefate` command line is developed by go and can run easily on Linux, Mac OS and Windows.

In addition to downloading and using the release package, it can also be compiled and installed on different platforms.

##### Linux

```bash
go build -o bin/kubefate kubefate.go
```

##### Mac OS

```bash
go build -o bin/kubefate kubefate.go
```

##### Windows

```bash
go build -o bin/kubefate.exe -buildmode=exe kubefate.go
```

Add `./bin` to the 'PATH' environment variable.

#### Modify configuration

Configuration in current working directory `config.yaml` file.

```bash
log:
  level: info
user:
  username: admin
  password: admin

serviceurl: example.com
```

#### Check

Use `kubefate version` to verify that the installation is successful.

### KubeFATE command line

If you have successfully installed `kubefate` command line, you can use these commands.

The `kubefate` command contains command actions and parameters.

#### cluster

Related operations of cluster command

##### install

Install a cluster

```bash
kubefate cluster install -f <cluster_config_yaml>
```

OPTIONS:

```bash
   --file value, -f value  Required, chart cluster.yaml
   --cover                 If the cluster already exists, overwrite the installation (default: false)
   --help, -h              show help (default: false)
```

If it runs successfully, a job_UUID will be returned. The cluster installation status can be obtained according to the job_UUID. 

*`<cluster_config_yaml>` means `cluster.yaml` `cluster- spart.yaml`  `cluster- serving.yaml` and so on.*

##### update

Update a cluster

```bash
kubefate cluster update -f <cluster_config_yaml>
```

OPTIONS:

```bash
   --file value, -f value  Required, chart cluster.yaml
   --help, -h              show help (default: false)
```

If it runs successfully, a job_UUID will be returned. According to the job_UUID, the cluster update status can be obtained. 

##### delete

Delete a cluster

```bash
kubefate cluster delete <cluster_uuid>
```
If it runs successfully, a job_UUID  will be returned. According to the job_UUID, the cluster deletion status can be obtained.

##### list

Get the list of currently running clusters.

```bash
kubefate cluster list
```

OPTIONS:
```bash
   --all, -A   List all clusters including deleted ones (default: false)
   --help, -h  show help (default: false)
```
##### describe

Get the description information of the given cluster.

```bash
kubefate cluster describe <cluster_uuid>
```

##### logs

Gets the component log for a given cluster. (If no component is specified, all logs will be obtained.)

```bash
kubefate cluster logs <cluster_uuid> [component]
```

OPTIONS:
```bash
   --follow, -f         Specify if the logs should be streamed. (default: false)
   --previous           If true, print the logs for the previous instance of the container in a pod if it exists. (defau
lt: false)
   --since value        Only return logs newer than a relative duration like 5s, 2m, or 3h. Defaults to all logs. Only o
ne of since-time since may be used. (default: 0s)
   --since-time value   Only return logs after a specific date (RFC3339). Defaults to all logs. Only one of since-time s
ince may be used. (default: (*time.Time)(nil))
   --timestamps         Include timestamps on each line in the log output. (default: false)
   --tail value         Lines of recent log file to display. Defaults to -1 with no selector, showing all log lines othe
rwise 10, if a selector is provided. (default: -1)
   --limit-bytes value  Maximum bytes of logs to return. Defaults to no limit. (default: 0)
   --help, -h           show help (default: false)
```
#### job

Through the install, update and delete of the cluster, the corresponding jobs will be generated.

##### list

Get the list of all jobs

```bash
kubefate job list
```

##### stop

Cancel the Running job

```bash
kubefate job stop <job_uuid>
```

*This only works for Running jobs of type ClusterInstall*

##### describe

Get the description information of the given job

```bash
kubefate job describe <job_uuid>
```

##### delete

Delete the record of the given job

```bash
kubefate job delete <job_uuid>
```

#### chart 

Chart is the management of the chart needed to install cluster.

##### upload

Upload chart file to KubeFATE service. The chart file must be generated by the `helm package`

```bash
kubefate chart upload -f <chart_file>
```

##### list

Get the chart list of existing KubeFATE services.

```bash
kubefate chart list
```

##### delete

Delete chart file from KubeFATE service.

```bash
kubefate chart delete <chart_uuid>
```

#### namespace

##### list

Get the namespace list of Kubernetes.

```bash
kubefate namespace list
```

#### user

##### list

Get the user list of KubeFATE service.

```bash
kubefate user list
```

##### describe

Obtain the specific user description information of KubeFATE service.

```bash
kubefate user describe <user_uuid>
```

#### version

View the corresponding version of KubeFATE service of the current command line and connection.

```bash
kubefate version
```

#### help

Get command line help.

```bash
kubefate help
```

All commands can be add OPTIONS of  `--help` to view help information.





If you have any questions, you can get help through issue.