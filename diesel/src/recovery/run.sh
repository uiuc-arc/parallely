#!/bin/bash

# trap "exit" INT TERM ERR
# trap "kill 0" EXIT

go run main.go &
go run worker.go 1 &
         
wait
