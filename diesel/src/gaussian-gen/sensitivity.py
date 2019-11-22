import subprocess
import re
import numpy as np

src_size = 262144
num_sample = 10

# sizes = {128, 256, 512 1024}

data_set = {}
for scale_factor in [128, 256, 384, 512]:
    dest_size = scale_factor * scale_factor
    slice_size = dest_size / 8

    template_str = open("gaussian.template", 'r').readlines()
    with open("__temp_gen.par", "w") as fout:
        for line in template_str:
            newline = line.replace('__DESTSIZE__', str(dest_size))
            newline = newline.replace('__SLICESIZE__', str(slice_size))
            # newline = newline.replace('__GAUSSIANFACTOR__', str(gaussian_factor))
            fout.write(newline)

    template_str = open("__basic_go.templ", 'r').readlines()
    with open("__temp_gen.txt", "w") as fout:
        for line in template_str:
            newline = line.replace('__DESTSIZE__', str(dest_size))
            newline = newline.replace('__SLICESIZE__', str(slice_size))
            # newline = newline.replace('__GAUSSIANFACTOR__', str(gaussian_factor))
            fout.write(newline)

    commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o gaussian.go; go build -tags instrument"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    result_test = subprocess.check_output("./gaussian-gen temp-{}.ppm temp.ppm".format(scale_factor),
                                          shell=True)
    matches2 = re.findall("Memory through channels : .*\n", result_test)
    mem_used = float(matches2[0].split(' : ')[-1])
    no_track_memory = mem_used

    commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o gaussian.go; go build"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    times = []
    for i in range(num_sample):
        print "Running Iteration : ", i
        result_test = subprocess.check_output("./gaussian-gen temp-{}.ppm temp.ppm".format(scale_factor),
                                              shell=True)
        print result_test
        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1])
        print time_spent
        times.append(time_spent)

    no_track_time = np.mean(times)

    commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o gaussian.go -dyn; go build -tags instrument"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    result_test = subprocess.check_output("./gaussian-gen temp-{}.ppm temp.ppm".format(scale_factor),
                                          shell=True)
    matches2 = re.findall("Memory through channels : .*\n", result_test)
    mem_used = float(matches2[0].split(' : ')[-1])
    track_memory = mem_used

    commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o gaussian.go -dyn; go build"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    times = []
    for i in range(num_sample):
        print "Running Iteration : ", i
        result_test = subprocess.check_output("./gaussian-gen temp-{}.ppm temp.ppm".format(scale_factor),
                                              shell=True)
        print result_test
        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1])
        print time_spent
        times.append(time_spent)

        # matches2 = re.findall("Memory through channels : .*\n", result_test)
        # mem_used = float(matches2[0].split(' : ')[-1])
        # print mem_used

    track_time = np.mean(times)

    commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o gaussian.go -dyn -a; go build -tags instrument"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    result_test = subprocess.check_output("./gaussian-gen temp-{}.ppm temp.ppm".format(scale_factor),
                                          shell=True)
    matches2 = re.findall("Memory through channels : .*\n", result_test)
    mem_used = float(matches2[0].split(' : ')[-1])
    opt_track_memory = mem_used

    commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o gaussian.go -dyn -a; go build"""
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    times = []
    for i in range(num_sample):
        print "Running Iteration : ", i
        result_test = subprocess.check_output("./gaussian-gen temp-{}.ppm temp.ppm".format(scale_factor),
                                              shell=True)
        print result_test
        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1])
        print time_spent
        times.append(time_spent)

        # matches2 = re.findall("Memory through channels : .*\n", result_test)
        # mem_used = float(matches2[0].split(' : ')[-1])
        # print mem_used

    opt_track_time = np.mean(times)

    data_set[scale_factor] = (no_track_time, track_time, opt_track_time,
                              no_track_memory, track_memory, opt_track_memory)
    print data_set

print "*************************"
print data_set
print "*************************"
