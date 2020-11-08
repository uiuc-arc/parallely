import subprocess
import re
import numpy as np

def geo_mean(iterable):
    a = np.array(iterable)
    return a.prod()**(1.0 / len(a))

times = []
numsamples = 50

print "Running without dynamic tracking"
# Compile
commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f combined-dist.par -tm template.txt -tw template_worker.txt -o _ -i"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(numsamples):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000.0
    print time_spent
    times.append(time_spent)

no_track_time = geo_mean(times)
print "Runtime without tracking: ", no_track_time

print "------------------------------------------"

# Compile
print "Running with dynamic tracking"
times = []

commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f combined-dist.par -tm template.txt -tw template_worker.txt -o _ -dyn -i"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(numsamples):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000.0
    print time_spent
    times.append(time_spent)

track_time = geo_mean(times)
print "Runtime with tracking: ", track_time

# print "Running with array optimization"
# times = []

# commstr = """python ../../../parser/crosscompiler-diesel.py -f pagerank.par -t __basic_go.txt -o pagerank.go -dyn -a; go build"""

# result_test = subprocess.check_output(commstr, shell=True)
# print result_test

# for i in range(10):
#     print "Running Iteration : ", i
#     result_test = subprocess.check_output("./pagerank-gen", shell=True)

#     matches = re.findall("Elapsed time : .*\n", result_test)
#     time_spent = float(matches[0].split(' : ')[-1])
#     print time_spent
#     times.append(time_spent)

# opt_time = np.mean(times)
# print "Runtime with tracking: ", track_time

print "Overhead : ", ((track_time - no_track_time) / no_track_time) * 100
# print "Overhead After Optimization : ", ((opt_time - no_track_time) / no_track_time) * 100