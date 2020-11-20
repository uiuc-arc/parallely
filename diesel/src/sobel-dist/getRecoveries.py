import subprocess
import re
import numpy as np
import time

recoveries = []
numsamples = 20

commstr = """python ../../../parser/crosscompiler-diesel-dist.py -f sobel.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o out.go -i -dyn -a -acc"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(numsamples):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)

    matches = re.findall("recoveries: .*\n", result_test)
    _recoveries = int(matches[0].split(':')[-1])
    print _recoveries
    recoveries.append(_recoveries)
    time.sleep(2)

print 'Total recoveries:', sum(recoveries)
print 'Recovery rate:', sum(recoveries)/10.0/numsamples
print 'Details:'
print recoveries
