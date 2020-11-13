#!/bin/bash



kubectl exec -n fate-9999 -it svc/fateflow -c python -- bash 
cd ../examples/toy_example/
python run_toy_example.py 9999 10000 1


# party Id 9999

kubectl -n fate-9999   exec -it   svc/fateboard -c python -- bash


cd ../examples/toy_example/
sed -i 's/    "backend": 0,/    "backend": 0,"spark_run": { "total-executor-cores": 12 },/g' toy_example_conf.json
sed -i 's/"partition": 48,/"partition": 4,/g' toy_example_conf.json
python run_toy_example.py 9999 10000 1 -b 1




# serving 

kubectl -n fate-10000   exec -it   svc/fateboard -c python -- bash

cd fate_flow;
sed -i "s/\"work_mode\": .*/\"work_mode\": 1,/g" examples/upload_host.json;
sed -i "s/\"backend\": .*/\"backend\": 1,/g" examples/upload_host.json;
python fate_flow_client.py -f upload -c examples/upload_host.json


kubectl  exec -it svc/fateflow -n fate-9999 -c python -- bash


cd fate_flow;
sed -i "s/\"work_mode\": .*/\"work_mode\": 1,/g" examples/upload_guest.json;
sed -i "s/\"backend\": .*/\"backend\": 1,/g" examples/upload_guest.json;
python fate_flow_client.py -f upload -c examples/upload_guest.json
