import subprocess
import re
import numpy as np
import scipy.stats as st
import pickle
from os import path
import os

inputsizes = [2500, 3000, 3500, 4000]
numsamples = 40

result = {}

for inputsize in inputsizes:
    # Optimized version
    template_str = open("worker-opt.tmpl", 'r').readlines()
    with open("worker-opt-gen.go", "w") as fout:
        for line in template_str:
            newline = line.replace('__DATASIZE__', str(inputsize))
            fout.write(newline)            
    template_str = open("main-opt.tmpl", 'r').readlines()
    with open("main-opt-gen.go", "w") as fout:
        for line in template_str:
            newline = line.replace('__DATASIZE__', str(inputsize))
            fout.write(newline)

    template_str = open("__basic_go_worker.tmpl", 'r').readlines()
    with open("__basic_go_worker.txt", "w") as fout:
        for line in template_str:
            newline = line.replace('__DATASIZE__', str(inputsize))
            fout.write(newline)            
    template_str = open("__basic_go_main.tmpl", 'r').readlines()
    with open("__basic_go_main.txt", "w") as fout:
        for line in template_str:
            newline = line.replace('__DATASIZE__', str(inputsize))
            fout.write(newline)
    template_str = open("regression.tmpl", 'r').readlines()
    with open("regression.par", "w") as fout:
        for line in template_str:
            newline = line.replace('__DATASIZE__', str(inputsize))
            fout.write(newline)

    commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f regression.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o _ -i"""
    result_test = subprocess.check_output(commstr, shell=True)
    print(result_test)

    times = []
    for i in range(numsamples):
        print("Running Iteration : ", i)
        result_test = subprocess.check_output("./run.sh", shell=True)

        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1]) / 1000000
        print(time_spent)

        times.append(time_spent)
        baseline = st.gmean(times)

    commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f regression.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o _ -i -dyn -rel"""
    result_test = subprocess.check_output(commstr, shell=True)
    print(result_test)

    times = []
    for i in range(numsamples):
        print("Running Iteration : ", i)
        result_test = subprocess.check_output("./run.sh", shell=True)

        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1]) / 1000000
        print(time_spent)

        times.append(time_spent)
        tracking_time = st.gmean(times)

    for i in range(numsamples):
        print("Running Iteration : ", i)
        result_test = subprocess.check_output("./run-opt-gen.sh", shell=True)

        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1]) / 1000000
        print(time_spent)
        times.append(time_spent)
        opt_time = st.gmean(times)

    overhead = ((tracking_time - baseline) / baseline) * 100
    overhead_opt = ((opt_time - baseline) / baseline) * 100

    result[inputsize] = (overhead, overhead_opt)
    print result
