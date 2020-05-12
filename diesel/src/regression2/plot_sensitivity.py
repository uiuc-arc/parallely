import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt

benchmarkName = 'Regression'

data = {10000:[1518361.32,  3891876.39,  1902442.66, 160160,  480480, 160800],
        20000:[3917359.50,  8320660.74,  4645970.45, 320160,  960480, 320800],
        30000:[5962077.99, 12313881.57,  7204270.29, 480160, 1440480, 480800],
        40000:[8223463.96, 18180880.07,  9324975.73, 640160, 1920480, 640800],
        50000:[9852202.10, 19518731.84, 11313165.00, 800160, 2400480, 800800]}

names = ['Baseline', 'Unoptimized', 'Optimized']
linestyles = ['-', '--', ':']
colors = ['orange', 'red', 'green']

sizes = []
times = ([], [], [])
comms = ([], [], [])

for datum in sorted(data.keys()):
  sizes.append(datum)
  times[0].append(data[datum][0]/1e6)
  times[1].append(data[datum][1]/1e6)
  times[2].append(data[datum][2]/1e6)
  comms[0].append(data[datum][3]/1e6)
  comms[1].append(data[datum][4]/1e6)
  comms[2].append(data[datum][5]/1e6)

# plt.figure(figsize=(2,2))
plt.title(benchmarkName, fontsize=20)
# plt.ticklabel_format(style='sci', axis='y', scilimits=(0,0))
for i in range(3):
  plt.plot(sizes,times[i],label=names[i],linestyle=linestyles[i],color=colors[i])
# plt.legend(loc='upper left')
plt.xlabel('Input Size', fontsize=20)
plt.ylabel('Time (ms)', fontsize=20)
plt.xticks(fontsize=20, rotation=90)
plt.yticks(fontsize=20)
plt.tight_layout()
plt.savefig('times-{}.png'.format(benchmarkName))
plt.close('all')

# plt.figure(figsize=(2,2))
plt.title(benchmarkName, fontsize=20)
# plt.ticklabel_format(style='sci', axis='y', scilimits=(0,0))
for i in range(3):
  plt.plot(sizes,comms[i],label=names[i],linestyle=linestyles[i],color=colors[i])
# plt.legend(loc='upper left')
plt.xlabel('Input Size', fontsize=20)
plt.ylabel('Data (MB)', fontsize=20)
plt.xticks(fontsize=20, rotation=90)
plt.yticks(fontsize=20)
plt.tight_layout()
plt.savefig('comms-{}.png'.format(benchmarkName))
plt.close('all')
