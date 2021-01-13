#/bin/bash

clean()
{
  rm -rf ${BASE_DIR}/*

  # Delete kind
  kind_status=`kind version`
  if [ $? -eq 0 ]; then
    echo "Deleting kind cluster..." 
    kind delete cluster
  fi
}

main()
{
  if [ "$1" != "" ]; then
    if [ "$1" == "failed" ]; then
      clean
      echo "exit with errors"
      exit 1
    fi
  else
    clean
  fi
}

main $1
