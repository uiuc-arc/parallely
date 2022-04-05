import matplotlib.pyplot as plt

benchmarkName = 'IncomeSVM'

#units are in ms
data = {20000: (10.31, 11.03),
        50000: (22.67, 23.86),
        100000: (44.15, 46.02),
        150000: (65.41, 67.27)}

names = ['Baseline', 'Unoptimized']
linestyles = ['-', '--']
colors = ['orange', 'red' ]

sizes = []
times = ([], [], [])
comms = ([], [], [])

for datum in sorted(data.keys()):
  sizes.append(datum)
  times[0].append(data[datum][0])
  times[1].append(data[datum][1])
  #times[2].append(data[datum][2]/1e9)
  '''
  comms[0].append(data[datum][3]/1e6)
  comms[1].append(data[datum][4]/1e6)
  comms[2].append(data[datum][5]/1e6)
  '''

# plt.figure(figsize=(2,2))
plt.title(benchmarkName, fontsize=20)
# plt.ticklabel_format(style='sci', axis='y', scilimits=(0,0))
for i in range(2):
  plt.plot(sizes,times[i],label=names[i],linestyle=linestyles[i],color=colors[i])
# plt.legend(loc='upper left')
plt.xlabel('Input Size', fontsize=20)
plt.ylabel('Time (ms)', fontsize=20)
plt.xticks(fontsize=20, rotation=90)
plt.yticks(fontsize=20)
plt.tight_layout()
plt.savefig('times-{}.png'.format(benchmarkName))
plt.close('all')
