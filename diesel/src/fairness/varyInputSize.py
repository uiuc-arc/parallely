import os

numTrials = 2

def avg(lst):
	float_lst = [float(x) for x in lst]
	return int(sum(float_lst)/float(len(float_lst)))


def geoMean(lst):
	float_lst = [float(x) for x in lst]
	return (gmean(float_lst))

def StdDev(lst):
	float_lst = [float(x) for x in lst]
	return std(float_lst)


def change_input(filename,numtrials):
	assert(type(numtrials) is int)
	newfilename = getTestFileName(filename,numtrials)
	newExName = getTestExecutableName(filename,numtrials)
	newStr = "const\ datasize\ =\ " + str(numtrials)
	replace_cmd = 'sed ' + 's/const\ datasize\ =\ 80000/' + newStr + '/g ' + filename + " > " + newfilename
	os.system(replace_cmd)
	compile_cmd = "go build -o " + newExName + " " + newfilename
	os.system(compile_cmd)


def getTestFileName(filename,numtrials):
	return "test" + str(numtrials) + ".go"

def getTestExecutableName(filename,numtrials):
	tmp = getTestFileName(filename,numtrials)
	return tmp[:-3]

def cleanup(filename,numtrials):
	assert(type(numtrials) is int)
	newfilename = "test" + str(numtrials) + ".go"
	newExName =newfilename[:-3] 
	remove_cmd = "rm " + newfilename + " " + newExName
	os.system(remove_cmd)

#returns the time in nanoseconds
def fetchOutputString(txtfile):
	f = open(txtfile)
	lines = f.readlines()
	assert(len(lines)==1)
	st = lines[0]
	st = st.split(" ")[-1]
	num = int(st)
	return num
	

#inputSizeRange = [10000,20000,40000,80000,120000]
inputSizeRange = [1000000,2000000,4000000,8000000]
sizes_str = str(inputSizeRange)
#inputSizeRange = [80,800,8000,80000,8000000,80000000]
def runForMultipleInputSizes(filename):
	times = []
	for inputSize in inputSizeRange:
		execName = getTestExecutableName(filename,inputSize)
		testFileName = execName + ".txt"
		change_input(filename,inputSize)
		runs = []
		for i in range(numTrials):
			run_cmd = "./" + execName + " > " + testFileName
			os.system(run_cmd)
			#wait_cmd = "sleep 1"
			#os.system(wait_cmd)
			runtime = fetchOutputString(testFileName)
			runs.append(runtime)
		times.append(avg(runs))
	os.system("rm test*")
	return times
		


output_str = ""
org_dir = os.getcwd()
for dirs in ["DecisionTree","Hiring","NN","SVM"]:
	os.chdir(dirs)
	intermediate_dir = os.getcwd()
	for method in ["specializedClass","specializedUninstrumented"]: 
		os.chdir(method)
		if method == "specializedClass":
			#print(os.getcwd())
			times = runForMultipleInputSizes("specialized.go")
			line = dirs + " tracked: " + str(times) + "\n"
			output_str += line
			os.chdir(intermediate_dir)
		else:
			#print(os.getcwd())
			times = runForMultipleInputSizes("uninstrumented.go")
			line = dirs + " uninstrumented: " + str(times) + "\n"
			output_str += line
			os.chdir(intermediate_dir)
	os.chdir(org_dir)




text_file = open("VariedInputSize.txt", "w")
output_str = sizes_str + "\n" + output_str
text_file.write(output_str)
text_file.close()
