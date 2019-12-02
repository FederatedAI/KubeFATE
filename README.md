# Overview
Federated learning involves multiple parties to collaborately train a machine learning model, therefore it is usually based on a distributed system. KubeteFATE manages federated learning workloads using cloud native technologies such as containers. KubeFATE enables federated learning jobs to run across public, private and hybrid cloud environments.

Currently, KubeFATE supports the deployment of [FATE](https://github.com/FederatedAI/FATE) via Docker Compose and Kubernetes. 

## Container images of FATE components
All components of a FATE release are pre-built into Docker images. They can be pulled from Docker Hub directly. It is a preferred approach to install FATE. It saves much time in building FATE from the source code.

[Harbor](https://github.com/goharbor/harbor) can be used as a local registry to store and serve images of FATE. It can replicate container images from Docker Hub for a local environemnt. Harbor significantly improves performance and reduces network consumption, hence it is recommended for environments using containers.

To build images of FATE components from source code, refer to [Building FATE images](https://github.com/FederatedAI/FATE/tree/master/docker-build). To set up Harbor registry for your environment, refer to this [guide](./registry/README.md).

## Deployment with Docker Compose
Docker Compose can deploy FATE components on a single host. By using Docker Compose, FATE can be set up for environments of multiple parties which are collaborating in a federated manner. Please refer to [Docker Compose Deployment](./docker-deploy/README.md) for more details.

## Deployment on Kubernetes
To deploy FATE in the cloud or in a multi-node environment, a convenient way is to use a Kubernetes cluster as the underlying infrastructure. Helm Charts can be used to deploy FATE on Kubernetes. Please refer to [Kubernetes Deployment](./k8s-deploy/README.md) for more details.

## Note on the usage of ".env"
By default, the installation script pulls the images from [Docker Hub](https://hub.docker.com/u/federatedai) during the deployment. A user could also modify `.env` to specify a local registry (such as Harbor) to pull images from.

## License
[Apache License 2.0](https://github.com/FederatedAI/FATE/blob/master/LICENSE)
