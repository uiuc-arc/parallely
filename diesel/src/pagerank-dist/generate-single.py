import subprocess
import re
import numpy as np

inputsets = {
    1: {"nodes": 8114, "slice": 1000, "file": "p2p-Gnutella09.txt"},
    2: {"nodes": 22687, "slice": 3000, "file": "p2p-Gnutella25.txt"},
    3: {"nodes": 36682, "slice": 5000, "file": "p2p-Gnutella30.txt"},
    4: {"nodes": 62586, "slice": 10000, "file": "p2p-Gnutella31.txt"},
    5: {"nodes": 1088092, "slice": 200000, "file": "roadNet-PA.txt"},
}

# src_size = 262144
# num_sample = 40

# data_set = {}

inputgraph = 1

max_degrees = 100
e_size = inputsets[inputgraph]["nodes"] * max_degrees
slice_size = inputsets[inputgraph]["slice"]

template_str = open("pagerank.template", 'r').readlines()
with open("__temp_gen.par", "w") as fout:
    for line in template_str:
        newline = line.replace('__NUMEDGES__', str(e_size))
        newline = newline.replace('__NUMNODES__', str(inputsets[inputgraph]["nodes"]))
        newline = newline.replace('__FILENAME__', str(inputsets[inputgraph]["file"]))
        newline = newline.replace('__SLICESIZE__', str(slice_size))
        newline = newline.replace('__MAX_DEGREE_', str(max_degrees))
        fout.write(newline)

template_str = open("__basic_go_main.tmpl", 'r').readlines()
with open("__temp_main.txt", "w") as fout:
    for line in template_str:
        newline = line.replace('__NUMEDGES__', str(e_size))
        newline = newline.replace('__NUMNODES__', str(inputsets[inputgraph]["nodes"]))
        newline = newline.replace('__FILENAME__', str(inputsets[inputgraph]["file"]))
        newline = newline.replace('__SLICESIZE__', str(slice_size))
        newline = newline.replace('__MAX_DEGREE_', str(max_degrees))        
        fout.write(newline)

template_str = open("__basic_go_worker.tmpl", 'r').readlines()
with open("__temp_worker.txt", "w") as fout:
    for line in template_str:
        newline = line.replace('__NUMEDGES__', str(e_size))
        newline = newline.replace('__NUMNODES__', str(inputsets[inputgraph]["nodes"]))
        newline = newline.replace('__FILENAME__', str(inputsets[inputgraph]["file"]))
        newline = newline.replace('__SLICESIZE__', str(slice_size))
        newline = newline.replace('__MAX_DEGREE_', str(max_degrees))        
        fout.write(newline)        
