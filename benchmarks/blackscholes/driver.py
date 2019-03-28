import subprocess
import time
import re
import numpy as np

query_str = "go run blackscholes-corrupt.go"
query_str2 = "go run blackscholes-corrupt-redo.go exact.out"
query_str3 = "go run blackscholes-corrupt-redo-off.go errors.out"

# approx_time = []
# for i in range(10):
#     temp = subprocess.check_output(query_str, shell=True)
#     matches = re.findall("Elapsed time : .*\n", temp)
#     time_spent = float(matches[0].split(' : ')[-1])
#     approx_time.append(time_spent)

redo_time = []
print "Running: ", query_str2
for i in range(30):
    print "."
    temp = subprocess.check_output(query_str2, shell=True)
    matches = re.findall("Elapsed time : .*\n", temp)
    time_spent = float(matches[0].split(' : ')[-1])
    redo_time.append(time_spent)

exact_results = []
result_file = open("exact.out", 'r').readlines()
for line in result_file:
    exact_results.append(float(line))

redo_time2 = []
approx_errors = []
print "Running: ", query_str3
for i in range(30):
    print "."
    temp = subprocess.check_output(query_str3, shell=True)
    matches = re.findall("Elapsed time : .*\n", temp)
    time_spent = float(matches[0].split(' : ')[-1])
    redo_time2.append(time_spent)

    approx_results = []
    approx_result_file = open("errors.out", 'r').readlines()
    for line in approx_result_file:
        approx_results.append(float(line))
    diffs = 0
    for i, element in enumerate(approx_results):
        if exact_results[i] == 0:
            continue
        diffs += abs(approx_results[i] - exact_results[i]) / exact_results[i]
    print diffs
    approx_errors.append(diffs / float(len(exact_results)))

# print approx_time
# print redo_time
# print np.mean(approx_time), np.mean(redo_time), np.mean(redo_time) / np.mean(approx_time)
print approx_errors, np.mean(approx_errors)

print redo_time2
print np.mean(redo_time2), np.mean(redo_time) / np.mean(redo_time2)
