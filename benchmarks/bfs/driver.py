import subprocess
import time
import re
import numpy as np

query_str = "go run bfs-noisy.go ../inputs/p2p-Gnutella31.txt 62586 out-exact.txt"
query_str2 = "go run bfs-noisy-off.go ../inputs/p2p-Gnutella31.txt 62586 out-approx.txt"

approx_time = []
skiped = []
print "running; ", query_str2
for i in range(10):
    print "."
    temp = subprocess.check_output(query_str, shell=True)
    matches = re.findall("Elapsed time : .*\n", temp)
    time_spent = float(matches[0].split(' : ')[-1])

    matches2 = re.findall("Retries : .*\n", temp)
    rets = float(matches2[0].split(' : ')[-1])
    approx_time.append(time_spent)
    skiped.append(rets)

exact_results = []
result_file = open("out-exact.txt", 'r').readlines()
for line in result_file:
    exact_results.append(float(line))

approx_errors = []
redo_time = []
print "running; ", query_str
for i in range(10):
    print "."
    temp = subprocess.check_output(query_str2, shell=True)
    matches = re.findall("Elapsed time : .*\n", temp)
    time_spent = float(matches[0].split(' : ')[-1])
    redo_time.append(time_spent)

    approx_results = []
    approx_result_file = open("out-approx.txt", 'r').readlines()
    for line in approx_result_file:
        approx_results.append(float(line))
    diffs = 0
    for i, element in enumerate(approx_results):
        # different maxes?
        if approx_results[i] - exact_results[i] != 0:
            diffs += 1
    approx_errors.append(diffs / float(len(exact_results)))
    approx_time.append(time_spent)

# redo_time2 = []
# for i in range(10):
#     temp = subprocess.check_output(query_str3, shell=True)
#     matches = re.findall("Elapsed time : .*\n", temp)
#     time_spent = float(matches[0].split(' : ')[-1])
#     redo_time2.append(time_spent)

print approx_time
print redo_time
print np.mean(approx_time), np.mean(redo_time), np.mean(redo_time) / np.mean(approx_time)
print skiped, np.mean(skiped)
print approx_errors, np.mean(approx_errors)
# print redo_time2
# print np.mean(redo_time2), np.mean(redo_time) / np.mean(redo_time2)
