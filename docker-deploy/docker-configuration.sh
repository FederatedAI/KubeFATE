#!/bin/bash

user=root
dir=/data/projects/fate
partylist=(9999 10000) 
partyiplist=(192.168.10.1 192.168.10.2)
venvdir=/data/projects/fate/venv

# party 1 will host the exchange by default
exchangeip=${partyiplist[0]}

