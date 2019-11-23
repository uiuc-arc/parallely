import matplotlib
import math
import matplotlib.pyplot as plt
from matplotlib.ticker import StrMethodFormatter

benchmarkName = 'PageRank'

plt.rcParams['axes.labelsize'] = 25
plt.rcParams['axes.labelweight'] = 'bold'
plt.rcParams["font.size"] = 25
plt.rcParams["font.weight"] = "bold"
plt.rcParams["axes.labelweight"] = "bold"
plt.rcParams['axes.titlepad'] = 20
matplotlib.rcParams.update({'legend.fontsize': 20})
matplotlib.rcParams.update({'font.weight': 'bold'})
matplotlib.rcParams.update({'axes.linewidth': 4})

names = ['Unoptimized','Optimized']
linestyles = ['--',':']
colors = ['red','green']

data = [[0.999999999, 0.9999999899999993, 0.999999961999997, 0.9999998849999907, 0.9999996749999736, 0.9999990999999283, 0.9999975119998056, 0.9999931489994726, 0.9999811289985642, 0.9999480609961008],
[0.999999999, 0.9999999899999993, 0.9999999089999927, 0.999999179999935, 0.9999926189994227, 0.9999335699948698, 0.9994021289544273, 0.9946191595952262, 0.9515724354054622, 0.564151918085008]]

plt.title(benchmarkName, fontsize=25, fontweight='bold')
for i in range(2):
  plt.plot(range(1,9),data[i][:8],label=names[i],linestyle=linestyles[i],color=colors[i],linewidth=3)
plt.xlabel("Iteration")
plt.ylabel("Reliability")
plt.xticks(range(1,9),[str(s) for s in range(1,9)])

plt.tight_layout()
plt.savefig('reliability-{}.png'.format(benchmarkName), pad_inches = 0)
plt.close('all')
