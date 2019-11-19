import matplotlib.pyplot as plt

# data = {4: (1, 2, 3, 1, 2, 3), 8: (1, 2, 3, 1, 2, 3), 16: (1, 2, 3, 1, 2, 3)}

data = {1: (169022414.375, 216057207.8, 153246880.425, 7291200.0, 14582400.0, 7292800.0),
        2: (391415299.15, 606179433.2, 414396086.475, 20549600.0, 41099200.0, 20551200.0),
        3: (832640887.725, 1011389912.475, 726658919.275, 33345600.0, 66691200.0, 33347200.0),
        4: (1196522856.7, 1922401998.05, 1366976605.9, 58068800.0, 116137600.0, 58070400.0)}

names = ['Baseline', 'Unoptimized', 'Optimized']
linestyles = ['-', '--', '-.']
colors = ['orange', 'red', 'green']

sizes = []
times = ([], [], [])
comms = ([], [], [])

for datum in sorted(data.keys()):
  sizes.append(datum)
  times[0].append(data[datum][0])
  times[1].append(data[datum][1])
  times[2].append(data[datum][2])
  comms[0].append(data[datum][3])
  comms[1].append(data[datum][4])
  comms[2].append(data[datum][5])

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
