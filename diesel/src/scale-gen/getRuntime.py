import subprocess
import re
import numpy as np

num_sample = 30
times = []

result_test = subprocess.check_output("go build", shell=True)
print result_test

for i in range(num_sample):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./scale-gen baboon.ppm temp.ppm", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1])
    print time_spent
    times.append(time_spent)

opt_time = np.mean(times)

orig_time = 8614203708

print "Runtime: ", opt_time
print "SD: ", np.std(times)
# print "Overhead: ", opt_time / orig_time
