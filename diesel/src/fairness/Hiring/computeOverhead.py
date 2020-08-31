
tracked = [1295194994, 1533673560, 1854841122, 2533516492]
uninstrumented = [1238203692, 1489975480, 1775365791, 2448989614]





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
