import os

def read(txtfile):
	f = open(txtfile,"r")
	lines = f.readlines()
	last = lines[-1]
	last = last.strip()
	time = last.split(":")[-1]
	time = time.replace(" ","")
	return float(time)

def avg(lst):
	assert(type(lst) is list)
	return sum(lst)/float(len(lst))
	

runs = {}
for executable in ["untracked","single_interval","multi-interval"]:
	runs[executable] = []
	for i in range(10):
		try:
			os.system("./" + executable + " > out.txt")
			time = read("out.txt")
			runs[executable].append(time)
			os.system("rm out.txt")
		except:
			pass



	print(executable + ": " , avg(runs[executable]))
