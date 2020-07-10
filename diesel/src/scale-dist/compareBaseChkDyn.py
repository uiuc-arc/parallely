#! /usr/bin/env python2

import subprocess
import re
import sys
import numpy as np
import scipy.stats as st

num_sample = 100

benchName = "scale"
genCommand = "python ../../../parser/crosscompiler-diesel.py -f={}.par -t=__basic_go.txt -o {}.go".format(benchName, benchName)

print "--------"
print "Baseline"
print "--------"
print "Generating executable"

subprocess.check_output(genCommand, shell=True)
subprocess.check_output("go build", shell=True)

times = []
for i in range(num_sample):
    print "Running Iteration", i+1, "of", num_sample
    result_test = subprocess.check_output("./{}-gen temp-512.ppm /dev/null".format(benchName), shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1])
    print time_spent
    times.append(time_spent)

baseMean = st.gmean(times)
baseStdv = np.exp(np.std(np.log(times)))

print "Baseline Runtime Geomean:", baseMean
print "Baseline Runtime Geostdv:", baseStdv

print "--------"
print "Checksum"
print "--------"
print "Generating executable"

subprocess.check_output("sed -i 's/SendFloat64Array/SendChkFloat64Array/g' {}.go".format(benchName), shell=True)
subprocess.check_output("sed -i 's/ReceiveFloat64Array/ReceiveChkFloat64Array/g' {}.go".format(benchName), shell=True)
subprocess.check_output("go build", shell=True)

times = []
for i in range(num_sample):
    print "Running Iteration", i+1, "of", num_sample
    result_test = subprocess.check_output("./{}-gen temp-512.ppm /dev/null".format(benchName), shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1])
    print time_spent
    times.append(time_spent)

chkMean = st.gmean(times)
chkStdv = np.exp(np.std(np.log(times)))

print "Checksum Runtime Geomean:", chkMean
print "Checksum Runtime Geostdv:", chkStdv

print "--------"
print "DynTrack"
print "--------"
print "Generating executable"

subprocess.check_output(genCommand+" -dyn", shell=True)
subprocess.check_output("go build", shell=True)

times = []
for i in range(num_sample):
    print "Running Iteration", i+1, "of", num_sample
    result_test = subprocess.check_output("./{}-gen temp-512.ppm /dev/null".format(benchName), shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1])
    print time_spent
    times.append(time_spent)

dynMean = st.gmean(times)
dynStdv = np.exp(np.std(np.log(times)))

print "DynTrack Runtime Geomean:", dynMean
print "DynTrack Runtime Geostdv:", dynStdv


print "--------"
print "Summary"
print "--------"

print "Baseline", baseMean, baseStdv
print "Checksum", chkMean, chkStdv, 100.0*(chkMean-baseMean)/baseMean
print "DynTrack", dynMean, dynStdv, 100.0*(dynMean-baseMean)/baseMean
