#! /usr/bin/env python2

import subprocess
import re
import sys
import numpy as np

num_sample = 30
times = []

result_test = subprocess.check_output("go build", shell=True)
print result_test

for i in range(num_sample):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./regression {}".format(sys.argv[1]), shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1])
    print time_spent
    times.append(time_spent)

print "Runtime: ", np.mean(times)
print "SD: ", np.std(times)
