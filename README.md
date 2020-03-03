# Overview
This branch is for the new designed KubeFATE preview version.

Federated learning involves multiple parties to collaborately train a machine learning model, therefore it is usually based on a distributed system. KubeteFATE manages federated learning workloads using cloud native technologies such as containers. KubeFATE enables federated learning jobs to run across public, private and hybrid cloud environments.

[FATE](https://github.com/FederatedAI/FATE) (Federated AI Technology Enabler) is an open-source project initiated by Webank's AI Department to provide a secure computing framework to support the federated AI ecosystem. It implements secure computation protocols based on homomorphic encryption and multi-party computation (MPC). It supports federated learning architectures and secure computation of various machine learning algorithms, including logistic regression, tree-based algorithms, deep learning and transfer learning.

KubeFATE supports the deployment of FATE via Docker Compose and Kubernetes. We recommend installing a quick development and playground FATE cluster with Docker Compose, while a production environment with Kubernetes. 

## Project structure
```
KubeFATE2
|-- docker-deploy   
|-- k8s-deploy   
```
`docker-deploy`: The traditional  Docker Compose installation. The pricipal of `docker-deploy` is simple and quickly to set the environment up. Docker Compose can deploy FATE components on a single host. By using Docker Compose, FATE can be set up for environments of multiple parties which are collaborating in a federated manner. Please refer to [Docker Compose Deployment](./docker-deploy/README.md) for more details.

`k8s-deploy`: The new KubeFATE preview version. The k8s deployment is design for a real production deployed and managed environment. It designed for flexibility to suit different variable environments.

### Major features of new KubeFATE k8s-deploy
  * Provide a single executable binary for initialing and managing FATE cluster
  * Full cycle FATE cluster management, includes deploying a new FATE cluster, querying existed FATE cluster, destroying a given FATE cluster and etc.
  * Support customized deployment
  * Support one KubeFATE to manage multiple FATE deployments
  * Provide cluster management service with RESTful APIs

For more details, please refer to [Kubernetes Deployment](./k8s-deploy/README.md).

## Note on the usage of ".env"
By default, the installation script pulls the images from [Docker Hub](https://hub.docker.com/u/federatedai) during the deployment. A user could also modify `.env` to specify a local registry (such as Harbor) to pull images from.

## License
[Apache License 2.0](https://github.com/FederatedAI/FATE/blob/master/LICENSE)
