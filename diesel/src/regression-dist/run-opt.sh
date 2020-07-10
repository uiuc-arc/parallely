#!/bin/bash

# trap "exit" INT TERM ERR
# trap "kill 0" EXIT

go run main-opt.go &

for i in {1..10}
do
    go run worker-opt.go $i &
done
         
wait
