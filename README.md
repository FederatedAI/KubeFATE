**Note**: The `master` branch is constantly in an *unstable* or *even broken* state during development. Please use [releases](https://github.com/FederatedAI/KubeFATE/releases) instead of the `master` branch in order to get a stable version.

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
# Overview
Federated learning involves multiple parties to collaboratively train a machine learning model, therefore it is usually based on a distributed system. KubeteFATE operationalizes federated learning workloads using cloud native technologies such as containers and Kubernetes. KubeFATE enables federated learning tasks to run across public, private and hybrid cloud environments.

[FATE](https://github.com/FederatedAI/FATE) (Federated AI Technology Enabler) is an open-source project to provide a secure computing framework to support the federated AI ecosystem. It implements secure computation protocols based on homomorphic encryption and secure multi-party computation (MPC). It supports federated learning architectures and secure computation of various machine learning algorithms, including logistic regression, tree-based algorithms, deep learning and transfer learning.

KubeFATE supports the deployment of FATE via Docker Compose and Kubernetes. We recommend installing a development environment of FATE via Docker Compose, and a production environment via Kubernetes. 

## Getting Involved
* For any frequently asked questions, you can check in [FAQ](https://github.com/FederatedAI/KubeFATE/wiki/KubeFATE#faqs).
* Please report bugs by submitting [issues](https://github.com/FederatedAI/KubeFATE/issues).
* Submit contributions using [pull requests](https://github.com/FederatedAI/KubeFATE/pulls)

## Project Structure
```
KubeFATE
|-- docker-deploy   
|-- k8s-deploy   
```
`docker-deploy`: The principle is to quickly set up an federated learning environment . Docker Compose can deploy FATE components on a single host. By using Docker Compose, FATE can be set up for environments of multiple parties which are collaborating in a federated manner. Please refer to [Docker Compose Deployment](./docker-deploy/README.md) for more details.

`k8s-deploy`: The deployment approach is designed for a production environment. It is designed with flexibility to operate FATE clusters in different environments. 

### Major features of KubeFATE k8s-deploy
  * Provide a single executable binary for initializing and managing FATE clusters.
  * Full life cycle management of FATE clusters, including deploying a new FATE cluster, querying an existing FATE cluster, destroying a given FATE cluster, etc.
  * Support customized deployment.
  * Support an instance of KubeFATE to manage multiple instances of FATE deployments.
  * Provide cluster management service with RESTful APIs.

For more details, please refer to [Kubernetes Deployment](./k8s-deploy/README.md).

## Building KubeFATE

To build the binary of KubeFATE (both CLI and KubeFATE service), a Golang development environment is needed.

```
$ git clone https://github.com/FederatedAI/KubeFATE.git
$ cd KubeFATE/k8s-deploy/
$ make kubefate-without-swag
```
To build the container image of KubeFATE service, a Docker environment is needed.

```
$ make docker-build
```

## Specifying an image repository
By default, the installation script pulls the images from [Docker Hub](https://hub.docker.com/u/federatedai) during the deployment. A user could also modify the file `.env` to specify a local registry (such as Harbor) to pull images from. A local registry can improve the efficiency of the deployment.

## License
[Apache License 2.0](https://github.com/FederatedAI/FATE/blob/master/LICENSE)
