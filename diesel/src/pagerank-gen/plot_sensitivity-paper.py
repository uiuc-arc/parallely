import matplotlib
import matplotlib.pyplot as plt

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

plt.title(benchmarkName, fontsize=25, fontweight='bold')
for i in range(3):
  plt.plot(sizes,comms[i],label=names[i],linestyle=linestyles[i],color=colors[i],linewidth=3)
# plt.legend(loc='upper left')
plt.xlabel('Input Size')
plt.ylabel('Data (MB)')
plt.xticks(list(data.keys()), ['8K','23K','37K','63K'])
#plt.yticks(range(4), ['0', '1', '2', '3'])
plt.tight_layout()
plt.savefig('comms-{}-paper.png'.format(benchmarkName), pad_inches = 0)
plt.close('all')

