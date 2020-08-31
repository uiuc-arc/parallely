
tracked = [1317723544, 1518565615, 1831797095, 2783908673]
uninstrumented = [1234536631, 1534821848, 1762227263, 2414685487]





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
