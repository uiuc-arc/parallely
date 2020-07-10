import subprocess
import re
import numpy as np
import scipy.stats as st

times = []

print "Running without dynamic tracking"
# Compile
commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f bfs.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o bfs.go -i"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(20):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)
    print "Finished running"

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1])
    print time_spent
    times.append(time_spent)

no_track_time =st.gmean(times) # np.mean(times)
print "Runtime without tracking: ", no_track_time

# no_track_time = 4443605867.5

print "------------------------------------------"

# Compile
print "Running with dynamic tracking"
times = []

commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f bfs.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o bfs.go -dyn -i"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(20):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1])
    print time_spent
    times.append(time_spent)

track_time = st.gmean(times) #np.mean(times)
print "Runtime with tracking: ", track_time

print "Running with array optimization"
times = []

commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f bfs.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o bfs.go -dyn -a -i"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(20):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1])
    print time_spent
    times.append(time_spent)

opt_time = st.gmean(times) # np.mean(times)
print "Runtime with opt: ", opt_time

print "Overhead : ", ((track_time - no_track_time) / no_track_time)
print "Overhead (Opt) : ", ((opt_time - no_track_time) / no_track_time)
