#!/bin/bash 
export GOPATH="/home/vimuth/projects/diesel/go-lib/:/home/vimuth/go/"

go build
./pagerank ../inputs/p2p-Gnutella31.txt 62586 4 8 0 temp.out -cpuprofile cpu.prof -memprofile mem.prof
