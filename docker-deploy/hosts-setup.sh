#! /bin/bash
groupadd -f docker
useradd -s /bin/bash -g docker -d /home/fate -m fate
passwd fate
mkdir -p /data/projects/fate
chown -R fate:docker /data/projects/fate

