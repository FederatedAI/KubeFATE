**Note**: The `master` branch undergoes frequent changes and can be *unstable* or even *broken* during development. To obtain a stable version, we recommend using the [releases](https://github.com/FederatedAI/KubeFATE/releases) instead of the `master` branch.


[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
# Overview
Federated learning involves multiple parties to collaboratively train a machine learning model, therefore it is usually based on a distributed system. KubeteFATE operationalizes federated learning workloads using cloud native technologies such as containers and Kubernetes. KubeFATE enables federated learning tasks to run across public, private and hybrid cloud environments.

[FATE](https://github.com/FederatedAI/FATE) (Federated AI Technology Enabler) is an open-source project to provide a secure computing framework to support the federated AI ecosystem. It implements secure computation protocols based on homomorphic encryption and secure multi-party computation (MPC). It supports federated learning architectures and secure computation of various machine learning algorithms, including logistic regression, tree-based algorithms, deep learning and transfer learning.

KubeFATE facilitates the deployment of FATE using both Docker Compose and Kubernetes. For a development environment of FATE, we recommend utilizing Docker Compose, while for a production environment, Kubernetes is the preferred option.

## Getting Involved
* To find answers to frequently asked questions, please refer to the [FAQ](https://github.com/FederatedAI/KubeFATE/wiki/KubeFATE#faqs) section.
* Please report bugs by submitting [issues](https://github.com/FederatedAI/KubeFATE/issues).
* Submit contributions using [pull requests](https://github.com/FederatedAI/KubeFATE/pulls).

## Project Structure
```
KubeFATE
|-- docker-deploy   
|-- k8s-deploy   
```
`docker-deploy`: The primary objective is to swiftly establish a federated learning environment. Docker Compose allows for the deployment of FATE components on a single host. Leveraging Docker Compose, FATE can be configured for multi-party environments, facilitating collaborative federated setups. For further information, kindly refer to the [Docker Compose Deployment](./docker-deploy/README.md) documentation.

`k8s-deploy`: The deployment approach is specifically tailored for production environments, providing a robust and scalable solution. Its design offers exceptional flexibility, enabling seamless operation of FATE clusters across various environments with ease and efficiency.

### Major features of KubeFATE k8s-deploy
  * Deliver an executable binary that simplifies the initialization and management of FATE clusters.
  * Provide the full life cycle management of FATE clusters, including deploying a new FATE cluster, querying an existing FATE cluster, destroying a given FATE cluster, etc.
  * Support customized deployment.
  * Efficiently manage multiple instances of FATE deployments simultaneously.
  * Offer a cluster management service with RESTful APIs.

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

## Specifying an image repository (registry)
By default, the installation script pulls the images from [Docker Hub](https://hub.docker.com/u/federatedai) during the deployment process.
* For docker compose mode, modify the `.env` file to specify the image registry.
* For K8s mode, check out this [offical document](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/).

## License
[Apache License 2.0](https://github.com/FederatedAI/FATE/blob/master/LICENSE)
