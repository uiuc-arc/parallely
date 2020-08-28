import subprocess
import re
import numpy as np
import time

def geo_mean(iterable):
    a = np.array(iterable)
    return a.prod()**(1.0 / len(a))

times = []
numsamples = 50

print "Running without dynamic tracking"
# Compile
commstr = """python ../../../parser/crosscompiler-diesel-dist-acc.py -f mm.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o out.go -i"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(numsamples):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000
    print time_spent
    times.append(time_spent)
    time.sleep(2)

no_track_time = geo_mean(times)
print "Runtime without tracking: ", no_track_time

print "------------------------------------------"

# Compile
# print "Running with old dynamic tracking"
# times = []

# commstr = """python ../../../parser/crosscompiler-diesel-dist-acc.py -f mm.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o out.go -i -dyn -a"""
# result_test = subprocess.check_output(commstr, shell=True)
# print result_test

# for i in range(numsamples):
#     print "Running Iteration : ", i
#     result_test = subprocess.check_output("./run.sh", shell=True)

#     matches = re.findall("Elapsed time : .*\n", result_test)
#     time_spent = float(matches[0].split(' : ')[-1]) / 1000000
#     print time_spent
#     times.append(time_spent)
#     time.sleep(2)

# Compile
old_opt_time = geo_mean(times)

print "Running with dynamic tracking"
times = []

commstr = """python ../../../parser/crosscompiler-diesel-dist-acc.py -f mm.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o out.go -i -dyn -a"""
result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(numsamples):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000
    print time_spent
    times.append(time_spent)
    time.sleep(2)

opt_time = geo_mean(times)
print "Runtime with optimizations: ", opt_time
print "Overhead After Optimization : ", ((opt_time - no_track_time) / no_track_time) * 100

print "------------------------------------------"

# Compile
print "Running with dynamic tracking"
times = []

# commstr = """python ../../../parser/crosscompiler-diesel-dist-acc.py -f mm.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o out.go -i -dyn -a"""
# result_test = subprocess.check_output(commstr, shell=True)
# print result_test

for i in range(numsamples):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run-dyn.sh", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000
    print time_spent
    times.append(time_spent)
    time.sleep(2)

opt_time = geo_mean(times)
print "Runtime with optimizations: ", opt_time
print "Old overheads : ", ((old_opt_time - no_track_time) / no_track_time) * 100
print "Overhead After Optimization : ", ((opt_time - no_track_time) / no_track_time) * 100