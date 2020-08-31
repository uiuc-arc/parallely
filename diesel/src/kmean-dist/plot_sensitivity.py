import matplotlib
matplotlib.use('Agg')

import matplotlib.pyplot as plt
import numpy as np

benchmarkName = 'KMeans'

data = {248: (0.3797132382950256, 0.02111444421293428), 512: (0.38087820476718415, 0.003284102709640509), 1024: (0.416349310392601, 0.016938769212795562), 2048: (0.46490315535838184, 0.04990946215861019)}

names = ['Unoptimized', 'Optimized']
linestyles = ['--', ':']
colors = ['red', 'green']

sizes = []
times = ([], [])
comms = ([], [])

for datum in sorted(data.keys()):
  sizes.append(datum)
  # baseline = np.mean(data[datum][0])/1e9
  # times[0].append(baseline/baseline)
  times[0].append(data[datum][0] * 100)
  times[1].append(data[datum][1] * 100)

# plt.figure(figsize=(2,2))
plt.title(benchmarkName, fontsize=20)
# plt.ticklabel_format(style='sci', axis='y', scilimits=(0,0))
for i in range(2):
  plt.plot(sizes, times[i],label=names[i],linestyle=linestyles[i],color=colors[i])
# plt.legend(loc='upper left')
plt.xlabel('Input Size', fontsize=20)
plt.ylabel('Overhead%', fontsize=20)
plt.xticks(fontsize=20, rotation=90)
plt.yticks(fontsize=20)
plt.tight_layout()
plt.savefig('times-{}.png'.format(benchmarkName))
plt.close('all')

# # plt.figure(figsize=(2,2))
# plt.title(benchmarkName, fontsize=20)
# # plt.ticklabel_format(style='sci', axis='y', scilimits=(0,0))
# for i in range(3):
#   plt.plot(sizes,comms[i],label=names[i],linestyle=linestyles[i],color=colors[i])
# # plt.legend(loc='upper left')
# plt.xlabel('Input Size', fontsize=20)
# plt.ylabel('Data (MB)', fontsize=20)
# plt.xticks(fontsize=20, rotation=90)
# plt.yticks(fontsize=20)
# plt.tight_layout()
# plt.savefig('comms-{}.png'.format(benchmarkName))
# plt.close('all')
