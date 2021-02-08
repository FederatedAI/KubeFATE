# Copyright 2019-2020 VMware, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# you may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/bash

set -e

PREFIX=federatedai
IMG_TAG=latest

source .env

# nginx
docker build --build-arg SOURCE_PREFIX=${PREFIX} --build-arg SOURCE_TAG=${IMG_TAG} -o ${PREFIX}/nginx:${IMG_TAG} nginx

# python-spark
docker build --build-arg SOURCE_PREFIX=${PREFIX} --build-arg SOURCE_TAG=${IMG_TAG} -o ${PREFIX}/python-spark:${IMG_TAG} python-spark

# spark
docker build --build-arg SOURCE_PREFIX=${PREFIX} --build-arg SOURCE_TAG=${IMG_TAG} -o ${PREFIX}/spark-base:${IMG_TAG} spark/base

docker build --build-arg SOURCE_PREFIX=${PREFIX} --build-arg SOURCE_TAG=${IMG_TAG} -o ${PREFIX}/spark-master:${IMG_TAG} spark/master
docker build --build-arg SOURCE_PREFIX=${PREFIX} --build-arg SOURCE_TAG=${IMG_TAG} -o ${PREFIX}/spark-worker:${IMG_TAG} spark/worker
