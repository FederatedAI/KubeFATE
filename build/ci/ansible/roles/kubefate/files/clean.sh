#/bin/bash

clean_when_timeout()
{
  rm -rf ${BASE_DIR}/*

  # Install Kind
  kind_status=`kind version`
  if [ $? -eq 0 ]; then
    echo "Deleting kind cluster..." 
    kind delete cluster
  fi

  echo "exit because of task timeout"
  exit 1
}

clean_when_timeout
