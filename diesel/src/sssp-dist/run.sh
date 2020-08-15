#!/bin/bash

# trap "exit" INT TERM ERR
# trap "kill 0" EXIT

go run main.go &

for i in {1..10}
do
    go run worker.go $i &
done
         
wait
