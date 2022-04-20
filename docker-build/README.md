# KubeFATE docker build

This contains the builds of some images for KubeFATE to deploy FATE.

- client
- nginx
- spark
- python-spark

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
