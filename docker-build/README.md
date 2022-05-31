# KubeFATE docker build

This contains the builds of some images for KubeFATE to deploy FATE.

Images list:

- federatedai/client
- federatedai/nginx
- federatedai/spark-master
- federatedai/spark-worker
- federatedai/python-spark
- federatedai/fate-test

## Prerequisites

1. A Linux host
2. Docker: 18+

## Build

All images build.

```bash
IMG_TAG=latest bash docker-build.sh all
```

## push

Push images to DockerHub

```bash
IMG_TAG=latest bash docker-build.sh all
```
