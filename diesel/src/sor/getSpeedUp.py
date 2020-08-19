import subprocess
import re
import numpy as np
import time

def geo_mean(iterable):
    a = np.array(iterable)
    return a.prod()**(1.0 / len(a))

times = []
numsamples = 20

print "Running without dynamic tracking"
# Compile
commstr = """python ../../../parser/crosscompiler-diesel-acc.py -f sor.par -t __basic_go.txt -o out.go"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(numsamples):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("go run out.go", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000
    print time_spent
    times.append(time_spent)
    # time.sleep(2)

no_track_time = geo_mean(times)
print "Runtime without tracking: ", no_track_time

print "------------------------------------------"

# maybe this will remove the random crashes
# time.sleep(20)

# Compile
print "Running with dynamic tracking"
times = []

commstr = """python ../../../parser/crosscompiler-diesel-acc.py -f sor.par -t __basic_go.txt -o out.go -dyn"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(numsamples):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("go run out.go", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000
    print time_spent
    times.append(time_spent)
    # time.sleep(2)    

track_time = geo_mean(times)
print "Runtime with tracking: ", track_time

# maybe this will remove the random crashes
# time.sleep(20)

print "Running with array optimization"
times = []

commstr = """python ../../../parser/crosscompiler-diesel-acc.py -f sor.par -t __basic_go.txt -o out.go -dyn -a"""

result_test = subprocess.check_output(commstr, shell=True)
print result_test

for i in range(numsamples):
    print "Running Iteration : ", i
    result_test = subprocess.check_output("go run out.go", shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1]) / 1000000
    print time_spent
    times.append(time_spent)
    # time.sleep(2)   
    
opt_time = geo_mean(times)
print "Runtime with optimizations: ", opt_time

print "Overhead : ", ((track_time - no_track_time) / no_track_time) * 100
print "Overhead After Optimization : ", ((opt_time - no_track_time) / no_track_time) * 100
