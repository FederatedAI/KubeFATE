# Overview
KubeteFATE provides tools to deploy FATE using Docker Compose and Kubernetes.

## Deployment with Docker Compose
A user can deploy FATE on a single host by using Docker Compose. Please refer to [Docker Deployment](./docker-deploy/README.md) for more details.

## Deployment on Kubernetes
For a multi-node deployment scenario, one of the solutions is to use a Kubernetes cluster as an underlying infrastructure to manage the FATE system. Please refer to [Kubernetes Deployment](./k8s-deploy/README.md) for more details.

## Usage of ".env"
By default, the script pulls the images from [Docker Hub](https://hub.docker.com/search?q=federatedai&type=image) during the deployment. A user could also modify `.env` to specify a registry to pull images from.

## License
[Apache License 2.0](https://github.com/FederatedAI/FATE/blob/master/LICENSE)
