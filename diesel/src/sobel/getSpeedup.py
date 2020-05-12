import subprocess
import re
import numpy as np


def geo_mean(iterable):
    a = np.array(iterable)
    return a.prod()**(1.0 / len(a))

times = []

# print "Running without dynamic tracking"
# # Compile
# commstr = """python ../../../parser/crosscompiler-diesel.py -f sobel.par -t __basic_go.txt -o sobel.go"""

# result_test = subprocess.check_output(commstr, shell=True)
# print result_test

# for i in range(20):
#     print "Running Iteration : ", i
#     result_test = subprocess.check_output("go run sobel.go", shell=True)

#     matches = re.findall("Elapsed time : .*\n", result_test)
#     time_spent = float(matches[0].split(' : ')[-1])
#     print time_spent
#     times.append(time_spent)

# no_track_time = geo_mean(times)
# print "Runtime without tracking: ", no_track_time

# times = []

no_track_time = 1237606

# print "Running with dynamic tracking"
# # Compile
# commstr = """python ../../../parser/crosscompiler-diesel.py -f sobel.par -t __basic_go.txt -o sobel.go -dyn"""
# result_test = subprocess.check_output(commstr, shell=True)
# print result_test

# for i in range(10):
#     print "Running Iteration : ", i
#     result_test = subprocess.check_output("go run sobel.go", shell=True)

#     matches = re.findall("Elapsed time : .*\n", result_test)
#     time_spent = float(matches[0].split(' : ')[-1])
#     print time_spent
#     times.append(time_spent)

# track_time = np.mean(times)
# print "Runtime with tracking (without opt): ", track_time

# print "Speedup (without opt) : ", track_time/no_track_time

times = []

print "Running with dynamic tracking"
# Compile
# commstr = """python ../../../parser/crosscompiler-diesel.py -f sobel.par -t __basic_go.txt -o sobel-opt.go -dyn -a"""
# result_test = subprocess.check_output(commstr, shell=True)
# print result_test

for i in range(20):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("go run sobel-opt.go", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1])
    print time_spent
    times.append(time_spent)

track_time = geo_mean(times)
print "Runtime with tracking: ", track_time

print "Speedup : ", track_time/no_track_time

print "Overhead : ", ((track_time - no_track_time) / no_track_time) * 100
# print "Overhead After Optimization : ", ((opt_time - no_track_time) / no_track_time) * 100
