import subprocess
import re
import numpy as np
import matplotlib.pyplot as plt

def getRelDegredation():
    min_rels = []
    for i in range(10):
        dyn_file = open("dynmap" + str(i), "r")
        dynvalues = [float(k) for k in dyn_file.readlines()[0][1:-2].split(' ')]
        min_rels.append(np.min(dynvalues))
    print min_rels
    return min_rels

print "Running with dynamic tracking"
times = []

commstr = """rm -f dynmap*; python ../../../parser/crosscompiler-diesel.py -f sssp.par -t __basic_go.txt -o sssp.go -dyn -i; go build"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

result_test = subprocess.check_output("./sssp-gen", shell=True)
print result_test
before_opt = getRelDegredation()

print "Running with array optimization"

commstr = """rm -f dynmap*; python ../../../parser/crosscompiler-diesel.py -f sssp.par -t __basic_go.txt -o sssp.go -dyn -i -a; go build"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

result_test = subprocess.check_output("./sssp-gen", shell=True)
print result_test
after_opt = getRelDegredation()

plt.plot(before_opt)
plt.plot(after_opt)

plt.ylabel("Tracked Reliability")
plt.xlabel("Iteration")

plt.savefig("reliability.png")
