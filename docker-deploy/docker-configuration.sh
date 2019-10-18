#!/bin/bash

user=luke
dir=/data/projects/fate
partylist=(3 4) 
partyiplist=(10.160.102.71 10.161.41.191)
venvdir=/data/projects/fate/venv

# party 1 will host the exchange by default
exchangeip=${partyiplist[0]}

