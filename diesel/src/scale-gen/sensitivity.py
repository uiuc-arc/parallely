import subprocess
import re
import numpy as np

src_size = 262144
num_sample = 5

data_set = {}


def parse_results_memory(resstring):
    matches2 = re.findall("Memory through channels : .*\n", resstring)
    return float(matches2[0].split(' : ')[-1])


def getRuntimes():
    times = []
    for i in range(num_sample):
        print "Running Iteration : ", i
        result_test = subprocess.check_output("./scale-gen baboon.ppm temp-{}.ppm".format(scale_factor),
                                              shell=True)
        print result_test
        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1])
        print time_spent
        times.append(time_spent)
    return np.mean(times)


for scale_factor in [2, 4, 8, 16]:
    dest_size = src_size * scale_factor * scale_factor
    slice_size = dest_size / 8

    # Creating the input specific files from the template
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

    # ----------------
    # Original Version
    # ----------------
    buildstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o scale.go; go build -tags instrument"""
    result_test = subprocess.check_output(buildstr, shell=True)
    print result_test
    result_test = subprocess.check_output("./scale-gen baboon.ppm temp-{}.ppm".format(scale_factor),
                                          shell=True)
    print result_test
    no_track_memory = parse_results_memory(result_test)
    commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o scale.go; go build"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test
    no_track_time = getRuntimes()

    # ----------------
    # Instrumented Version
    # ----------------
    commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o scale.go -dyn; go build -tags instrument"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    result_test = subprocess.check_output("./scale-gen baboon.ppm temp-{}.ppm".format(scale_factor),
                                          shell=True)
    print result_test
    matches2 = re.findall("Memory through channels : .*\n", result_test)
    mem_used = float(matches2[0].split(' : ')[-1])
    track_memory = mem_used

    commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o scale.go; go build"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    track_time = getRuntimes()

    # ----------------
    # Optimized Version
    # ----------------
    commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o scale.go -dyn -a; go build -tags instrument"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test
    result_test = subprocess.check_output("./scale-gen baboon.ppm temp-{}.ppm".format(scale_factor),
                                          shell=True)
    print result_test
    matches2 = re.findall("Memory through channels : .*\n", result_test)
    mem_used = float(matches2[0].split(' : ')[-1])
    opt_track_memory = mem_used

    commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o scale.go; go build"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    opt_track_time = getRuntimes()

    data_set[scale_factor] = (no_track_time, track_time, opt_track_time,
                              no_track_memory, track_memory, opt_track_memory)
    print data_set

print "*************************"
print data_set
print "*************************"
