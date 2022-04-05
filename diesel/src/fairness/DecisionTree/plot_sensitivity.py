import matplotlib.pyplot as plt

benchmarkName = 'DecisionTree'

#units are in ms!
tracked = [1317723544, 1518565615, 1831797095, 2783908673]
uninstrumented = [1234536631, 1534821848, 1762227263, 2414685487]


tracked = [float(x) for x in tracked]
uninstrumented = [float(x) for x in uninstrumented]
data = {1000000: (uninstrumented[0], tracked[0],),
        2000000: (uninstrumented[1], tracked[1],),
        4000000: (uninstrumented[2], tracked[2]),
        8000000: (uninstrumented[3], tracked[3])}





names = ['Uninstrumented', 'Tracked']
linestyles = ['-', '--']
colors = ['orange', 'red' ]

sizes = []
times = ([], [], [])
comms = ([], [], [])

for datum in sorted(data.keys()):
  sizes.append(datum)
  times[0].append(data[datum][0]/1e6)
  times[1].append(data[datum][1]/1e6)
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
plt.xticks([1000000,2000000,3000000,4000000,5000000,6000000,7000000,8000000], ['1M','2M','3M','4M','5M','6M','7M','8M'], fontsize=20, rotation=90)
plt.yticks(fontsize=20)
plt.tight_layout()
plt.savefig('times-{}.png'.format(benchmarkName))
plt.close('all')