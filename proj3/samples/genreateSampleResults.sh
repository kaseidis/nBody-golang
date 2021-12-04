#!/bin/bash

if [[ -z "${NBODY_THREAD}" ]]; then
  THREAD_COUNT="4"
else
  THREAD_COUNT="${NBODY_THREAD}"
fi

echo "Generating result from test1.json file"
go run proj3/main $THREAD_COUNT < test1.json > test1.result.json 

echo "Generating result from test2.json file"
go run proj3/main $THREAD_COUNT < test2.json > test2.result.json 

echo "Generating result from test3.json file"
go run proj3/main $THREAD_COUNT < test3.json > test3.result.json

echo "Results generated, please using result_visulize.htm to visulize result"