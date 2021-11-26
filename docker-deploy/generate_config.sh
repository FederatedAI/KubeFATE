#!/bin/bash

# Copyright 2019-2020 VMware, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# you may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e
BASEDIR=$(dirname "$0")
cd $BASEDIR
WORKINGDIR=$(pwd)
deploy_dir=/data/projects/fate

# fetch fate-python image
source ${WORKINGDIR}/.env
source ${WORKINGDIR}/parties.conf

echo "Generate Config"
echo "Info:"
echo "  RegistryURI: ${RegistryURI}"
echo "  Tag: ${TAG}"
echo "  Serving_Tag: ${SERVING_TAG}"
echo "  Backend: ${backend}"

echo "  Party_List:"
for ((i = 0; i < ${#party_list[*]}; i++)); do
	echo "  - Party_ID:${party_list[${i}]}"
	echo "    Party_IP:${party_ip_list[${i}]}"
done
echo ""
echo ""
cd ${WORKINGDIR}

GenerateConfig() {
	for ((i = 0; i < ${#party_list[*]}; i++)); do
		eval party_id=\${party_list[${i}]}
		eval party_ip=\${party_ip_list[${i}]}
		eval serving_ip=\${serving_ip_list[${i}]}

		eval processor_count=2
		eval venv_dir=/data/projects/python/venv
		eval python_path=${deploy_dir}/python:${deploy_dir}/eggroll/python
		eval data_dir=${deploy_dir}/data-dir

		eval nodemanager_ip=nodemanager
		eval nodemanager_port=4671
		eval nodemanager_port_db=9461

		eval clustermanager_ip=clustermanager
		eval clustermanager_port=4670
		eval clustermanager_port_db=9460

		eval proxy_ip=rollsite
		eval proxy_port=9370

		eval fateboard_ip=fateboard
		eval fateboard_port=8080
		eval fateboard_username=${fateboard_username}
		eval fateboard_password=${fateboard_password}

		eval fate_flow_ip=python
		eval fate_flow_grpc_port=9360
		eval fate_flow_http_port=9380
		eval fml_agent_port=8484

		eval db_ip=${mysql_ip}
		eval db_user=${mysql_user}
		eval db_password=${mysql_password}
		eval db_name=${mysql_db}

		eval exchange_ip=${exchangeip}

		rm -rf confs-$party_id/
		mkdir -p confs-$party_id/confs
		cp -r training_template/public/* confs-$party_id/confs/
		# handle spark backend here
		if [ "$backend" == "spark_rabbitmq" ]; then
			cp -r training_template/backends/spark/hadoop confs-$party_id/confs/
			cp -r training_template/backends/spark/nginx confs-$party_id/confs/
			cp -r training_template/backends/spark/spark confs-$party_id/confs/
			cp -r training_template/backends/spark/rabbitmq confs-$party_id/confs/
			
			cp training_template/docker-compose-spark.yml confs-$party_id/docker-compose.yml
			sed -i '155,170d' confs-$party_id/docker-compose.yml
		fi

		if [ "$backend" == "spark_pulsar" ]; then
			cp -r training_template/backends/spark/hadoop confs-$party_id/confs/
			cp -r training_template/backends/spark/nginx confs-$party_id/confs/
			cp -r training_template/backends/spark/spark confs-$party_id/confs/
			cp -r training_template/backends/spark/pulsar confs-$party_id/confs/
			
			cp training_template/docker-compose-spark.yml confs-$party_id/docker-compose.yml
			sed -i '138,153d' confs-$party_id/docker-compose.yml
		fi

		if [ "$backend" == "eggroll" ]; then
			# if the computing backend is not spark, use eggroll anyway
			cp -r training_template/backends/eggroll confs-$party_id/confs/
			cp training_template/docker-compose-eggroll.yml confs-$party_id/docker-compose.yml

			# eggroll config
			#db connect inf
			# use the fixed db name here
			sed -i "s#<jdbc.url>#jdbc:mysql://${db_ip}:3306/${db_name}?useSSL=false\&serverTimezone=UTC\&characterEncoding=utf8\&allowPublicKeyRetrieval=true#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
			sed -i "s#<jdbc.username>#${db_user}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
			sed -i "s#<jdbc.password>#${db_password}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties

			#clustermanager & nodemanager
			sed -i "s#<clustermanager.host>#${clustermanager_ip}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
			sed -i "s#<clustermanager.port>#${clustermanager_port}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
			sed -i "s#<nodemanager.port>#${nodemanager_port}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
			sed -i "s#<party.id>#${party_id}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties

			#pythonpath, very import, do not modify."
			sed -i "s#<python.path>#/data/projects/fate/python:/data/projects/fate/eggroll/python#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties

			#javahome
			sed -i "s#<java.home>#/usr/lib/jvm/java-1.8.0-openjdk#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
			sed -i "s#<java.classpath>#conf/:lib/*#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties

			sed -i "s#<rollsite.host>#${proxy_ip}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
			sed -i "s#<rollsite.port>#${proxy_port}#g" ./confs-$party_id/confs/eggroll/conf/eggroll.properties
		fi

		# generate conf dir
		cp ${WORKINGDIR}/.env ./confs-$party_id
		
		# check if use python-nn
		if [ "$enabled_nn" = "true" ]; then
			sed -i 's#image: "federatedai/python:${TAG}"#image: "federatedai/python-nn:${TAG}"#g' ./confs-$party_id/docker-compose.yml
		fi

		if [ "$RegistryURI" != "" ]; then
			sed -i 's#federatedai#${RegistryURI}/federatedai#g' ./confs-$party_id/docker-compose.yml
			sed -i 's#image: "mysql:8"#image: ${RegistryURI}/federatedai/mysql:8#g' ./confs-$party_id/docker-compose.yml
		fi

		# replace namenode in training_template/public/fate_flow/conf/service_conf.yaml
		if [ "$name_node" != "" ]; then
			sed -i "s#name_node: hdfs://namenode:9000#name_node: ${name_node}#g" ./confs-$party_id/confs/fate_flow/conf/service_conf.yaml
		fi

		# update serving ip
		sed -i "s/fate-serving/${serving_ip}/g" ./confs-$party_id/docker-compose.yml

		# update the path of shared_dir
		shared_dir="confs-${party_id}/shared_dir"

		# create directories
		for value in "examples" "federatedml" "data"; do
			mkdir -p ${shared_dir}/${value}
		done

		sed -i "s|{/path/to/host/dir}|${dir}/${shared_dir}|g" ./confs-$party_id/docker-compose.yml

		# Start the general config rendering
		# fateboard
		sed -i "s#^server.port=.*#server.port=${fateboard_port}#g" ./confs-$party_id/confs/fateboard/conf/application.properties
		sed -i "s#^fateflow.url=.*#fateflow.url=http://${fate_flow_ip}:${fate_flow_http_port}#g" ./confs-$party_id/confs/fateboard/conf/application.properties
		sed -i "s#<jdbc.username>#${db_user}#g" ./confs-$party_id/confs/fateboard/conf/application.properties
		sed -i "s#<jdbc.password>#${db_password}#g" ./confs-$party_id/confs/fateboard/conf/application.properties
		sed -i "s#<jdbc.url>#jdbc:mysql://${db_ip}:3306/${db_name}?characterEncoding=utf8\&characterSetResults=utf8\&autoReconnect=true\&failOverReadOnly=false\&serverTimezone=GMT%2B8#g" ./confs-$party_id/confs/fateboard/conf/application.properties
		sed -i "s#<fateboard.username>#${fateboard_username}#g" ./confs-$party_id/confs/fateboard/conf/application.properties
		sed -i "s#<fateboard.password>#${fateboard_password}#g" ./confs-$party_id/confs/fateboard/conf/application.properties
		echo fateboard module of $party_id done!

		# mysql
		
		echo >./confs-$party_id/confs/mysql/init/insert-node.sql
		echo "CREATE DATABASE IF NOT EXISTS ${db_name};" >>./confs-$party_id/confs/mysql/init/insert-node.sql
		echo "CREATE USER '${db_user}'@'%' IDENTIFIED BY '${db_password}';" >>./confs-$party_id/confs/mysql/init/insert-node.sql
		echo "GRANT ALL ON *.* TO '${db_user}'@'%';" >>./confs-$party_id/confs/mysql/init/insert-node.sql
		
		if [[ "$backend" == "eggroll" ]]; then
			echo 'USE `'${db_name}'`;' >>./confs-$party_id/confs/mysql/init/insert-node.sql
			echo "INSERT INTO server_node (host, port, node_type, status) values ('${clustermanager_ip}', '${clustermanager_port_db}', 'CLUSTER_MANAGER', 'HEALTHY');" >>./confs-$party_id/confs/mysql/init/insert-node.sql
			for ((j = 0; j < ${#nodemanager_ip[*]}; j++)); do
				echo "INSERT INTO server_node (host, port, node_type, status) values ('${nodemanager_ip[j]}', '${nodemanager_port_db}', 'NODE_MANAGER', 'HEALTHY');" >>./confs-$party_id/confs/mysql/init/insert-node.sql
			done
			echo "show tables;" >>./confs-$party_id/confs/mysql/init/insert-node.sql
			echo "select * from server_node;" >>./confs-$party_id/confs/mysql/init/insert-node.sql
		
			sed -i "s/eggroll_meta/${db_name}/g" ./confs-$party_id/confs/mysql/init/create-eggroll-meta-tables.sql
		else
			rm -f ./confs-$party_id/confs/mysql/init/create-eggroll-meta-tables.sql
		fi
		echo mysql module of $party_id done!

		# fate_flow
		sed -i "s/name: <db_name>/name: '${db_name}'/g" ./confs-$party_id/confs/fate_flow/conf/service_conf.yaml
		sed -i "s/user: <db_user>/user: '${db_user}'/g" ./confs-$party_id/confs/fate_flow/conf/service_conf.yaml
		sed -i "s/passwd: <db_passwd>/passwd: '${db_password}'/g" ./confs-$party_id/confs/fate_flow/conf/service_conf.yaml
		sed -i "s/host: <db_host>/host: '${db_ip}'/g" ./confs-$party_id/confs/fate_flow/conf/service_conf.yaml
		sed -i "s/127.0.0.1:8000/${serving_ip}:8000/g" ./confs-$party_id/confs/fate_flow/conf/service_conf.yaml


		if [[ "$backend" == "spark_rabbitmq" ]]; then
			sed -i "s/proxy: rollsite/proxy: nginx/g" ./confs-$party_id/confs/fate_flow/conf/service_conf.yaml
			# 
			sed -i "s/  computing: eggroll/  computing: spark/g" ./confs-$party_id/confs/fate_flow/conf/service_conf.yaml
			sed -i "s/  federation: eggroll/  federation: rabbitmq/g" ./confs-$party_id/confs/fate_flow/conf/service_conf.yaml
			sed -i "s/  storage: eggroll/  storage: hdfs/g" ./confs-$party_id/confs/fate_flow/conf/service_conf.yaml
		fi

		if [[ "$backend" == "spark_pulsar" ]]; then
			sed -i "s/proxy: rollsite/proxy: nginx/g" ./confs-$party_id/confs/fate_flow/conf/service_conf.yaml
			# 
			sed -i "s/  computing: eggroll/  computing: spark/g" ./confs-$party_id/confs/fate_flow/conf/service_conf.yaml
			sed -i "s/  federation: eggroll/  federation: pulsar/g" ./confs-$party_id/confs/fate_flow/conf/service_conf.yaml
			sed -i "s/  storage: eggroll/  storage: hdfs/g" ./confs-$party_id/confs/fate_flow/conf/service_conf.yaml
		fi

		echo fate_flow module of $party_id done!

		# now we handles the route table
		# eggroll
		if [[ "$backend" == "eggroll" ]]; then
			cat >./confs-$party_id/confs/eggroll/conf/route_table.json <<EOF
{
	"route_table": {
		"default": {
			"default": [
				{
$(if [ "$exchange_ip" != "" ]; then
				echo "
				\"ip\": \"${exchange_ip}\",
				\"port\": 9371
	"
			else
				echo " 
				\"ip\": \"${proxy_ip}\",
				\"port\": \"${proxy_port}\"
	"
			fi)
				}
			]
		},
$(for ((j = 0; j < ${#party_list[*]}; j++)); do
				if [ "${party_id}" == "${party_list[${j}]}" ]; then
					continue
				fi
				echo "
		\"${party_list[${j}]}\": {
			\"default\": [{
		 		\"ip\": \"${party_ip_list[${j}]}\",
				\"port\": 9370
			    }]
		},
	"
			done)
		"${party_id}": {
			"default": [{
				"ip": "${proxy_ip}",
				"port": ${proxy_port}
			}],
			"fateflow": [{
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
		fi

		# nginx
        # TODO nginx 不需要对方fateflow IP PORT
		if [[ "$backend" == "spark"* ]]; then
			cat >./confs-$party_id/confs/nginx/route_table.yaml <<EOF
default:
  proxy:
    - host: nginx
      http_port: 9300
      grpc_port: 9310
$(for ((j = 0; j < ${#party_list[*]}; j++)); do
				if [ "${party_id}" == "${party_list[${j}]}" ]; then
					continue
				fi
				echo "${party_list[${j}]}:
  proxy:
    - host: ${party_ip_list[${j}]} 
      http_port: 9300
      grpc_port: 9310
  fateflow:
    - host: ${party_ip_list[${j}]}
      grpc_port: ${fate_flow_grpc_port}
      http_port: ${fate_flow_http_port}
"
			done)
${party_id}:
  proxy:
    - host: nginx
      http_port: 9300
      grpc_port: 9310
  fateflow:
    - host: ${fate_flow_ip}
      grpc_port: ${fate_flow_grpc_port}
      http_port: ${fate_flow_http_port}
EOF
		fi

		# spark_pulsar
		if [[ "$backend" == "spark_pulsar" ]]; then
			cat >./confs-$party_id/confs/fate_flow/conf/pulsar_route_table.yaml <<EOF
$(for ((j = 0; j < ${#party_list[*]}; j++)); do
				if [ "${party_id}" == "${party_list[${j}]}" ]; then
					continue
				fi
				echo "${party_list[${j}]}:
    host: ${party_ip_list[${j}]}
    port: 6650
"
			done)
${party_id}:
    host: pulsar
    port: 6650
EOF

		fi

		# spark_rabbitmq
		if [[ "$backend" == "spark_rabbitmq" ]]; then
			cat >./confs-$party_id/confs/fate_flow/conf/rabbitmq_route_table.yaml <<EOF
$(for ((j = 0; j < ${#party_list[*]}; j++)); do
				if [ "${party_id}" == "${party_list[${j}]}" ]; then
					continue
				fi
				echo "${party_list[${j}]}:
    host: ${party_ip_list[${j}]}
    port: 5672
"
			done)
${party_id}:
    host: rabbitmq
    port: 5672
EOF
		fi
		echo proxy module of $party_id done!

		# package of $party_id
		tar -czf ./outputs/confs-$party_id.tar ./confs-$party_id
		rm -rf ./confs-$party_id
		echo package $party_id done!

		if [ "$exchange_ip" != "" ]; then
			# handle exchange
			echo "handle exchange"
			module_name=exchange
			cd ${WORKINGDIR}
			rm -rf confs-exchange/
			mkdir -p confs-exchange/conf/
			cp ${WORKINGDIR}/.env confs-exchange/
			cp training_template/docker-compose-exchange.yml confs-exchange/docker-compose.yml
			cp -r training_template/backends/eggroll/conf/* confs-exchange/conf/

			if [ "$RegistryURI" != "" ]; then
				sed -i 's#federatedai#${RegistryURI}/federatedai#g' ./confs-exchange/docker-compose.yml
			fi

			sed -i "s#<rollsite.host>#${proxy_ip}#g" ./confs-exchange/conf/eggroll.properties
			sed -i "s#<rollsite.port>#${proxy_port}#g" ./confs-exchange/conf/eggroll.properties
			sed -i "s#<party.id>#exchange#g" ./confs-exchange/conf/eggroll.properties
			sed -i "s/coordinator=.*/coordinator=exchange/g" ./confs-exchange/conf/eggroll.properties
			sed -i "s/ip=.*/ip=0.0.0.0/g" ./confs-exchange/conf/eggroll.properties

			cat >./confs-exchange/conf/route_table.json <<EOF
{
    "route_table": {
$(for ((j = 0; j < ${#party_list[*]}; j++)); do
				echo "        \"${party_list[${j}]}\": {
            \"default\": [
                {
                    \"ip\": \"${party_ip_list[${j}]}\",
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
	for ((i = 0; i < ${#serving_ip_list[*]}; i++)); do
		eval party_id=\${party_list[${i}]}
		eval party_ip=\${party_ip_list[${i}]}
		eval serving_ip=\${serving_ip_list[${i}]}

		rm -rf serving-$party_id/
		mkdir -p serving-$party_id/confs
		cp -r serving_template/docker-serving/* serving-$party_id/confs/

		cp serving_template/docker-compose-serving.yml serving-$party_id/docker-compose.yml
		if [ "$RegistryURI" != "" ]; then
			sed -i 's#federatedai#${RegistryURI}/federatedai#g' ./serving-$party_id/docker-compose.yml
			# should not use federatedai in here
			sed -i 's#image: "redis:5"#image: "${RegistryURI}/federatedai/redis:5"#g' ./serving-$party_id/docker-compose.yml
		fi
		# generate conf dir
		cp ${WORKINGDIR}/.env ./serving-$party_id

		# serving server
		sed -i "s/127.0.0.1:9380/${party_ip}:9380/g" ./serving-$party_id/confs/serving-server/conf/serving-server.properties
		sed -i "s/<redis.ip>/${redis_ip}/g" ./serving-$party_id/confs/serving-server/conf/serving-server.properties
		sed -i "s/<redis.port>/${redis_port}/g" ./serving-$party_id/confs/serving-server/conf/serving-server.properties
		sed -i "s/<redis.password>/${redis_password}/g" ./serving-$party_id/confs/serving-server/conf/serving-server.properties
		sed -i "s/<redis.password>/${redis_password}/g" ./serving-$party_id/docker-compose.yml

		# network
		sed -i "s/name: <fate-network>/name: confs-${party_id}_fate-network/g" serving-$party_id/docker-compose.yml
		
		

		# serving proxy
		sed -i "s/coordinator=9999/coordinator=${party_id}/g" ./serving-$party_id/confs/serving-proxy/conf/application.properties
		cat >./serving-$party_id/confs/serving-proxy/conf/route_table.json <<EOF
{
    "route_table": {
$(for ((j = 0; j < ${#party_list[*]}; j++)); do
			if [ "${party_id}" == "${party_list[${j}]}" ]; then
				echo "        \"${party_list[${j}]}\": {
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
				echo "        \"${party_list[${j}]}\": {
            \"default\": [
                {
                    \"ip\": \"${serving_ip_list[${j}]}\",
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

ShowUsage() {
	echo "Usage: "
	echo "Generate configuration: bash generate_config.sh"
}

CleanOutputDir() {
	if [ -d ${WORKINGDIR}/outputs ]; then
		rm -rf ${WORKINGDIR}/outputs
	fi
	mkdir ${WORKINGDIR}/outputs
}

main() {
	if [ "$1" != "" ]; then
		ShowUsage
		exit 1
	else
		CleanOutputDir
		GenerateConfig
	fi

	exit 0
}

main $@
