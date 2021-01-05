**Note**: The `master` branch may be in an *unstable* or *even broken* state during development. Please use [releases](https://github.com/FederatedAI/KubeFATE/releases) instead of the `master` branch in order to get a stable set of binaries.

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
# Overview
Federated learning involves multiple parties to collaborately train a machine learning model, therefore it is usually based on a distributed system. KubeteFATE manages federated learning workloads using cloud native technologies such as containers. KubeFATE enables federated learning jobs to run across public, private and hybrid cloud environments.

[FATE](https://github.com/FederatedAI/FATE) (Federated AI Technology Enabler) is an open-source project initiated by Webank's AI Department to provide a secure computing framework to support the federated AI ecosystem. It implements secure computation protocols based on homomorphic encryption and multi-party computation (MPC). It supports federated learning architectures and secure computation of various machine learning algorithms, including logistic regression, tree-based algorithms, deep learning and transfer learning.

KubeFATE supports the deployment of FATE via Docker Compose and Kubernetes. We recommend installing a quick development and playground FATE cluster with Docker Compose, while a production environment with Kubernetes. 

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
`docker-deploy`: The pricipal of `docker-deploy` is simple and quickly to set the environment up. Docker Compose can deploy FATE components on a single host. By using Docker Compose, FATE can be set up for environments of multiple parties which are collaborating in a federated manner. Please refer to [Docker Compose Deployment](./docker-deploy/README.md) for more details.

`k8s-deploy`: The k8s deployment is design for a real production deployed and managed environment. It designed for flexibility to suit different various environments.

### Major features of new KubeFATE k8s-deploy
  * Provide a single executable binary for initialing and managing FATE cluster
  * Full cycle FATE cluster management, includes deploying a new FATE cluster, querying existed FATE cluster, destroying a given FATE cluster and etc.
  * Support customized deployment
  * Support one KubeFATE to manage multiple FATE deployments
  * Provide cluster management service with RESTful APIs

For more details, please refer to [Kubernetes Deployment](./k8s-deploy/README.md).

## Build KubeFATE
##### To use docker-deploy for docker compose deployment, you need to make sure [Docker Compose] installed
Refer to: [Docker Compose Deployment](./docker-deploy/README.md) for more details

##### To build KubeFATE binary, you need a [Go environment] 

```
git clone https://github.com/FederatedAI/KubeFATE.git
cd KubeFATE
make build-linux-binary
```
##### To build KubeFATE service image, you need a [Docker environment]

```
git clone https://github.com/FederatedAI/KubeFATE.git
cd KubeFATE
make build-docker-image
```

## Note on the usage of ".env"
By default, the installation script pulls the images from [Docker Hub](https://hub.docker.com/u/federatedai) during the deployment. A user could also modify `.env` to specify a local registry (such as Harbor) to pull images from.

## License
[Apache License 2.0](https://github.com/FederatedAI/FATE/blob/master/LICENSE)
