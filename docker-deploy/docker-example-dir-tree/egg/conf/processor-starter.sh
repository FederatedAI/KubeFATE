#!/bin/bash

ONE_ARG_LIST=(
    "venv"
    "python-path"
    "engine-path"
    "port"
    "data-dir"
    "logs-dir"
    "node-manager"
    "engine-addr"
)

opts=$(getopt \
    --longoptions "$(printf "%s:," "${ONE_ARG_LIST[@]}")" \
    --name "$(basename "$0")" \
    --options "" \
    -- "$@"
)

while [[ $# -gt 0 ]]; do
   case "$1" in
        --venv)
            venv=$2
            shift 2
            ;;
        --python-path)
            python_path=$2
            shift 2
            ;;
        --engine-path)
            engine_path=$2
            shift 2
            ;;
        --port)
            port=$2
            shift 2
            ;;
        --data-dir)
            data_dir=$2
            shift 2
            ;;
        --logs-dir)
            logs_dir=$2
            shift 2
            ;;
        --node-manager)
            node_manager=$2
            shift 2
            ;;
        --engine-addr)
            engine_addr=$2
            shift 2
            ;;
        *)
            break
            ;;
    esac
done

source ${venv}/bin/activate
export PYTHONPATH=$PYTHONPATH:${python_path}
PYTHON=${venv}/bin/python
cmd="python ${engine_path} -p ${port} -d ${data_dir} -m ${node_manager} -a $(hostname -I | awk '{print $1}')"

mkdir -p ${logs_dir}
${cmd} >> ${logs_dir}/processor-${port}.log 2>&1
