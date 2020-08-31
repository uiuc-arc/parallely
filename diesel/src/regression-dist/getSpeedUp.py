import subprocess
import re
import numpy as np
import scipy.stats as st
import pickle
from os import path
import os

checkpointing = {"baseline": [0, []],
                 "dyn": [0, []],
                 "opt": [0, []]}

numsamples = 100

if path.exists("checkpoint.pickle"):
    with open("checkpoint.pickle", 'rb') as f:
        checkpointing = pickle.load(f)

if checkpointing["baseline"][0]<numsamples:
    remaining = numsamples - checkpointing["baseline"][0]
    print("Running without dynamic tracking")
    # Compile
    commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f regression2.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o _ -i"""

    result_test = subprocess.check_output(commstr, shell=True)
    print(result_test)

    for i in range(remaining):
        print("Running Iteration : ", checkpointing["baseline"][0])
        result_test = subprocess.check_output("./run.sh", shell=True)

        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1]) / 1000000
        print(time_spent)
        checkpointing["baseline"][0] += 1
        checkpointing["baseline"][1].append(time_spent)
        with open("checkpoint.pickle", 'wb') as f:
            pickle.dump(checkpointing, f)
        
no_track_time =st.gmean(checkpointing["baseline"][1])
print("Runtime without tracking: ", no_track_time)
print("------------------------------------------")

# Compile
print("Running with dynamic tracking")

if checkpointing["dyn"][0]<numsamples:
    remaining = numsamples - checkpointing["opt"][0]
    commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f regression2.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o _ -i -dyn -rel"""

    result_test = subprocess.check_output(commstr, shell=True)
    print(result_test)

    for i in range(remaining):
        print("Running Iteration : ", checkpointing["dyn"][0])
        result_test = subprocess.check_output("./run.sh", shell=True)

        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1]) / 1000000
        print(time_spent)
        checkpointing["dyn"][0] += 1
        checkpointing["dyn"][1].append(time_spent)
        with open("checkpoint.pickle", 'wb') as f:
            pickle.dump(checkpointing, f)

track_time = st.gmean(checkpointing["dyn"][1]) #np.mean(times)
print("Runtime with tracking: ", track_time)

print("Running with optimization")
# commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f regression.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -i -o pg.go -dyn -a"""

# result_test = subprocess.check_output(commstr, shell=True)
# print result_test

if checkpointing["opt"][0]<numsamples:
    remaining = numsamples - checkpointing["opt"][0]
    for i in range(remaining):
        print("Running Iteration : ", i)
        result_test = subprocess.check_output("./run-opt.sh", shell=True)

        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1]) / 1000000
        print(time_spent)
        checkpointing["opt"][0] += 1
        checkpointing["opt"][1].append(time_spent)
        with open("checkpoint.pickle", 'wb') as f:
            pickle.dump(checkpointing, f)        

opt_time = st.gmean(checkpointing["opt"][1]) # np.mean(times)
print("optimized tracking: ", track_time)

print("Overhead : ", ((track_time - no_track_time) / no_track_time) * 100)
print("Overhead After Optimization : ", ((opt_time - no_track_time) / no_track_time) * 100)
os.remove("checkpoint.pickle") 
