########################################################
# Copyright 2019-2020 program was created VMware, Inc. #
# SPDX-License-Identifier: Apache-2.0                  #
########################################################

#!/bin/bash

BASEDIR=$(dirname "$0")
cd $BASEDIR
WORKINGDIR=`pwd`

# fetch fate-python image
source ${WORKINGDIR}/../.env
source ${WORKINGDIR}/kube.cfg

for ((i=0;i<${#partylist[*]};i++))
do
  eval partyid=\${partylist[${i}]}
  eval redispass=fate_dev
  eval jdbcrootpassword=fate_dev
  eval jdbcdbname=fate
  eval jdbcuser=fate
  eval jdbcpasswd=fate_dev
  
  rm -rf fate-$partyid/
  mkdir -p fate-$partyid/
  
  cp -r helm/fate-party/Chart.yaml fate-$partyid/Chart.yaml
  cp -r helm/fate-party/templates fate-$partyid/templates
  
  cat > fate-$partyid/values.yaml << EOF
#nfspath: /data/fate-data
#nfsserver: 192.168.0.2

image:
$( if [[ "$RegistryURI" != "" ]]
   then
     echo "  registry: ${RegistryURI}"
     echo "  isThridParty: true"
   else
     echo "  registry: ${PREFIX}"
     echo "  isThridParty:"
   fi )
  tag: ${TAG}
  pullPolicy: IfNotPresent
partyId: ${partyid}
host:
  fateboard: ${partyid}.fateboard.fedai.org
nodePort: ${partyiplist[${i}]#*:}

exchange:
  partyIp: ${exchangeip%:*}
  partyPort: ${exchangeip#*:}
mysql:
  mysql_root_password: ${jdbcrootpassword}
  mysql_database: ${jdbcdbname}
  mysql_user: ${jdbcuser}
  mysql_password: ${jdbcpasswd}
redis:
  requirepass: ${redispass}
eggList:
$( if [ ${#eggList[*]} == 0 ]
   then
     echo "  - egg: egg${j}"
     echo "    nodeLabel: ${nodeLabel}"
     echo "    value: ${eggList[${j}]}"
   else
     for ((j=0;j<${#eggList[*]};j++))
       do
         echo "  - egg: egg${j}"
         echo "    nodeLabel: ${nodeLabel}"
         echo "    value: ${eggList[${j}]}"
       done
   fi )
nodeSelector:
  fateboard:
    nodeLabel: ${nodeLabel}
    value: ${fateboard}
  federation:
    nodeLabel: ${nodeLabel}
    value: ${federation}
  metaService:
    nodeLabel: ${nodeLabel}
    value: ${metaService}
  mysql:
    nodeLabel: ${nodeLabel}
    value: ${mysql}
  proxy:
    nodeLabel: ${nodeLabel}
    value: ${proxy}
  python:
    nodeLabel: ${nodeLabel}
    value: ${python}
  redis:
    nodeLabel: ${nodeLabel}
    value: ${redis}
  roll:
    nodeLabel: ${nodeLabel}
    value: ${roll}
  servingServer:
    nodeLabel: ${nodeLabel}
    value: ${servingServer}
EOF

  echo fate-$partyid done!
done

# exchange

rm -rf fate-exchange/
mkdir -p fate-exchange/

cp -r helm/fate-exchange/Chart.yaml fate-exchange/Chart.yaml
cp -r helm/fate-exchange/templates fate-exchange/templates

cat > fate-exchange/values.yaml << EOF
image:
$( if [[ "$RegistryURI" != "" ]]
   then
     echo "  registry: ${RegistryURI}"
     echo "  isThridParty: true"
   else
     echo "  registry: ${PREFIX}"
     echo "  isThridParty:"
   fi )
  tag: ${TAG}
  pullPolicy: IfNotPresent
partyId: 0000

nodePort: ${exchangeip#*:}

partyList:
$( for ((j=0;j<${#partylist[*]};j++))
   do
     if [ ${i} -eq ${j} ]
     then
       continue
     fi
     echo "  - partyId: ${partylist[${j}]}"
     echo "    partyIp: ${partyiplist[${j}]%:*}"
     echo "    partyPort: ${partyiplist[${j}]#*:}"
   done )

nodeSelector:
  proxy:
    nodeLabel: fedai.org
    value: 
EOF
echo fate-exchange done!