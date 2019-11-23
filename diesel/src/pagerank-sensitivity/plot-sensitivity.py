import matplotlib
import math
import matplotlib.pyplot as plt
from matplotlib.ticker import StrMethodFormatter

plt.rcParams['axes.labelsize'] = 25
plt.rcParams['axes.labelweight'] = 'bold'
plt.rcParams["font.size"] = 25
plt.rcParams["font.weight"] = "bold"
plt.rcParams["axes.labelweight"] = "bold"
plt.rcParams['axes.titlepad'] = 20
matplotlib.rcParams.update({'legend.fontsize': 20})
matplotlib.rcParams.update({'font.weight': 'bold'})
matplotlib.rcParams.update({'axes.linewidth': 4})

orig = [0.999999999, 0.9999999899999993, 0.999999961999997,
        0.9999998849999907, 0.9999996749999736, 0.9999990999999283,
        0.9999975119998056, 0.9999931489994726, 0.9999811289985642,
        0.9999480609961008]

opt = [0.999999999, 0.9999999899999993, 0.9999999089999927,
       0.999999179999935, 0.9999926189994227, 0.9999335699948698,
       0.9994021289544273, 0.9946191595952262, 0.9515724354054622,
       0.564151918085008]

# plt.plot([1 + i for i in range(8)], [-math.log(1-orig[i]) for i in range(8)], label="unoptimized")
# plt.plot([1 + i for i in range(8)], [-math.log(1-opt[i]) for i in range(8)], label="optimized")

fig = plt.figure()
ax = fig.add_subplot(111)
ax.plot([1 + i for i in range(8)], [orig[i] for i in range(8)], label="unoptimized", marker='o', linewidth=2)
ax.plot([1 + i for i in range(8)], [opt[i] for i in range(8)], label="optimized", marker='o', linewidth=2)

# ax.set_yscale('log')
plt.legend()
plt.ylim(0.99, 1.001)
# plt.xlim(1, 8.2)
plt.ylabel("Reliability")
plt.xlabel("Iteration")
# plt.gca().yaxis.grid(True)
plt.xticks(range(1, 9))
plt.tight_layout()
plt.savefig("rel.png", pad_inches = 0)
