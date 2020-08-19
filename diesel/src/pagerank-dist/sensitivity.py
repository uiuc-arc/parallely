import subprocess
import re
import numpy as np

inputsets = {
    1: {"nodes": 8114, "slice": 1000, "file": "p2p-Gnutella09.txt"},
    2: {"nodes": 22687, "slice": 3000, "file": "p2p-Gnutella25.txt"},
    3: {"nodes": 36682, "slice": 5000, "file": "p2p-Gnutella30.txt"},
    4: {"nodes": 62586, "slice": 10000, "file": "p2p-Gnutella31.txt"},
}

src_size = 262144
num_sample = 5

data_set = {}

for inputgraph in inputsets:
    max_degrees = 100
    e_size = inputsets[inputgraph]["nodes"] * max_degrees
    slice_size = inputsets[inputgraph]["slice"]

    template_str = open("pagerank.template", 'r').readlines()
    with open("__temp_gen.par", "w") as fout:
        for line in template_str:
            newline = line.replace('__NUMEDGES__', str(e_size))
            newline = newline.replace('__NUMNODES__', str(inputsets[inputgraph]["nodes"]))
            newline = newline.replace('__FILENAME__', str(inputsets[inputgraph]["file"]))
            newline = newline.replace('__SLICESIZE__', str(slice_size))
            newline = newline.replace('__MAX_DEGREE_', str(max_degrees))
            fout.write(newline)

    template_str = open("__basic_go_main.tmpl", 'r').readlines()
    with open("__temp_main.txt", "w") as fout:
        for line in template_str:
            newline = line.replace('__NUMEDGES__', str(e_size))
            newline = newline.replace('__NUMNODES__', str(inputsets[inputgraph]["nodes"]))
            newline = newline.replace('__FILENAME__', str(inputsets[inputgraph]["file"]))
            newline = newline.replace('__SLICESIZE__', str(slice_size))
            newline = newline.replace('__MAX_DEGREE_', str(max_degrees))        
            fout.write(newline)

    template_str = open("__basic_go_worker.tmpl", 'r').readlines()
    with open("__temp_worker.txt", "w") as fout:
        for line in template_str:
            newline = line.replace('__NUMEDGES__', str(e_size))
            newline = newline.replace('__NUMNODES__', str(inputsets[inputgraph]["nodes"]))
            newline = newline.replace('__FILENAME__', str(inputsets[inputgraph]["file"]))
            newline = newline.replace('__SLICESIZE__', str(slice_size))
            newline = newline.replace('__MAX_DEGREE_', str(max_degrees))        
            fout.write(newline)        

    commstr = """python2 ../../../parser/crosscompiler-diesel-dist-rel.py -f __temp_gen.par -tm __temp_main.txt -tw __temp_worker.txt -o bfs.go -i"""
    result_test = subprocess.check_output(commstr, shell=True)
    print(result_test)

    no_track_time = []
    for i in range(num_sample):
        result_test = subprocess.check_output("./run.sh", shell=True)
        # print(result_test)
        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1])
        print(time_spent)
        no_track_time.append(time_spent)

    commstr = """python2 ../../../parser/crosscompiler-diesel-dist-rel.py -f __temp_gen.par -tm __temp_main.txt -tw __temp_worker.txt -o bfs.go -i -dyn"""
    result_test = subprocess.check_output(commstr, shell=True)
    print(result_test)

    track_time = []
    for i in range(num_sample):
        result_test = subprocess.check_output("./run.sh", shell=True)
        # print(result_test)
        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1])
        print(time_spent)
        track_time.append(time_spent)

    commstr = """python2 ../../../parser/crosscompiler-diesel-dist-rel.py -f __temp_gen.par -tm __temp_main.txt -tw __temp_worker.txt -o bfs.go -i -dyn -a"""
    result_test = subprocess.check_output(commstr, shell=True)
    print(result_test)

    opt_track_time = []
    for i in range(num_sample):
        result_test = subprocess.check_output("./run.sh", shell=True)
        # print(result_test)
        matches = re.findall("Elapsed time : .*\n", result_test)
        time_spent = float(matches[0].split(' : ')[-1])
        print(time_spent)
        opt_track_time.append(time_spent)

    print(no_track_time, track_time, opt_track_time)
    data_set[inputgraph] = (no_track_time, track_time, opt_track_time)


print("*************************")
print(data_set)
print("*************************")
    
#     commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o pagerank.go; go build"""
#     result_test = subprocess.check_output(commstr, shell=True)
#     print result_test

#     times = []
#     memory = []
#     for i in range(num_sample):
#         print "Running Iteration : ", i
#         result_test = subprocess.check_output("./pagerank-gen", shell=True)
#         print result_test
#         matches = re.findall("Elapsed time : .*\n", result_test)
#         time_spent = float(matches[0].split(' : ')[-1])
#         print time_spent
#         times.append(time_spent)
#         # matches2 = re.findall("Memory through channels : .*\n", result_test)
#         # mem_used = float(matches2[0].split(' : ')[-1])
#         # print mem_used
#         # memory.append(mem_used)

#     no_track_time = np.mean(times)
#     no_track_memory = mem_used

#     commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o pagerank.go -dyn; go build -tags instrument"""
#     result_test = subprocess.check_output(commstr, shell=True)
#     print result_test
#     result_test = subprocess.check_output("./pagerank-gen", shell=True)
#     print result_test
#     matches2 = re.findall("Memory through channels : .*\n", result_test)
#     mem_used = float(matches2[0].split(' : ')[-1])

#     commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o pagerank.go -dyn; go build"""
#     result_test = subprocess.check_output(commstr, shell=True)
#     print result_test

#     times = []
#     memory = []
#     for i in range(num_sample):
#         print "Running Iteration : ", i
#         result_test = subprocess.check_output("./pagerank-gen", shell=True)
#         print result_test
#         matches = re.findall("Elapsed time : .*\n", result_test)
#         time_spent = float(matches[0].split(' : ')[-1])
#         print time_spent
#         times.append(time_spent)

#         # matches2 = re.findall("Memory through channels : .*\n", result_test)
#         # mem_used = float(matches2[0].split(' : ')[-1])
#         # print mem_used
#         # memory.append(mem_used)

#     track_time = np.mean(times)
#     track_memory = mem_used

#     commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o pagerank.go -dyn -a; go build -tags instrument"""
#     result_test = subprocess.check_output(commstr, shell=True)
#     print result_test
#     result_test = subprocess.check_output("./pagerank-gen", shell=True)
#     print result_test
#     matches2 = re.findall("Memory through channels : .*\n", result_test)
#     mem_used = float(matches2[0].split(' : ')[-1])
#     print mem_used

#     commstr = """python ../../../parser/crosscompiler-diesel.py -f __temp_gen.par -t __temp_gen.txt -o pagerank.go -dyn -a; go build"""
#     result_test = subprocess.check_output(commstr, shell=True)
#     print result_test

#     times = []
#     memory = []
#     for i in range(num_sample):
#         print "Running Iteration : ", i
#         result_test = subprocess.check_output("./pagerank-gen", shell=True)
#         print result_test
#         matches = re.findall("Elapsed time : .*\n", result_test)
#         time_spent = float(matches[0].split(' : ')[-1])
#         print time_spent
#         times.append(time_spent)
#         # matches2 = re.findall("Memory through channels : .*\n", result_test)
#         # mem_used = float(matches2[0].split(' : ')[-1])
#         # print mem_used
#         # memory.append(mem_used)

#     opt_track_time = np.mean(times)
#     opt_track_memory = mem_used

#     data_set[inputgraph] = (no_track_time, track_time, opt_track_time,
#                             no_track_memory, track_memory, opt_track_memory)
#     print data_set

