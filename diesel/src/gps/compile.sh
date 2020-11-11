#! /usr/bin/env bash

python ../../../parser/crosscompiler-diesel-dist.py -f gps.par -tm boilerplate_main.tmpl -tw boilerplate_worker.tmpl -i -dyn -acc
