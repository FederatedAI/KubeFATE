#!/bin/bash

user=root
dir=/data/projects/fate
partylist=(9999 10000) 
partyiplist=(0.0.0.0 0.0.0.0)
venvdir=/data/projects/fate/venv

# party 1 will host the exchange by default
exchangeip=${partyiplist[0]}

