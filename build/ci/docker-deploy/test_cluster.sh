#!/bin/bash
echo "# toy_example test"
docker exec confs-10000_client_1 bash -c 'flow test toy --guest-party-id 10000 --host-party-id 10000'
echo "# test is ok!"
