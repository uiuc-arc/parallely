import subprocess
import re
import numpy as np
import time

numsamples = 100
recoveries = 0

for i in range(numsamples):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)
    matches = re.findall("Error : .*\n", result_test)
    error = float(matches[0].split(' : ')[-1])
    print error
    if error > 1e-7:
        recoveries += 1

print recoveries, float(recoveries)/numsamples
