ARG SOURCE_PREFIX=federatedai
ARG SOURCE_TAG=latest
FROM ${SOURCE_PREFIX}/python:${SOURCE_TAG}

WORKDIR /data/projects/fate/

RUN cd fate; pip install -e python/fate_client; pip install -e python/fate_test; pip install torchvision

COPY test_suite ./test_suite
