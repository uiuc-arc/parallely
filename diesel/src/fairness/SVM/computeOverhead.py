
tracked = [1300241798, 1508535761, 1892845297, 2559020386]
uninstrumented = [1251105192, 1457639129, 1847797433, 2510187457]




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
