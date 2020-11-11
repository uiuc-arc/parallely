#!/usr/bin/env bash

go build -gcflags=-B main.go
go build -gcflags=-B worker.go 

./main &
./worker &

wait
