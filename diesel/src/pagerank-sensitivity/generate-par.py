import subprocess
import re
import numpy as np
import sys

inputsets = {
    1: {"nodes": 8114, "slice": 1000, "file": "p2p-Gnutella09.txt"},
    2: {"nodes": 22687, "slice": 3000, "file": "p2p-Gnutella25.txt"},
    3: {"nodes": 36682, "slice": 5000, "file": "p2p-Gnutella30.txt"},
    4: {"nodes": 62586, "slice": 10000, "file": "p2p-Gnutella31.txt"},
}

# src_size = 262144
# num_sample = 40

data_set = {}

inputgraph = sys.argv[1]
nodes = sys.argv[2]
slice_size = sys.argv[3]
e_size = int(nodes) * 10

print inputgraph, nodes, slice_size

template_str = open("pagerank.template", 'r').readlines()
with open("pagerank.par", "w") as fout:
    for line in template_str:
        newline = line.replace('__NUMEDGES__', str(e_size))
        newline = newline.replace('__NUMNODES__', str(nodes))
        newline = newline.replace('__FILENAME__', str(inputgraph))
        newline = newline.replace('__SLICESIZE__', str(slice_size))
        fout.write(newline)

template_str = open("__basic_go.template", 'r').readlines()
with open("__basic_go.txt", "w") as fout:
    for line in template_str:
        newline = line.replace('__NUMEDGES__', str(e_size))
        newline = newline.replace('__NUMNODES__', str(nodes))
        newline = newline.replace('__FILENAME__', str(inputgraph))
        newline = newline.replace('__SLICESIZE__', str(slice_size))
        fout.write(newline)
