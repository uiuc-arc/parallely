import subprocess
import sys

BlocksPerWorker = int(sys.argv[1])

template_str = open("motion.tmpl", 'r').readlines()
with open("motion.par", "w") as fout:
    for line in template_str:
        newline = line.replace('__BLOCKSPERWORKER__', str(BlocksPerWorker))
        fout.write(newline)

template_str = open("__basic_go.tmpl", 'r').readlines()
with open("__basic_go.txt", "w") as fout:
    for line in template_str:
        newline = line.replace('__BLOCKSPERWORKER__', str(BlocksPerWorker))
        fout.write(newline)
