import matplotlib
matplotlib.use('Agg')

import matplotlib.pyplot as plt
import numpy as np

benchmarkName = 'SOR'

data = {200: (0.00011042792950580116, 5.883356799603508e-06), 300: (0.00012495085391220153, 8.683904918742102e-06), 400: (0.00012871759019120673, 8.030793982149319e-06), 100: (6.992116158573908e-05, 3.954419222109069e-06)}

names = ['Baseline', 'Diesel']
linestyles = ['--', ':']
colors = ['red', 'green']

sizes = []
times = ([], [])
comms = ([], [])

for datum in sorted(data.keys()):
  sizes.append(datum)
  # times[0].append(1)
  times[0].append(data[datum][0] * 1000000)
  times[1].append(data[datum][1] * 1000000)

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
