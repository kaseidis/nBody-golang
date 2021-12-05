#!/bin/bash

# Check if NBODY_THREAD is defined in environment
if [[ -z "${NBODY_THREAD}" ]]; then
  THREAD_COUNT="4"
else
  THREAD_COUNT="${NBODY_THREAD}"
fi

# Run program
echo "Running in %THREAD_COUNT% threads, you can set up environment varible 'NBODY_THREAD' to change this"
go run proj3/main $THREAD_COUNT < test1.json > test1.result.json 

echo "Generating result from test2.json file"
go run proj3/main $THREAD_COUNT < test2.json > test2.result.json 

echo "Generating result from test3.json file"
go run proj3/main $THREAD_COUNT < test3.json > test3.result.json

# Print usage
echo "Results generated, please using 'result_visulize.htm' in this folder to visulize result"