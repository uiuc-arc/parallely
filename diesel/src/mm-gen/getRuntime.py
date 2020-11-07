import subprocess
import re
import numpy as np

def geo_mean(iterable):
    a = np.array(iterable)
    return a.prod()**(1.0 / len(a))

times = []
numsamples = 50

for i in range(numsamples):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("go run mm.go", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 100000
    print time_spent
    times.append(time_spent)

opt_time = geo_mean(times)
print "Runtime: ", opt_time
