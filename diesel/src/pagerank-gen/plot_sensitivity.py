import matplotlib.pyplot as plt

benchmarkName = 'pagerank'

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

plt.figure(figsize=(2,2))
plt.rcParams.update({'font.size': 12})
plt.ticklabel_format(style='sci', axis='y', scilimits=(0,0))
for i in range(3):
  plt.plot(sizes,times[i],label=names[i],linestyle=linestyles[i],color=colors[i])
# plt.legend(loc='upper left')
plt.xlabel('Input Size')
plt.ylabel('Time (ms)')
plt.tight_layout()
plt.savefig('times-{}.png'.format(benchmarkName))
plt.close()

plt.figure(figsize=(2,2))
plt.rcParams.update({'font.size': 12})
plt.ticklabel_format(style='sci', axis='y', scilimits=(0,0))
for i in range(3):
  plt.plot(sizes,comms[i],label=names[i],linestyle=linestyles[i],color=colors[i])
# plt.legend(loc='upper left')
plt.xlabel('Input Size')
plt.ylabel('Data (MB)')
plt.tight_layout()
plt.savefig('comms-{}.png'.format(benchmarkName))

'''
# fig, axs = plt.subplots(2, 1)
for i in range(3):
    plt.plot(sizes, times[i], label=names[i], linestyle=linestyles[i], color=colors[i])
    # plt.plot(sizes, comms[i], label=names[i], linestyle=linestyles[i], color=colors[i])
# for i in range(2):
#   axs[i].legend(loc='upper left')
#   axs[i].set_xlabel('Input Size')

plt.legend(loc='upper left')
plt.xlabel('Input Size (Number of Nodes)')
# 1: {"nodes": 8114, "slice": 1000, "file": "p2p-Gnutella09.txt"},
# 2: {"nodes": 22687, "slice": 3000, "file": "p2p-Gnutella25.txt"},
# 3: {"nodes": 36682, "slice": 5000, "file": "p2p-Gnutella30.txt"},
# 4: {"nodes": 62586, "slice": 10000, "file": "p2p-Gnutella31.txt"},
plt.ylabel('Time (s)')
plt.xticks(sorted(data.keys()), ['8114', '22687', '36682', '62586'])
plt.tight_layout()
plt.savefig("time.png")
plt.close("all")

for i in range(3):
    plt.plot(sizes, comms[i], label=names[i], linestyle=linestyles[i], color=colors[i])
plt.legend(loc='upper left')
plt.ylabel('Communicated Data (MB)')
plt.xlabel('Input Size (Number of Nodes)')
plt.xticks(sorted(data.keys()), ['8114', '22687', '36682', '62586'])
plt.tight_layout()
plt.savefig("memory.png")
'''
