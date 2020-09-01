#!/bin/bash

# trap "exit" INT TERM ERR
# trap "kill 0" EXIT

go run -gcflags=-B main.go &

for i in {1..10}
do
    go run -gcflags=-B worker.go $i &
done
         
wait
