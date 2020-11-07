#!/bin/bash

# trap "exit" INT TERM ERR
# trap "kill 0" EXIT

go build -gcflags=-B main.go
go build -gcflags=-B worker.go 

./main &

for i in {1..10}
do
    ./worker $i &
done
         
wait
