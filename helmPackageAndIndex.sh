#!/bin/bash

helm lint ./fate/
helm package ./fate/ -d package/

helm lint ./fate-serving/
helm package ./fate-serving/ -d package/
helm repo index .