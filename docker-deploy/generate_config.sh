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
    
        eval nodemanager_ip=(nodemanager)
        eval nodemanager_port=4671
        eval nodemanager_port_db=9461
    
        eval clustermanager_ip=clustermanager
        eval clustermanager_port=4670
        eval clustermanager_port_db=9460

        eval proxy_ip=rollsite
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
	# use the fixed db name here
        sed -i "s#<jdbc.url>#jdbc:mysql://${db_ip}:3306/eggroll_meta?useSSL=false\&serverTimezone=UTC\&characterEncoding=utf8\&allowPublicKeyRetrieval=true#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
	sed -i "s#<jdbc.username>#${db_user}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
        sed -i "s#<jdbc.password>#${db_password}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties

        #clustermanager & nodemanager
	sed -i "s#<clustermanager.host>#${clustermanager_ip}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties 
	sed -i "s#<clustermanager.port>#${clustermanager_port}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties  
	sed -i "s#<nodemanager.port>#${nodemanager_port}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
	sed -i "s#<party.id>#${party_id}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties

        #python env
        sed -i "s#<venv>#${venv_dir}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
        #pythonpath, very import, do not modify."
	sed -i "s#<python.path>#/data/projects/fate/python:/data/projects/fate/eggroll/python#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties

        #javahome
        sed -i "s#<java.home>#/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.252.b09-2.el7_8.x86_64/jre#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
	sed -i "s#<java.classpath>#conf/:lib/*#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties

	sed -i "s#<rollsite.host>#${proxy_ip}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
	sed -i "s#<rollsite.port>#${proxy_port}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties

        # fateboard
        sed -i "s#^server.port=.*#server.port=${fateboard_port}#g" ./confs-$party_id/confs/fateboard/conf/application.properties
        sed -i "s#^fateflow.url=.*#fateflow.url=http://${fate_flow_ip}:${fate_flow_http_port}#g" ./confs-$party_id/confs/fateboard/conf/application.properties
        sed -i "s#<jdbc.username>#${db_user}#g" ./confs-$party_id/confs/fateboard/conf/application.properties
        sed -i "s#<jdbc.password>#${db_password}#g" ./confs-$party_id/confs/fateboard/conf/application.properties
        sed -i "s#<jdbc.url>#jdbc:mysql://${db_ip}:3306/${db_name}?characterEncoding=utf8\&characterSetResults=utf8\&autoReconnect=true\&failOverReadOnly=false\&serverTimezone=GMT%2B8#g" ./confs-$party_id/confs/fateboard/conf/application.properties
        echo fateboard module of $party_id done!
    
        # mysql
        # sed -i "s/eggroll_meta/${db_name}/g" ./confs-$party_id/confs/mysql/init/create-eggroll-meta-tables.sql
        echo > ./confs-$party_id/confs/mysql/init/insert-node.sql
	echo "GRANT ALL ON *.* TO '${db_user}'@'%';" >> ./confs-$party_id/confs/mysql/init/insert-node.sql
	echo 'USE `eggroll_meta`;' >> ./confs-$party_id/confs/mysql/init/insert-node.sql
        echo "INSERT INTO server_node (host, port, node_type, status) values ('${clustermanager_ip}', '${clustermanager_port_db}', 'CLUSTER_MANAGER', 'HEALTHY');" >> ./confs-$party_id/confs/mysql/init/insert-node.sql
        for ((j=0;j<${#nodemanager_ip[*]};j++))
        do
            echo "INSERT INTO server_node (host, port, node_type, status) values ('${nodemanager_ip[j]}', '${nodemanager_port_db}', 'NODE_MANAGER', 'HEALTHY');" >> ./confs-$party_id/confs/mysql/init/insert-node.sql
        done
        echo "show tables;" >> ./confs-$party_id/confs/mysql/init/insert-node.sql
        echo "select * from node;" >> ./confs-$party_id/confs/mysql/init/insert-node.sql
        echo mysql module of $party_id done!

        # fate_flow
	sed -i "s/WORK_MODE =.*/WORK_MODE = 1/g" ./confs-$party_id/confs/fate_flow/conf/settings.py
        sed -i "s/user:.*/user: '${db_user}'/g" ./confs-$party_id/confs/fate_flow/conf/base_conf.yaml
        sed -i "s/passwd:.*/passwd: '${db_password}'/g" ./confs-$party_id/confs/fate_flow/conf/base_conf.yaml
        sed -i "s/host: 192.168.0.1*/host: '${db_ip}'/g" ./confs-$party_id/confs/fate_flow/conf/base_conf.yaml
        sed -i "s/name:.*/name: '${db_name}'/g" ./confs-$party_id/confs/fate_flow/conf/base_conf.yaml
        sed -i "s/serving:8000/${serving_ip}:8000/g" ./confs-$party_id/confs/fate_flow/conf/server_conf.json
    
        cat > ./confs-$party_id/confs/fate_flow/conf/server_conf.json <<EOF
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
        "fateflow": {
            "host": "${fate_flow_ip}",
            "grpc.port": ${fate_flow_grpc_port},
            "http.port": ${fate_flow_http_port}
        },
        "servings": [
          "serving:8000"
        ]
    }
}
EOF
        echo fate_flow module of $party_id done!
        # rollsite 
    	cat > ./confs-$party_id/confs/eggroll/conf/route_table.json <<EOF
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
            }]
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
  		sed -i "s/port=.*/port=${proxy_port}/g" ./confs-exchange/conf/proxy.properties
  		sed -i "s#route.table=.*#route.table=${deploy_dir}/proxy/conf/route_table.json#g" ./confs-exchange/conf/proxy.properties
  		sed -i "s/coordinator=.*/coordinator=${party_id}/g" ./confs-exchange/conf/proxy.properties
  		sed -i "s/ip=.*/ip=0.0.0.0/g" ./confs-exchange/conf/proxy.properties
  		
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
	done
  
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
    	    sed -i "s/127.0.0.1:9380/${party_ip}:9380/g" ./serving-$party_id/confs/serving-server/conf/serving-server.properties

    	    # serving proxy
    	    sed -i "s/coordinator=partyid/coordinator=${party_id}/g" ./serving-$party_id/confs/serving-proxy/conf/application.properties
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
    sed -i "s/port=.*/port=9370/g" ./confs-${party_id}/conf/proxy.properties
    sed -i "s#route.table=.*#route.table=${deploy_dir}/proxy/conf/route_table.json#g" ./confs-${party_id}/conf/proxy.properties
    sed -i "s/coordinator=.*/coordinator=${party_id}/g" ./confs-${party_id}/conf/proxy.properties
    sed -i "s/ip=.*/ip=0.0.0.0/g" ./confs-${party_id}/conf/proxy.properties
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
