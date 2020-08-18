import subprocess
import re
import numpy as np

def geo_mean(iterable):
    a = np.array(iterable)
    return a.prod()**(1.0 / len(a))

times = []
numsamples = 20

print("Running without dynamic tracking")

for i in range(numsamples):
    print("Running Iteration : ", i)
    result_test = subprocess.check_output("go run pagerank.go", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000
    print(time_spent)
    times.append(time_spent)

no_track_time = geo_mean(times)
print("Runtime without tracking: ", no_track_time)

print("------------------------------------------")

print("Running with optimization")
times = []

for i in range(numsamples):
    print("Running Iteration : ", i)
    result_test = subprocess.check_output("go run pagerank-opt.go", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000
    print(time_spent)
    times.append(time_spent)

opt_time = geo_mean(times)
print("Runtime with tracking: ", opt_time)

print("Overhead After Optimization : ", ((opt_time - no_track_time) / no_track_time) * 100)
