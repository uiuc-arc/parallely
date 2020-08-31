import subprocess
import re
import numpy as np

def geo_mean(iterable):
    a = np.array(iterable)
    return a.prod()**(1.0 / len(a))

results = {}
times = []
numsamples = 1

for numsensors in [248, 512, 1024, 2048]:
    template_str = open("kmeans.template", 'r').readlines()
    with open("kmeansgen.par", "w") as fout:
        for line in template_str:
            newline = line.replace('__NUMSENSORS__', str(numsensors))
            fout.write(newline)

    template_str = open("kmean_go_main.tmpl", 'r').readlines()
    with open("__basic_go_main.txt", "w") as fout:
        for line in template_str:
            newline = line.replace('__NUMSENSORS__', str(numsensors))
            fout.write(newline)

    template_str = open("kmean_go_worker.tmpl", 'r').readlines()
    with open("__basic_go_worker.txt", "w") as fout:
        for line in template_str:
            newline = line.replace('__NUMSENSORS__', str(numsensors))
            fout.write(newline)    

    print "Running without dynamic tracking"
    # Compile
    commstr = "python2 ../../../parser/crosscompiler-diesel-dist.py -f kmeansgen.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o kmeans -i"

    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    for i in range(numsamples):
        print "Running Iteration : ", i
        result_test = subprocess.check_output("./run.sh", shell=True)

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

    commstr = "python2 ../../../parser/crosscompiler-diesel-dist.py -f kmeansgen.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o kmeans -i -dyn"
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    for i in range(numsamples):
        print "Running Iteration : ", i
        result_test = subprocess.check_output("./run.sh", shell=True)

        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1]) / 1000000
        print time_spent
        times.append(time_spent)

    track_time = geo_mean(times)
    print "Runtime with tracking: ", track_time

    print "Running with array optimization"
    times = []

    commstr = "python2 ../../../parser/crosscompiler-diesel-dist.py -f kmeansgen.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o kmeans -i -dyn -a"
    result_test = subprocess.check_output(commstr, shell=True)
    print result_test

    for i in range(numsamples):
        print "Running Iteration : ", i
        result_test = subprocess.check_output("./run.sh", shell=True)

        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1]) / 1000000
        print time_spent
        times.append(time_spent)

    opt_time = geo_mean(times)
    print "Runtime with opt: ", opt_time

    overhead = ((track_time - no_track_time) / no_track_time)
    opt_overhead = ((opt_time - no_track_time) / no_track_time)
    results[numsensors] = (overhead, opt_overhead)
