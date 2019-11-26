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
source ${WORKINGDIR}/parties.conf

cd ${WORKINGDIR}

Deploy() {
	if [ "$1" = "" ];then
		echo "No party id was provided, please check your arguments "
		exit 1
	fi

	while [ "$1" != "" ]; do
    	case $1 in
			splitting_proxy)
				shift
				DeployPartyInternal $1 $2
				break
				;;
			all)
				for party in ${partylist[*]}
				do
				    DeployPartyInternal $party
				done

				DeployPartyInternal exchange
				break
				;;
         	*)
				DeployPartyInternal $1
            	;;
    	esac
    	shift

	done
}

DeployPartyInternal() {
	target_party_id=$1
	# should not use localhost at any case
	target_party_ip="127.0.0.1"

	# check configuration files
	if [ ! -d ${WORKINGDIR}/outputs ];then
		echo "Unable to find outputs dir, please generate config files first."
		exit 1
	fi
	if [ ! -f ${WORKINGDIR}/outputs/confs-${target_party_id}.tar ];then
		echo "Unable to find deployment file for party $target_party_id, please generate it first."
		exit 1
	fi
	# extract the ip address of the target party
	if [ "$target_party_id" = "exchange" ];then
		target_party_ip=${exchangeip}
	elif [ "$2" != "" ]; then
		target_party_ip="$2"
	else
		for ((i=0;i<${#partylist[*]};i++))
		do
			if [ "${partylist[$i]}" = "$target_party_id" ];then
				target_party_ip=${partyiplist[$i]}
			fi
		done
	fi
	# verify the target_party_ip
	if [ "$target_party_ip" = "127.0.0.1" ]; then
		echo "Unable to find Party: $target_party_id, please check you input."
		exit 1
	fi
    scp ${WORKINGDIR}/outputs/confs-$target_party_id.tar $user@$target_party_ip:~/
    #rm -f ${WORKINGDIR}/outputs/confs-$target_party_id.tar
    echo "$target_party_ip copy is ok!"
    ssh -tt $user@$target_party_ip<< eeooff
mkdir -p $dir
mv ~/confs-$target_party_id.tar $dir
cd $dir
tar -xzf confs-$target_party_id.tar
cd confs-$target_party_id
docker-compose down
docker-compose up -d
cd ../
rm -f confs-$target_party_id.tar
exit
eeooff
    echo "party $target_party_id deploy is ok!"
}

ShowUsage() {
	echo "Usage: "
    echo "Deploy all parties or specified partie(s): bash docker_deploy.sh partyid1[partyid2...] | all"
}

main() {
	if [ "$1" = "" ] || [ "$" = "--help" ]; then
		ShowUsage
		exit 1
	else
		Deploy "$@"
	fi

	exit 0
}

main $@
