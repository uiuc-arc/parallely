import subprocess
import re
import pickle

nthreads = 8

# for inputsize in [1024, 2048, 4096, 8192]:

def genKmeansFromTemplate(inputsize):
    temp_file = "kmean.template"
    template_str = open(temp_file, 'r').readlines()
    with open("kmeans.par", "w") as fout:
        for line in template_str:
            newline = line.replace('NUMSENSORS', str(inputsize))
            fout.write(newline)

    temp_file = "__basic_go.template"
    template_str = open(temp_file, 'r').readlines()
    with open("_kmeans_go.txt", "w") as fout:
        for line in template_str:
            newline = line.replace('NUMSENSORS', str(inputsize))
            fout.write(newline)    

total_results = {}

genKmeansFromTemplate(1024)

# for inputsize in [1024, 2048, 4096, 8192]:
#     orig_times = []
#     new_times = []
#     opt_times = []
#     for i in xrange(samplesize):

#         commstr = """go build; ./kmean"""
#         result_test = subprocess.check_output(commstr, shell=True)

#         matches = re.findall("Elapsed time : .*\n", result_test)
#         time_spent = float(matches[0].split(' : ')[-1])
#         print time_spent
#         orig_times.append(time_spent)


#         genKmeansFromTemplate("kmeans_template.txt", inputsize)
#         # commstr = """go build; ./kmean"""
#         result_test = subprocess.check_output(commstr, shell=True)

#         matches = re.findall("Elapsed time : .*\n", result_test)
#         time_spent = float(matches[0].split(' : ')[-1])
#         print time_spent
#         new_times.append(time_spent)

#         genKmeansFromTemplate("kmeans_template_opt.txt", inputsize)
#         commstr = """go build; ./kmean"""
#         result_test = subprocess.check_output(commstr, shell=True)

#         matches = re.findall("Elapsed time : .*\n", result_test)
#         time_spent = float(matches[0].split(' : ')[-1])
#         print time_spent
#         opt_times.append(time_spent)
#         total_results[inputsize] = (orig_times, new_times, opt_times)

# with open("sensitivity-results.txt", "wb") as fout:        
#     pickle.dump(total_results, fout)
