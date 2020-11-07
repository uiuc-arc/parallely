#!/bin/bash

# trap "exit" INT TERM ERR
# trap "kill 0" EXIT

for exp in {0..100}
do 
    go run main-error.go $exp &

    for i in {1..8}
    do
        go run worker-error.go $i &
    done
         
    wait
done
