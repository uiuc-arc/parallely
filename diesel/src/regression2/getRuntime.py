#! /usr/bin/env python2

import subprocess
import re
import sys
import numpy as np
import scipy.stats as st

num_sample = 100
times = []

WorkPerThread = sys.argv[1]
runType = sys.argv[2]
flags = None
if runType=='base':
    flags = ""
elif runType=='track':
    flags = "-dyn"
elif runType=='opt':
    flags = "-dyn -a"
else:
    raise Exception('bad run type')

result_test = subprocess.check_output("sed 's/WorkPerThread/{}/g' regression.tmpl > regression.par".format(WorkPerThread), shell=True)
# print result_test

result_test = subprocess.check_output("sed 's/WorkPerThread/{}/g' __basic_go.tmpl > __basic_go.txt".format(WorkPerThread), shell=True)
# print result_test

result_test = subprocess.check_output("python ../../../parser/crosscompiler-diesel.py -f=regression.par -t=__basic_go.txt -o regression.go {}".format(flags), shell=True)
# print result_test

result_test = subprocess.check_output("go build", shell=True)
# print result_test

for i in range(num_sample):
    print "Running Iteration", i+1, "of", num_sample
    result_test = subprocess.check_output("./regression2", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1])
    print time_spent
    times.append(time_spent)

times.sort()

print "Runtime Geomean:", st.gmean(times)
print "Runtime Std dev:", np.exp(np.std(np.log(times)))
