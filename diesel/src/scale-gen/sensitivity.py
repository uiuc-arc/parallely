import subprocess
import re
import numpy as np

src_size = 262144
num_sample = 10

data_set = {}

for scale_factor in [2, 4, 8, 16]:
    dest_size = src_size * scale_factor * scale_factor
    slice_size = dest_size / 8

    template_str = open("scale.template", 'r').readlines()
    with open("__temp_gen.par", "w") as fout:
        for line in template_str:
            newline = line.replace('__DESTSIZE__', str(dest_size))
            newline = newline.replace('__SLICESIZE__', str(slice_size))
            newline = newline.replace('__SCALEFACTOR__', str(scale_factor))
            fout.write(newline)

    template_str = open("__basic_go.template", 'r').readlines()
    with open("__temp_gen.txt", "w") as fout:
        for line in template_str:
            newline = line.replace('__DESTSIZE__', str(dest_size))
            newline = newline.replace('__SLICESIZE__', str(slice_size))
            newline = newline.replace('__SCALEFACTOR__', str(scale_factor))
            fout.write(newline)

    commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o scale.go; go build -tags instrument"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    times = []
    memory = []
    for i in range(num_sample):
        print "Running Iteration : ", i
        result_test = subprocess.check_output("./scale-gen baboon.ppm temp.ppm", shell=True)
        print result_test
        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1])
        print time_spent
        times.append(time_spent)

        matches2 = re.findall("Memory through channels : .*\n", result_test)
        mem_used = float(matches2[0].split(' : ')[-1])
        print mem_used
        memory.append(mem_used)

    no_track_time = np.mean(times)
    no_track_memory = np.mean(memory)

    commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o scale.go -dyn; go build -tags instrument"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    times = []
    memory = []
    for i in range(num_sample):
        print "Running Iteration : ", i
        result_test = subprocess.check_output("./scale-gen baboon.ppm temp.ppm", shell=True)
        print result_test
        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1])
        print time_spent
        times.append(time_spent)

        matches2 = re.findall("Memory through channels : .*\n", result_test)
        mem_used = float(matches2[0].split(' : ')[-1])
        print mem_used
        memory.append(mem_used)

    track_time = np.mean(times)
    track_memory = np.mean(memory)

    commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o scale.go -dyn -a; go build -tags instrument"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    times = []
    memory = []
    for i in range(num_sample):
        print "Running Iteration : ", i
        result_test = subprocess.check_output("./scale-gen baboon.ppm temp.ppm", shell=True)
        print result_test
        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1])
        print time_spent
        times.append(time_spent)

        matches2 = re.findall("Memory through channels : .*\n", result_test)
        mem_used = float(matches2[0].split(' : ')[-1])
        print mem_used
        memory.append(mem_used)

    opt_track_time = np.mean(times)
    opt_track_memory = np.mean(memory)

    data_set[scale_factor] = (no_track_time, track_time, opt_track_time,
                              no_track_memory, track_memory, opt_track_memory)
    print data_set

print "*************************"
print data_set
print "*************************"
