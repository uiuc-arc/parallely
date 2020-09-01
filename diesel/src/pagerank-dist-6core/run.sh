#!/bin/bash

# trap "exit" INT TERM ERR
# trap "kill 0" EXIT

taskset --cpu-list 0 go run main.go &

for i in {1..5}
do
    taskset --cpu-list $i go run worker.go $i &
done
         
wait
