#!/bin/bash
echo "# toy_example test"
docker exec confs-10000_python_1 bash -c 'sed -i s/\"partition\":\ 48/\"partition\":\ 4/g /data/projects/fate/examples/toy_example/toy_example_conf.json'
docker exec confs-10000_python_1 bash -c 'python /data/projects/fate/examples/toy_example/run_toy_example.py 10000 10000 1'
echo "# test is ok!"
