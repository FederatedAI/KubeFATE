#!/bin/bash
echo "# toy_example test"
sudo docker exec confs-10000_python_1 bash -c "source /data/projects/python/venv/bin/activate&&/data/projects/python/venv/bin/python /data/projects/fate/python/examples/toy_example/run_toy_example.py 10000 10000 1"
echo "# test is ok!"
