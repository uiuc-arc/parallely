import matplotlib.pyplot as plt

benchmarkName = 'PageRank'

data = {8114: (135290433.55, 193474567.7, 139523218.175, 5832960.0, 11665920.0, 5834240.0),
        22687: (324506229.775, 488433193.25, 311194275.875, 16439680.0, 32879360.0, 16440960.0),
        36682: (473615164.6, 746029766.65, 490723435.125, 26676480.0, 53352960.0, 26677760.0),
        62586: (904140282.3, 1506761906.6, 935868980.525, 46455040.0, 92910080.0, 46456320.0)}

names = ['Baseline', 'Unoptimized', 'Optimized']
linestyles = ['-', '--', ':']
colors = ['orange', 'red', 'green']

sizes = []
times = ([], [], [])
comms = ([], [], [])

for datum in sorted(data.keys()):
  sizes.append(datum)
  times[0].append(data[datum][0]/1e9)
  times[1].append(data[datum][1]/1e9)
  times[2].append(data[datum][2]/1e9)
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
