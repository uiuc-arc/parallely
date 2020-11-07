import subprocess
import re
import numpy as np
import scipy.stats as st

times = []
numsamples = 42

print "Running without dynamic tracking"
# Compile
commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f pagerank.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -i -o pg.go -rel"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(numsamples):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000
    print time_spent
    times.append(time_spent)

no_track_time =st.gmean(times) # np.mean(times)
print "Runtime without tracking: ", no_track_time

print "------------------------------------------"

# Compile
print "Running with dynamic tracking"
times = []

commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f pagerank.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -i -o pg.go -dyn -rel -n"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(numsamples):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000
    print time_spent
    times.append(time_spent)

track_time = st.gmean(times) #np.mean(times)
print "Runtime with tracking: ", track_time

print "Running with array optimization"
times = []

commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f pagerank.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -i -o pg.go -dyn -a -rel -n"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(numsamples):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000
    print time_spent
    times.append(time_spent)

opt_time = st.gmean(times) # np.mean(times)
print "Runtime with tracking: ", track_time

print "Overhead : ", ((track_time - no_track_time) / no_track_time) * 100
print "Overhead After Optimization : ", ((opt_time - no_track_time) / no_track_time) * 100
