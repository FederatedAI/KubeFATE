#!/bin/bash
dir=$(dirname $0)
cd ${dir}/../../docker-deploy

mv ./parties.conf.back ./parties.conf
docker rm -f $(docker ps -aq)
