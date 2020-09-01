import subprocess
import re
import numpy as np
import scipy.stats as st

src_size = 262144
num_sample = 40

# sizes = {128, 256, 512 1024}

data_set = {}
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
    for i in range(num_sample):
        print "Running Iteration : ", i
        result_test = subprocess.check_output("./run.sh", shell=True)
        print result_test
        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1])
        print time_spent
        times.append(time_spent)

    no_track_time = st.gmean(times)

    commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f __temp_gen.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o gaussian -i -dyn -rel"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    times = []
    for i in range(num_sample):
        print "Running Iteration : ", i
        result_test = subprocess.check_output("./run.sh", shell=True)
        print result_test
        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1])
        print time_spent
        times.append(time_spent)

    track_time = st.gmean(times)

    commstr = """python2 ../../../parser/crosscompiler-diesel-dist.py -f __temp_gen.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o gaussian -i -dyn -a -rel"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    times = []
    for i in range(num_sample):
        print "Running Iteration : ", i
        result_test = subprocess.check_output("./run.sh", shell=True)
        print result_test
        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1])
        print time_spent
        times.append(time_spent)

    opt_track_time = st.gmean(times)

    data_set[scale_factor] = (no_track_time, track_time, opt_track_time)
    print data_set

print "*************************"
print data_set
print "*************************"
