import subprocess
import re
import pickle

nthreads = 8

# for inputsize in [1024, 2048, 4096, 8192]:
def genKmeansFromTemplate(temp_file, inputsize):
    datasize = inputsize*2
    centeridsize = 8
    centersize = centeridsize * 2
    TOTALSIZEQ = datasize + centersize + 1 + centersize + 1
    TOTALSIZE0 = datasize + centersize *3
    mywork = inputsize/nthreads

    template_str = open(temp_file, 'r').readlines()
    with open("kmeans_gen.go", "w") as fout:
        for line in template_str:
            newline = line.replace('DATASIZE', str(datasize))
            newline = newline.replace('CENTERIDSIZE', str(centeridsize))
            newline = newline.replace('CENTERSIZE', str(centersize))
            newline = newline.replace('TOTALSIZEQ', str(TOTALSIZEQ))
            newline = newline.replace('TOTALSIZE0', str(TOTALSIZE0))
            newline = newline.replace('MYWORK', str(mywork))
            fout.write(newline)

total_results = {}
total_results_memory = {}

samplesize = 32

for inputsize in [1024, 2048, 4096, 8192]:
    orig_times = []
    new_times = []
    opt_times = []

    genKmeansFromTemplate("kmeans_template_notrack.txt", inputsize)
    commstr = """go build -tags instrument; ./kmean"""
    result_test = subprocess.check_output(commstr, shell=True)
    # print result_test
    matches2 = re.findall("Memory through channels : .*\n", result_test)
    # print matches2[0].split(' : ')[-1]
    mem_used_0 = float(matches2[0].split(' : ')[-1])

    genKmeansFromTemplate("kmeans_template.txt", inputsize)
    commstr = """go build -tags instrument; ./kmean"""
    result_test = subprocess.check_output(commstr, shell=True)
    matches2 = re.findall("Memory through channels : .*\n", result_test)
    mem_used_1 = float(matches2[0].split(' : ')[-1])

    genKmeansFromTemplate("kmeans_template_opt.txt", inputsize)
    commstr = """go build -tags instrument; ./kmean"""
    result_test = subprocess.check_output(commstr, shell=True)
    matches2 = re.findall("Memory through channels : .*\n", result_test)
    mem_used_2 = float(matches2[0].split(' : ')[-1])

    total_results_memory[inputsize] = (mem_used_0, mem_used_1, mem_used_2)
    print total_results_memory
    
#     for i in xrange(samplesize):
#         genKmeansFromTemplate("kmeans_template_notrack.txt", inputsize)
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

with open("sensitivity-results-memory.txt", "wb") as fout:        
    pickle.dump(total_results_memory, fout)    
