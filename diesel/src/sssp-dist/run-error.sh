#!/bin/bash

# trap "exit" INT TERM ERR
# trap "kill 0" EXIT

go run main-error.go &

for i in {1..10}
do
    go run worker-error.go $i &
done
         
wait
