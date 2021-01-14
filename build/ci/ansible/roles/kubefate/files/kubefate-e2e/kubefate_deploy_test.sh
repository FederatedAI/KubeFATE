#!/bin/bash

dir=$(cd $(dirname $0) && pwd)
source $dir/color.sh

source $dir/common.sh

binary_install

if check_kubectl; then
  loginfo "kubectl ready"
else
  exit 1
fi

if kubefate_install; then
  loginfo "kubefate install success"
else
  exit 1
fi

set_host

if check_kubefate_version; then
  loginfo "kubefate CLI ready"
else
  exit 1
fi

kubefate_uninstall

clean_host

loginfo "kubefate_deploy_test done."

exit 0
