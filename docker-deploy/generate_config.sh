########################################################
# Copyright 2019-2020 program was created VMware, Inc. #
# SPDX-License-Identifier: Apache-2.0                  #
########################################################

#!/bin/bash

BASEDIR=$(dirname "$0")
cd $BASEDIR
WORKINGDIR=`pwd`
deploy_dir=/data/projects/fate

# fetch fate-python image
source ${WORKINGDIR}/.env
source ${WORKINGDIR}/parties.conf

cd ${WORKINGDIR}

GenerateConfig() {
    for ((i=0;i<${#partylist[*]};i++))
    do
        eval party_id=\${partylist[${i}]}
        eval party_ip=\${partyiplist[${i}]}
        eval serving_ip=\${servingiplist[${i}]}
    
        eval processor_count=2
        eval venv_dir=/data/projects/python/venv
        eval python_path=${deploy_dir}/python:${deploy_dir}/eggroll/python
        eval data_dir=${deploy_dir}/data-dir
    
        eval egg_ip=(egg)
        eval egg_port=7888
    
        eval proxy_ip=proxy
        eval proxy_port=9370
    
        eval fateboard_ip=fateboard
        eval fateboard_port=8080
    
        eval fate_flow_ip=python
        eval fate_flow_grpc_port=9360
        eval fate_flow_http_port=9380
    
        eval db_ip=mysql
        eval db_user=fate
        eval db_password=fate_dev
        eval db_name=fate
    
        eval exchange_ip=${exchangeip}
    
        rm -rf confs-$party_id/
        mkdir -p confs-$party_id/confs
        cp -r docker-example-dir-tree/* confs-$party_id/confs/
        
        cp ./docker-compose.yml confs-$party_id/
        # generate conf dir
        cp ${WORKINGDIR}/.env ./confs-$party_id
        if [ "$RegistryURI" != "" ]
        then
          sed -i "s#PREFIX#${RegistryURI}#g" ./confs-$party_id/docker-compose.yml
          sed -i 's#image: "mysql:8"#image: "${RegistryURI}/mysql:8"#g' ./confs-$party_id/docker-compose.yml
          #sed -i 's#image: "redis:5"#image: "${RegistryURI}/redis:5"#g' ./confs-$party_id/docker-compose.yml
        fi

        # update the path of shared_dir
        ## handle examples
        shared_example_dir="confs-${party_id}/shared_dir/examples"
        mkdir -p "./$shared_example_dir"
        sed -i "s|/path/to/host/dir/examples|${dir}/${shared_example_dir}|g" ./confs-$party_id/docker-compose.yml
    
        ## handle federatedml
        shared_federatedml_dir="confs-${party_id}/shared_dir/federatedml"
        mkdir -p "./$shared_federatedml_dir"
        sed -i "s|/path/to/host/dir/federatedml|${dir}/${shared_federatedml_dir}|g" ./confs-$party_id/docker-compose.yml

        # egg config
        module_name=eggroll
	#db connect inf
        sed -i.bak "s#^eggroll.resourcemanager.clustermanager.jdbc.url=#eggroll.resourcemanager.clustermanager.jdbc.url=jdbc:mysql://${db_ip}:3306/${db_name}?characterEncoding=utf8\&characterSetResults=utf8\&autoReconnect=true\&failOverReadOnly=false\&serverTimezone=GMT%2B8#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
        sed -i.bak "s#^eggroll.resourcemanager.clustermanager.jdbc.username=fate"
        sed -i.bak "s#^eggroll.resourcemanager.clustermanager.jdbc.password=fate_dev"
        sed -i.bak "s#^eggroll.data.dir=data/"
        sed -i.bak "s#^eggroll.logs.dir=logs/"
        #clustermanager & nodemanager
        sed -i.bak "s#^eggroll.resourcemanager.clustermanager.host="
        sed -i.bak "s#^eggroll.resourcemanager.clustermanager.port=4670"
        sed -i.bak "s#^eggroll.resourcemanager.nodemanager.port=4671"
        sed -i.bak "s#^eggroll.resourcemanager.process.tag=fate-host"
        sed -i.bak "s#^eggroll.bootstrap.root.script=bin/eggroll_boot.sh"
        sed -i.bak "s#^eggroll.resourcemanager.bootstrap.egg_pair.exepath=bin/roll_pair/egg_pair_bootstrap.sh"
        #python env
        sed -i.bak "eggroll.resourcemanager.bootstrap.egg_pair.venv=${venv_dir}"
        #pythonpath, very import, do not modify."
        sed -i.bak "eggroll.resourcemanager.bootstrap.egg_pair.pythonpath=/data/projects/fate/python:/data/projects/fate/eggroll/python"
        sed -i.bak "eggroll.resourcemanager.bootstrap.egg_pair.filepath=python/eggroll/roll_pair/egg_pair.py"
        sed -i.bak "eggroll.resourcemanager.bootstrap.roll_pair_master.exepath=bin/roll_pair/roll_pair_master_bootstrap.sh"
        #javahome
        sed -i.bak "eggroll.resourcemanager.bootstrap.roll_pair_master.javahome=/data/projects/fate/common/jdk/jdk-8u192"
        sed -i.bak "eggroll.resourcemanager.bootstrap.roll_pair_master.classpath=conf/:lib/*"
        sed -i.bak "eggroll.resourcemanager.bootstrap.roll_pair_master.mainclass=com.webank.eggroll.rollpair.RollPairMasterBootstrap"
        sed -i.bak "eggroll.resourcemanager.bootstrap.roll_pair_master.jvm.options="
        # for roll site. rename in the next round
        sed -i.bak "eggroll.rollsite.coordinator=webank"
        sed -i.bak "eggroll.rollsite.host=192.168.0.1"
        sed -i.bak "eggroll.rollsite.port=9370"
        sed -i.bak "eggroll.rollsite.party.id=10000"
        sed -i.bak "eggroll.rollsite.route.table.path=conf/route_table.json"
        sed -i.bak "eggroll.session.processors.per.node=4"
        sed -i.bak "eggroll.session.start.timeout.ms=180000"
        sed -i.bak "eggroll.rollsite.adapter.sendbuf.size=1048576"
        sed -i.bak "eggroll.rollpair.transferpair.sendbuf.size=4150000"
        sed -i.bak "s#processor.venv=.*#processor.venv=${venv_dir}#g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s#processor.python-path=.*#processor.python-path=${python_path}#g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s#processor.engine-path=.*#processor.engine-path=${deploy_dir}/eggroll/python/eggroll/computing/processor.py#g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s#data-dir=.*#data-dir=${data_dir}#g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s#processor.logs-dir=.*#processor.logs-dir=${deploy_dir}/eggroll/python/eggroll/logs/processor#g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s#count=.*#count=${processor_count}#g" ./confs-$party_id/confs/egg/conf/egg.properties
        echo  >> ./confs-$party_id/confs/egg/conf/egg.properties
        echo "eggroll.computing.processor.python-path=${python_path}" >> ./confs-$party_id/confs/egg/conf/egg.properties
        echo egg module of $party_id done!
    
        # fateboard
        sed -i.bak "s#^server.port=.*#server.port=${fateboard_port}#g" ./confs-$party_id/confs/fateboard/conf/application.properties
        sed -i.bak "s#^fateflow.url=.*#fateflow.url=http://${fate_flow_ip}:${fate_flow_http_port}#g" ./confs-$party_id/confs/fateboard/conf/application.properties
        sed -i.bak "s/^spring.datasource.username=.*/spring.datasource.username=${db_user}/g" ./confs-$party_id/confs/fateboard/conf/application.properties
        sed -i.bak "s/^spring.datasource.password=.*/spring.datasource.password=${db_password}/g" ./confs-$party_id/confs/fateboard/conf/application.properties
        echo fateboard module of $party_id done!
    
        # mysql
        sed -i.bak "s/eggroll_meta/${db_name}/g" ./confs-$party_id/confs/mysql/init/create-meta-service.sql
        echo > ./confs-$party_id/confs/mysql/init/insert-node.sql
        echo "INSERT INTO node (ip, port, type, status) values ('${roll_ip}', '${roll_port}', 'ROLL', 'HEALTHY');" >> ./confs-$party_id/confs/mysql/init/insert-node.sql
        echo "INSERT INTO node (ip, port, type, status) values ('${proxy_ip}', '${proxy_port}', 'PROXY', 'HEALTHY');" >> ./confs-$party_id/confs/mysql/init/insert-node.sql
        for ((j=0;j<${#egg_ip[*]};j++))
        do
            echo "INSERT INTO node (ip, port, type, status) values ('${egg_ip[j]}', '${egg_port}', 'EGG', 'HEALTHY');" >> ./confs-$party_id/confs/mysql/init/insert-node.sql
        done
        for ((j=0;j<${#storage_service_ip[*]};j++))
        do
            echo "INSERT INTO node (ip, port, type, status) values ('${storage_service_ip[j]}', '${storage_service_port}', 'STORAGE', 'HEALTHY');" >> ./confs-$party_id/confs/mysql/init/insert-node.sql
        done
        echo "show tables;" >> ./confs-$party_id/confs/mysql/init/insert-node.sql
        echo "select * from node;" >> ./confs-$party_id/confs/mysql/init/insert-node.sql
        echo mysql module of $party_id done!

        # fate_flow
        sed -i.bak "s/user:.*/user: '${db_user}',/g" ./confs-$party_id/confs/fate_flow/conf/base_conf.yaml
        sed -i.bak "s/passwd:.*/passwd: '${db_password}',/g" ./confs-$party_id/confs/fate_flow/conf/base_conf.yaml
        sed -i.bak "s/host:.*/host: '${db_ip}',/g" ./confs-$party_id/confs/fate_flow/conf/base_conf.yaml
        sed -i.bak "s/name:.*/name: '${db_name}',/g" ./confs-$party_id/confs/fate_flow/conf/base_conf.yaml
        sed -i.bak "s/serving:8000/${serving_ip}:8000/g" ./confs-$party_id/confs/fate_flow/conf/server_conf.json
    
        cat > ./confs-$party_id/confs/fate_flow/server_conf.json <<EOF
{
    "servers": {
        "proxy": {
            "host": "${proxy_ip}",
            "port": ${proxy_port}
        },
        "fateboard": {
            "host": "${fateboard_ip}",
            "port": ${fateboard_port}
        },
        "roll": {
            "host": "${roll_ip}",
            "port": ${roll_port}
        },
        "fateflow": {
            "host": "${fate_flow_ip}",
            "grpc.port": ${fate_flow_grpc_port},
            "http.port": ${fate_flow_http_port}
        },
        "federation": {
            "host": "${federation_ip}",
            "port": ${federation_port}
        },
        "clustercomm": {
            "host": "${federation_ip}",
            "port": ${federation_port}
        }
    }
}
EOF
        echo fate_flow module of $party_id done!
    
        # rollsite 
        module_name="rollsite"
        sed -i.bak "s/port=.*/port=${proxy_port}/g" ./confs-$party_id/confs/proxy/conf/${module_name}.properties
        sed -i.bak "s#route.table=.*#route.table=${deploy_dir}/${module_name}/conf/route_table.json#g" ./confs-$party_id/confs/proxy/conf/${module_name}.properties
        sed -i.bak "s/coordinator=.*/coordinator=${party_id}/g" ./confs-$party_id/confs/proxy/conf/${module_name}.properties
    
    cat > ./confs-$party_id/confs/proxy/conf/route_table.json <<EOF
{
    "route_table": {
        "default": {
            "default": [
                {
$( if [ "$exchange_ip" != "" ]; then
echo "                    \"ip\": \"${exchange_ip}\",
                    \"port\": 9371 "
else
echo "                    \"ip\": \"proxy\",
                    \"port\": 9370 "
fi )
                }
            ]
        },
$( for ((j=0;j<${#partylist[*]};j++));do
if [ "${party_id}" == "${partylist[${j}]}" ]; then
continue
fi
echo "        \"${partylist[${j}]}\": {
            \"default\": [
                {
                    \"ip\": \"${partyiplist[${j}]}\",
                    \"port\": 9370
                }
            ]
        },"
        done)
        "${party_id}": {
            "fateflow": [
            {
                "ip": "${fate_flow_ip}",
                "port": ${fate_flow_grpc_port}
            }],
            "fate": [
            {
                "ip": "${federation_ip}",
                "port": ${federation_port}
            }
            ]
        }
    },
    "permission": {
        "default_allow": true
    }
}
EOF
    tar -czf ./outputs/confs-$party_id.tar ./confs-$party_id
    rm -rf ./confs-$party_id
    echo proxy module of $party_id done!
    done

  if [ "$exchange_ip" != "" ]; then
  
  
    # handle exchange
    echo "handle exchange"
    module_name=exchange
    cd ${WORKINGDIR}
    rm -rf confs-exchange/
    mkdir -p confs-exchange/conf
    cp ${WORKINGDIR}/.env confs-exchange/
    cp docker-compose-exchange.yml confs-exchange/docker-compose.yml
    cp -r docker-example-dir-tree/proxy/conf confs-exchange/
    
    if [ "$RegistryURI" != "" ]; then
        sed -i "s#PREFIX#RegistryURI#g" ./confs-exchange/docker-compose.yml
    fi
    sed -i.bak "s/port=.*/port=${proxy_port}/g" ./confs-exchange/conf/proxy.properties
    sed -i.bak "s#route.table=.*#route.table=${deploy_dir}/proxy/conf/route_table.json#g" ./confs-exchange/conf/proxy.properties
    sed -i.bak "s/coordinator=.*/coordinator=${party_id}/g" ./confs-exchange/conf/proxy.properties
    sed -i.bak "s/ip=.*/ip=0.0.0.0/g" ./confs-exchange/conf/proxy.properties
    
    cat > ./confs-exchange/conf/route_table.json <<EOF
{
    "route_table": {
$( for ((j=0;j<${#partylist[*]};j++));do
echo "        \"${partylist[${j}]}\": {
            \"default\": [
                {
                    \"ip\": \"${partyiplist[${j}]}\",
                    \"port\": 9370
                }
            ]
        },"
        done)
        "default": {
            "default": [
                {
                }
            ]
        }
    },
    "permission": {
        "default_allow": true
    }
}
EOF
    tar -czf ./outputs/confs-exchange.tar ./confs-exchange
    rm -rf ./confs-exchange
    echo exchange module done!
    
  fi
  
    # handle serving
    echo "handle serving"
    for ((i=0;i<${#servingiplist[*]};i++))
    do
        eval party_id=\${partylist[${i}]}
        eval party_ip=\${partyiplist[${i}]}
        eval serving_ip=\${servingiplist[${i}]}

        rm -rf serving-$party_id/
        mkdir -p serving-$party_id/confs
        cp -r docker-serving/* serving-$party_id/confs/
        
        cp ./docker-compose-serving.yml serving-$party_id/docker-compose.yml
        # generate conf dir
        cp ${WORKINGDIR}/.env ./serving-$party_id


        # serving server
        sed -i.bak "s/127.0.0.1:9380/${party_ip}:9380/g" ./serving-$party_id/confs/serving-server/conf/serving-server.properties

        # serving proxy
        sed -i.bak "s/coordinator=partyid/coordinator=${party_id}/g" ./serving-$party_id/confs/serving-proxy/conf/application.properties
        cat > ./serving-$party_id/confs/serving-proxy/conf/route_table.json <<EOF
{
    "route_table": {
$( for ((j=0;j<${#partylist[*]};j++));do
if [ "${party_id}" == "${partylist[${j}]}" ]; then
echo "        \"${partylist[${j}]}\": {
            \"default\": [
                {
                    \"ip\": \"serving-proxy\",
                    \"port\": 8059
                }
            ],
            \"serving\": [
                {
                    \"ip\": \"serving-server\",
                    \"port\": 8000
                }
            ]
        },"
else
echo "        \"${partylist[${j}]}\": {
            \"default\": [
                {
                    \"ip\": \"${servingiplist[${j}]}\",
                    \"port\": 8869
                }
            ]
        },"
fi
done)
        "default": {
            "default": [
                {
                    "ip": "default-serving-proxy",
                    "port": 8869
                }
            ]
        }
    },
    "permission": {
        "default_allow": true
    }
}
EOF
    tar -czf ./outputs/serving-$party_id.tar ./serving-$party_id
    rm -rf ./serving-$party_id
    echo serving of $party_id done!
    done
}

# only used in the k8s deployment
# TODO modularize the components

GenerateSplittingProxy() {
    # The default proxy/${party_id} port is 9370
    # "$#" return the number of args with $0 exclude
    if [ "$#" -ne 6 ]; then
        echo ""
        echo "[ERROR] Illegal number of parameters: want 8 you input $#"
        echo "The params are party_id federation_ip federation_port fate_flow_ip fate_flow_port exchange_ip"
        exit 1
    fi
    # params:
    party_id=$1
    federation_ip=$2
    federation_port=$3
    fate_flow_ip=$4
    fate_flow_port=$5
    exchange_ip=$6

    echo "Handle Splitting Proxy"
    module_name=proxy
    cd ${WORKINGDIR}
    rm -rf confs-${party_id}/
    mkdir -p confs-${party_id}/conf
    cp ${WORKINGDIR}/.env confs-${party_id}
    cp docker-compose-exchange.yml confs-${party_id}/docker-compose.yml
    cp -r docker-example-dir-tree/proxy/conf confs-${party_id}/
    if [ "$RegistryURI" != "" ]; then
        sed -i "s#PREFIX#RegistryURI#g" ./confs-${party_id}/docker-compose.yml
    fi
    sed -i "s#9371:9370#9370:9370#g" ./confs-${party_id}/docker-compose.yml
    sed -i.bak "s/port=.*/port=9370/g" ./confs-${party_id}/conf/proxy.properties
    sed -i.bak "s#route.table=.*#route.table=${deploy_dir}/proxy/conf/route_table.json#g" ./confs-${party_id}/conf/proxy.properties
    sed -i.bak "s/coordinator=.*/coordinator=${party_id}/g" ./confs-${party_id}/conf/proxy.properties
    sed -i.bak "s/ip=.*/ip=0.0.0.0/g" ./confs-${party_id}/conf/proxy.properties
    cat > ./confs-${party_id}/conf/route_table.json <<EOF
{
    "route_table": {
        "default": {
            "default": [
                {
                    "ip": "${exchange_ip}",
                    "port": 9370
                }
            ]
        },
        "${party_id}": {
            "fateflow": [
                {
                    "ip": "${fate_flow_ip}",
                    "port": ${fate_flow_port}
                }
            ],
            "fate": [
                {
                    "ip": "${federation_ip}",
                    "port": ${federation_port}
                }
            ]
        }
    },
    "permission": {
        "default_allow": true
    }
}
EOF
    tar -czf ./outputs/confs-${party_id}.tar ./confs-${party_id}
    rm -rf ./confs-${party_id}
    echo Splitting proxy of ${party_id} done!
}

ShowUsage() {
    echo "Usage: "
    echo "Generate configuration: bash generate_config.sh"
}

CleanOutputDir() {
    if [ -d ${WORKINGDIR}/outputs ];then
        rm -rf ${WORKINGDIR}/outputs
    fi
    mkdir ${WORKINGDIR}/outputs
}

main() {
    if [ "$1" = "splitting_proxy" ]; then
        CleanOutputDir
        shift
        GenerateSplittingProxy $@
    elif [ "$1" != "" ]; then
        ShowUsage
        exit 1
    else
        CleanOutputDir
        GenerateConfig
    fi

    exit 0
}

main $@
