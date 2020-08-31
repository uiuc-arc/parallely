
tracked = [1299277499, 1585449742, 1834624075, 2538168609]
uninstrumented = [1220158201, 1558396009, 1825615984, 2498836172]





def geoMean(lst):
	float_lst = [float(x) for x in lst]
	from scipy.stats import gmean
	return (gmean(float_lst))

def StdDev(lst):
	float_lst = [float(x) for x in lst]
	return std(float_lst)



overheads = []
for i in range(len(tracked)):
	ovhd = float(tracked[i])/float(uninstrumented[i])
	overheads.append(ovhd)
print geoMean(overheads)
