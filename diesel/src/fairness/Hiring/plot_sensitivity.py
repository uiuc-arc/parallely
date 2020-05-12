import matplotlib.pyplot as plt

benchmarkName = 'Hiring'

#units are in ms!
tracked = [1295194994, 1533673560, 1854841122, 2533516492]
uninstrumented = [1238203692, 1489975480, 1775365791, 2448989614]


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
plt.ylabel('Time (s)', fontsize=20)
plt.xticks(fontsize=20, rotation=90)
plt.yticks(fontsize=20)
plt.tight_layout()
plt.savefig('times-{}.png'.format(benchmarkName))
plt.close('all')
