#!/bin/bash

helm package ./fate/ -d package/
helm package ./fate-serving/ -d package/
helm repo index .