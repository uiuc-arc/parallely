import subprocess
import re
import numpy as np

ArrayDim = 100
ArraySz = ArrayDim*ArrayDim
SliceSz = ArraySz/10

template_str = open("sobel.tmpl", 'r').readlines()
with open("sobel.par", "w") as fout:
    for line in template_str:
        newline = line.replace('__ARRAYDIM__', str(ArrayDim))
        newline = newline.replace('__ARRAYSZ__', str(ArraySz))
        newline = newline.replace('__SLICESZ__', str(SliceSz))
        fout.write(newline)

template_str = open("__basic_go_main.tmpl", 'r').readlines()
with open("__basic_go_main.txt", "w") as fout:
    for line in template_str:
        newline = line.replace('__ARRAYDIM__', str(ArrayDim))
        newline = newline.replace('__ARRAYSZ__', str(ArraySz))
        newline = newline.replace('__SLICESZ__', str(SliceSz))
        fout.write(newline)

template_str = open("__basic_go_worker.tmpl", 'r').readlines()
with open("__basic_go_worker.txt", "w") as fout:
    for line in template_str:
        newline = line.replace('__ARRAYDIM__', str(ArrayDim))
        newline = newline.replace('__ARRAYSZ__', str(ArraySz))
        newline = newline.replace('__SLICESZ__', str(SliceSz))
        fout.write(newline)
