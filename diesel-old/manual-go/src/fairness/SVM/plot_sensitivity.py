import matplotlib.pyplot as plt

benchmarkName = 'SVM_Income'

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
plt.figure(figsize=(2,2))
plt.rcParams.update({'font.size': 12})
plt.ticklabel_format(style='sci', axis='y', scilimits=(0,0))
for i in range(2):
  plt.plot(sizes,times[i],label=names[i],linestyle=linestyles[i],color=colors[i])
# plt.legend(loc='upper left')
plt.xlabel('Input Size')
plt.ylabel('Time (ms)')
plt.tight_layout()
plt.savefig('times-{}.png'.format(benchmarkName))
plt.close()

'''
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
