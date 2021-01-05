# Offline Installation Guide

## Goal
Since Federated Learning is used in banking and financial organizationss, which are usually offline or internal network environment. KubeFATE supports offline installation to fit these use cases. The keys to install FATE and FATE-Serving is:
1. Kubernetes/Docker-compose has installed. We will discuss how to install Kubernetes offline in this guide, but not the key points; 
2. Container images is ready locally (Docker-Compose cases) or in local images registry. We have provides all images used in FATE and FATE-Serving in cloud storage. We need to download them and load into local image registry;
3. KubeFATE has set to pull images from local image registry.

We focus on the installing FATE and FATE-Serving on Kubernetes in this guide, which is more common-use in production environment and more complex. We will also discuss how it work in Docker-compose cases.

## Install Kubernetes offline
When you design to use KubeFATE to install FATE or FATE-Serving on Kubernetes, it suppose you have verified it before and the Kubernetes has been installed in offline environment by IT teams. If not, we can refer to several way to install Kubernetes:
1. Offline install Kubernetes with Kubespray: https://kubespray.io/#/docs/offline-environment. [Kubespray](https://kubespray.io/) is one of the most popular Kubernetes deploy tool. The advantage of it are: 
	1. can deploy Kubernetes on different cloud providers or baremetal with unified interface;
	2. HA Kubernetes cluster supported;
	3. Composable support: different network and storag plugin;
	4. Based on Ansible, which is powerful and friendly to DevOps.
2. Offline install Kubernetes with Kubeadm: https://gist.github.com/onuryilmaz/89a29261652299d7cf768223fd61da02. [Kubeadm](https://kubernetes.io/docs/reference/setup-tools/kubeadm/) is a tranditional tool to deploy Kubernetes. 

## Build local images registry (with Harbor)