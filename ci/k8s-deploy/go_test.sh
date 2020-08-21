#!/bin/bash
dir=$(dirname $0)

cd ${dir}/../../k8s-deploy/pkg/api
go test

