import subprocess
import time
import re
import numpy as np

query_str = "go run sssp-noisy.go ../inputs/p2p-Gnutella31.txt 62586"
query_str2 = "go run sssp-noisy-off.go ../inputs/p2p-Gnutella31.txt 62586"

approx_time = []
skiped = []
for i in range(10):
    temp = subprocess.check_output(query_str2, shell=True)
    matches = re.findall("Elapsed time : .*\n", temp)
    time_spent = float(matches[0].split(' : ')[-1])

    matches2 = re.findall("Retries : .*\n", temp)
    rets = float(matches2[0].split(' : ')[-1])
    approx_time.append(time_spent)
    skiped.append(rets)

redo_time = []
for i in range(10):
    temp = subprocess.check_output(query_str, shell=True)
    matches = re.findall("Elapsed time : .*\n", temp)
    time_spent = float(matches[0].split(' : ')[-1])
    redo_time.append(time_spent)

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

# print redo_time2
# print np.mean(redo_time2), np.mean(redo_time) / np.mean(redo_time2)