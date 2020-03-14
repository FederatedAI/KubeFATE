#!/bin/bash

helm package ./fate/ -d package/
helm repo index .