#!/bin/bash

# trap "exit" INT TERM ERR
# trap "kill 0" EXIT

go run main-dyn.go &

for i in {1..10}
do
    go run worker-dyn.go $i &
done
         
wait
