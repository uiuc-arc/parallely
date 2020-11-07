import subprocess
import re
import numpy as np
import scipy.stats as st
from os import path
import os
import pickle

src_size = 262144
numsamples = 20

# sizes = {128, 256, 512 1024}

checkpointing = {256: [[], [], []],
                 384: [[], [], []],
                 512: [[], [], []],
                 1024: [[], [], []]}

if path.exists("checkpoint.pickle"):
    with open("checkpoint.pickle", 'rb') as f:
        checkpointing = pickle.load(f)

for scale_factor in [256, 384, 512, 1024]:
    dest_size = scale_factor * scale_factor
    slice_size = dest_size / 8

    template_str = open("gaussian.template", 'r').readlines()
    with open("__temp_gen.par", "w") as fout:
        for line in template_str:
            newline = line.replace('__DESTSIZE__', str(dest_size))
            newline = newline.replace('__SLICESIZE__', str(slice_size))
            # newline = newline.replace('__GAUSSIANFACTOR__', str(gaussian_factor))
            fout.write(newline)

    template_str = open("__basic_go_main.tmpl", 'r').readlines()
    with open("__basic_go_main.txt", "w") as fout:
        for line in template_str:
            newline = line.replace('__DESTSIZE__', str(dest_size))
            newline = newline.replace('__SLICESIZE__', str(slice_size))
            newline = newline.replace('__FILENAME__', str(scale_factor))
            # newline = newline.replace('__GAUSSIANFACTOR__', str(gaussian_factor))
            fout.write(newline)

    template_str = open("__basic_go_worker.tmpl", 'r').readlines()
    with open("__basic_go_worker.txt", "w") as fout:
        for line in template_str:
            newline = line.replace('__DESTSIZE__', str(dest_size))
            newline = newline.replace('__SLICESIZE__', str(slice_size))
            newline = newline.replace('__FILENAME__', str(scale_factor))
            # newline = newline.replace('__GAUSSIANFACTOR__', str(gaussian_factor))
            fout.write(newline)

    commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f __temp_gen.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o gaussian -i -rel"""
    result_test = subprocess.check_output(commstr, shell=True)
    print(result_test)

    times = []

    if len(checkpointing[scale_factor][0])<numsamples:
        remaining = numsamples - len(checkpointing[scale_factor][0])
        print "+++++++++++++ Remaining samples: ", remaining
        for i in range(remaining):
            print "Running Iteration : ", i
            result_test = subprocess.check_output("./run.sh", shell=True)
            print result_test
            matches = re.findall("Elapsed time : .*\n", result_test)
            time_spent = float(matches[0].split(' : ')[-1])
            print time_spent
            # times.append(time_spent)
            checkpointing[scale_factor][0].append(time_spent)
            with open("checkpoint.pickle", 'wb') as f:
                pickle.dump(checkpointing, f)

    # no_track_time = st.gmean(times)

    commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f __temp_gen.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o gaussian -i -dyn -rel"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    times = []
    if len(checkpointing[scale_factor][1])<numsamples:
        remaining =  numsamples - len(checkpointing[scale_factor][1])
        for i in range(remaining):
            print "Running Iteration : ", i
            result_test = subprocess.check_output("./run.sh", shell=True)
            print result_test
            matches = re.findall("Elapsed time : .*\n", result_test)
            time_spent = float(matches[0].split(' : ')[-1])
            print time_spent
            times.append(time_spent)
            checkpointing[scale_factor][1].append(time_spent)
            with open("checkpoint.pickle", 'wb') as f:
                pickle.dump(checkpointing, f)


    commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f __temp_gen.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o gaussian -i -dyn -a -rel"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    times = []
    if len(checkpointing[scale_factor][2])<numsamples:
        remaining =  numsamples - len(checkpointing[scale_factor][2])
        for i in range(remaining):
            print "Running Iteration : ", i
            result_test = subprocess.check_output("./run.sh", shell=True)
            print result_test
            matches = re.findall("Elapsed time : .*\n", result_test)
            time_spent = float(matches[0].split(' : ')[-1])
            print time_spent
            times.append(time_spent)
            checkpointing[scale_factor][2].append(time_spent)
            with open("checkpoint.pickle", 'wb') as f:
                pickle.dump(checkpointing, f)

    # opt_track_time = st.gmean(times)

    # data_set[scale_factor] = (no_track_time, track_time, opt_track_time)
#     print data_set

# print "*************************"
# print data_set
# print "*************************"
