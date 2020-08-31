import matplotlib
matplotlib.use('Agg')

import matplotlib.pyplot as plt
import numpy as np

benchmarkName = 'regression'

data = {1000: (35.4623878558934, 15.400499619350938), 1500: (44.855101316498406, 20.958483297735793), 500: (56.253179014675766, 25.81276001462917), 2000: (53.81771204830667, 26.596745066445244)}

names = ['Unoptimized', 'Optimized']
linestyles = ['--', ':']
colors = ['red', 'green']

sizes = []
times = ([], [])
comms = ([], [])

for datum in sorted(data.keys()):
  sizes.append(datum)
  # times[0].append(baseline/baseline)
  times[0].append(data[datum][0])
  times[1].append(data[datum][1])

# plt.figure(figsize=(2,2))
plt.title(benchmarkName, fontsize=20)
# plt.ticklabel_format(style='sci', axis='y', scilimits=(0,0))
for i in range(2):
  plt.plot(sizes, times[i],label=names[i],linestyle=linestyles[i],color=colors[i])
plt.legend(loc='upper left')
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
