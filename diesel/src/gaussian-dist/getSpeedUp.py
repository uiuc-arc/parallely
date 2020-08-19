import subprocess
import re
import numpy as np

def geo_mean(iterable):
    a = np.array(iterable)
    return a.prod()**(1.0 / len(a))

num_sample = 20

times = []

print "Running without dynamic tracking"
# Compile
commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f gaussian.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o gaussian -i"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(num_sample):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)
    # print result_test

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000
    print time_spent
    times.append(time_spent)

no_track_time = geo_mean(times)
print "Runtime without tracking: ", no_track_time

print "------------------------------------------"

# Compile
print "Running with dynamic tracking"
times = []

commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f gaussian.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o gaussian -i -dyn"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(num_sample):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)
    # print result_test

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000
    print time_spent
    times.append(time_spent)

track_time = geo_mean(times)
print "Runtime with tracking: ", track_time

print "Running with array optimization"
times = []

commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f gaussian.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o gaussian -i -dyn -a"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(num_sample):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("./run.sh", shell=True)
    # print result_test

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000
    print time_spent
    times.append(time_spent)

opt_time = geo_mean(times)
print "Runtime after optimizations: ", opt_time

print "---------------------------------"

print "Overhead : ", ((track_time - no_track_time) / no_track_time) * 100
print "Overhead (Opt) : ", ((opt_time - no_track_time) / no_track_time) * 100
print "---------------------------------"

