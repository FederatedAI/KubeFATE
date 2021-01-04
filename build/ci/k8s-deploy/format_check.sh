#!/bin/bash
dir=$(dirname $0)

cd ${dir}/../../k8s-deploy/pkg/
result=$(gofmt -l . | wc -l)
if [ $result -ne 0 ]
then
    gofmt -l .
    exit 1
fi
echo "# fmt is ok!"
exit 0
