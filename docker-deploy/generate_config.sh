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
        eval version=1.2
        eval java_dir=
        eval source_code_dir=
        eval output_packages_dir=
        eval deploy_packages_dir=
    
        eval party_id=\${partylist[${i}]}
        eval party_ip=\${partyiplist[${i}]}
        eval serving_ip=\${servingiplist[${i}]}
    
        eval processor_port=50000
        eval processor_count=16
        eval venv_dir=/data/projects/python/venv
        eval python_path=${deploy_dir}/python:${deploy_dir}/eggroll/python
        eval data_dir=${deploy_dir}/data-dir
    
        eval egg_ip=(egg)
        eval egg_port=7888
    
        eval meta_service_ip=meta-service
        eval meta_service_port=8590
    
        eval roll_ip=roll
        eval roll_port=8011
    
        eval proxy_ip=proxy
        eval proxy_port=9370
    
        eval fateboard_ip=fateboard
        eval fateboard_port=8080
    
        eval fate_flow_ip=python
        eval fate_flow_grpc_port=9360
        eval fate_flow_http_port=9380
    
        eval storage_service_ip=(egg)
        eval storage_service_port=7778
    
        eval db_ip=mysql
        eval db_user=fate
        eval db_password=fate_dev
        eval db_name=fate
    
        eval node_list=()
    
        eval redis_ip=redis
        eval redis_password=fate_dev
    
        eval exchange_ip=${exchangeip}
        eval federation_ip=federation
        eval federation_port=9394
        eval serving_ip1=
        eval eval serving_ip2=
    
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
          sed -i 's#image: "redis:5"#image: "${RegistryURI}/redis:5"#g' ./confs-$party_id/docker-compose.yml
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
        sed -i.bak "s/party.id=.*/party.id=${party_id}/g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s/service.port=.*/service.port=${egg_port}/g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s/engine.names=.*/engine.names=processor/g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s#bootstrap.script=.*#bootstrap.script=${deploy_dir}/${module_name}/egg/conf/processor-starter.sh#g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s#start.port=.*#start.port=${processor_port}#g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s#processor.venv=.*#processor.venv=${venv_dir}#g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s#processor.python-path=.*#processor.python-path=${python_path}#g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s#processor.engine-path=.*#processor.engine-path=${deploy_dir}/eggroll/python/eggroll/computing/processor.py#g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s#data-dir=.*#data-dir=${data_dir}#g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s#processor.logs-dir=.*#processor.logs-dir=${deploy_dir}/eggroll/python/eggroll/logs/processor#g" ./confs-$party_id/confs/egg/conf/egg.properties
        sed -i.bak "s#count=.*#count=${processor_count}#g" ./confs-$party_id/confs/egg/conf/egg.properties
        echo  >> ./confs-$party_id/confs/egg/conf/egg.properties
        echo "eggroll.computing.processor.python-path=${python_path}" >> ./confs-$party_id/confs/egg/conf/egg.properties
        echo egg module of $party_id done!
    
        # meta-service
        sed -i.bak "s/party.id=.*/party.id=${party_id}/g" ./confs-$party_id/confs/meta-service/conf/meta-service.properties
        sed -i.bak "s/service.port=.*/service.port=${meta_service_port}/g" ./confs-$party_id/confs/meta-service/conf/meta-service.properties
        sed -i.bak "s#//.*?#//${db_ip}:3306/${db_name}?#g" ./confs-$party_id/confs/meta-service/conf/meta-service.properties
        sed -i.bak "s/jdbc.username=.*/jdbc.username=${db_user}/g" ./confs-$party_id/confs/meta-service/conf/meta-service.properties
        sed -i.bak "s/jdbc.password=.*/jdbc.password=${db_password}/g" ./confs-$party_id/confs/meta-service/conf/meta-service.properties
        echo meta-service module of $party_id done!
    
        # roll
        sed -i.bak "s/party.id=.*/party.id=${party_id}/g" ./confs-$party_id/confs/roll/conf/roll.properties
        sed -i.bak "s/service.port=.*/service.port=${roll_port}/g" ./confs-$party_id/confs/roll/conf/roll.properties
        sed -i.bak "s/meta.service.ip=.*/meta.service.ip=${meta_service_ip}/g" ./confs-$party_id/confs/roll/conf/roll.properties
        sed -i.bak "s/meta.service.port=.*/meta.service.port=${meta_service_port}/g" ./confs-$party_id/confs/roll/conf/roll.properties
        echo roll module of $party_id done!
    
        # fateboard
        sed -i.bak "s#^server.port=.*#server.port=${fateboard_port}#g" ./confs-$party_id/confs/fateboard/conf/application.properties
        sed -i.bak "s#^fateflow.url=.*#fateflow.url=http://${fate_flow_ip}:${fate_flow_http_port}#g" ./confs-$party_id/confs/fateboard/conf/application.properties
        sed -i.bak "s#^spring.datasource.driver-Class-Name=.*#spring.datasource.driver-Class-Name=com.mysql.cj.jdbc.Driver#g" ./confs-$party_id/confs/fateboard/conf/application.properties
        sed -i.bak "s#^spring.datasource.url=.*#spring.datasource.url=jdbc:mysql://${db_ip}:3306/${db_name}?characterEncoding=utf8\&characterSetResults=utf8\&autoReconnect=true\&failOverReadOnly=false\&serverTimezone=GMT%2B8#g" ./confs-$party_id/confs/fateboard/conf/application.properties
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
    
        # redis
        sed -i.bak "s/bind 127.0.0.1/bind 0.0.0.0/g" ./confs-$party_id/confs/redis/conf/redis.conf
        sed -i.bak "s/# requirepass foobared/requirepass ${redis_password}/g" ./confs-$party_id/confs/redis/conf/redis.conf
        sed -i.bak "s/databases 16/databases 50/g" ./confs-$party_id/confs/redis/conf/redis.conf
        echo redis module of $party_id done!
    
        # fate_flow
        sed -i.bak "s/WORK_MODE =.*/WORK_MODE = 1/g" ./confs-$party_id/confs/fate_flow/conf/settings.py
        sed -i.bak "s/'user':.*/'user': '${db_user}',/g" ./confs-$party_id/confs/fate_flow/conf/settings.py
        sed -i.bak "s/'passwd':.*/'passwd': '${db_password}',/g" ./confs-$party_id/confs/fate_flow/conf/settings.py
        sed -i.bak "s/'host':.*/'host': '${db_ip}',/g" ./confs-$party_id/confs/fate_flow/conf/settings.py
        sed -i.bak "s/'name':.*/'name': '${db_name}',/g" ./confs-$party_id/confs/fate_flow/conf/settings.py
        sed -i.bak "s/'password':.*/'password': '${redis_password}',/g" ./confs-$party_id/confs/fate_flow/conf/settings.py
        sed -i.bak "/'host':.*/{x;s/^/./;/^\.\{2\}$/{x;s/.*/    'host': '${redis_ip}',/;x};x;}" ./confs-$party_id/confs/fate_flow/conf/settings.py
        sed -i.bak "s/serving:8000/${serving_ip}:8000/g" ./confs-$party_id/confs/fate_flow/conf/server_conf.json
        echo fate_flow module of $party_id done!
    
        # federatedml
        cat > ./confs-$party_id/confs/federatedml/conf/server_conf.json <<EOF
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
        echo federatedml module of $party_id done!
    
        # federation
        sed -i.bak "s/party.id=.*/party.id=${party_id}/g" ./confs-$party_id/confs/federation/conf/federation.properties
        sed -i.bak "s/service.port=.*/service.port=${federation_port}/g" ./confs-$party_id/confs/federation/conf/federation.properties
        sed -i.bak "s/meta.service.ip=.*/meta.service.ip=${meta_service_ip}/g" ./confs-$party_id/confs/federation/conf/federation.properties
        sed -i.bak "s/meta.service.port=.*/meta.service.port=${meta_service_port}/g" ./confs-$party_id/confs/federation/conf/federation.properties
        sed -i.bak "s/proxy.ip=.*/proxy.ip=${proxy_ip}/g" ./confs-$party_id/confs/federation/conf/federation.properties
        sed -i.bak "s/proxy.port=.*/proxy.port=${proxy_port}/g" ./confs-$party_id/confs/federation/conf/federation.properties
        echo federation module of $party_id done!
    
        # proxy
        module_name="proxy"
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
