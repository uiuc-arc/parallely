import os

def avg(lst):
	float_lst = [float(x) for x in lst]
	return int(sum(float_lst)/float(len(float_lst)))


def fetchOutputString(txtfile):
	f = open(txtfile)
	lines = f.readlines()
	assert(len(lines)==1)
	st = lines[0]
	st = st.split(" ")[-1]
	num = int(st)
	return num
	


cwd = os.getcwd()
numTrials = 25
for d in ["tracked","uninstrumented"]:
	os.chdir(d)
	if d == "tracked":
		tracked_times = []
		for i in range(numTrials):
			os.system("./fused > tracked_out.txt")
			time = fetchOutputString("tracked_out.txt")
			tracked_times.append(time)
			os.system("rm *_out.txt")
		#print(tracked_times)
	elif d == "uninstrumented":
		uninstrumented_times = []
		for i in range(numTrials):
			os.system("./uninstrumented > uninstrumented_out.txt")
			time = fetchOutputString("uninstrumented_out.txt")
			uninstrumented_times.append(time)
			os.system("rm *_out.txt")
		#print(uninstrumented_times)

	os.chdir(cwd)




result = "tracked: " + str(avg(tracked_times)) + "\n" + "uninstrumented: " + str(avg(uninstrumented_times))
#print("\n\n")
print(result)




