import subprocess
import re
import numpy as np

num_sample = 10

times = []

print "Running without dynamic tracking"
# Compile
commstr = """python ../../../parser/crosscompiler-diesel.py -f scale.par -t __basic_go.txt -o scale.go; go build"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(num_sample):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./scale-gen baboon.ppm temp.ppm", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1])
    print time_spent
    times.append(time_spent)

no_track_time = np.mean(times)
print "Runtime without tracking: ", no_track_time

print "------------------------------------------"

# Compile
print "Running with dynamic tracking"
times = []

commstr = """python ../../../parser/crosscompiler-diesel.py -f scale.par -t __basic_go.txt -o scale.go -dyn; go build"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(num_sample):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./scale-gen baboon.ppm temp.ppm", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1])
    print time_spent
    times.append(time_spent)

track_time = np.mean(times)
print "Runtime with tracking: ", track_time

print "Running with array optimization"
times = []

commstr = """python ../../../parser/crosscompiler-diesel.py -f scale.par -t __basic_go.txt -o scale.go -dyn -a; go build"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(num_sample):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./scale-gen baboon.ppm temp.ppm", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1])
    print time_spent
    times.append(time_spent)

opt_time = np.mean(times)
print "Runtime with tracking: ", track_time

print "Overhead : ", track_time / no_track_time
print "Overhead After Optimization : ", opt_time / no_track_time
