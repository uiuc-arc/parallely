import subprocess
import time
import re
import numpy as np

query_str = "go run scale-taskfail.go baboon.ppm scaled-baboon.ppm 4 64 1"
query_str2 = "go run scale-taskfail.go baboon.ppm scaled-baboon-errors.ppm 4 64 0"

psnr_str = "compare -metric psnr scaled-baboon.ppm scaled-baboon-errors.ppm _diff.jpp"

approx_time = []
# skiped = []
print "running; ", query_str2
for i in range(10):
    print "."
    temp = subprocess.check_output(query_str2, shell=True)
    matches = re.findall("Elapsed time : .*\n", temp)
    time_spent = float(matches[0].split(' : ')[-1])

    # matches2 = re.findall("Retries : .*\n", temp)
    # rets = float(matches2[0].split(' : ')[-1])
    approx_time.append(time_spent)
    # skiped.append(rets)

redo_time = []
errors = []
print "running; ", query_str
for i in range(10):
    print "."
    temp = subprocess.check_output(query_str, shell=True)
    matches = re.findall("Elapsed time : .*\n", temp)
    time_spent = float(matches[0].split(' : ')[-1])
    redo_time.append(time_spent)

    cmd = subprocess.Popen(psnr_str, shell=True,
                           stderr=subprocess.PIPE, stdout=subprocess.PIPE)
    cmd_out, cmd_err = cmd.communicate()
    output = float(cmd_err)
    print output
    errors.append(output)


print approx_time
print redo_time
print errors
print np.mean(approx_time), np.mean(redo_time), np.mean(redo_time) / np.mean(approx_time)
# print skiped, np.mean(skiped)

# print redo_time2
# print np.mean(redo_time2), np.mean(redo_time) / np.mean(redo_time2)
