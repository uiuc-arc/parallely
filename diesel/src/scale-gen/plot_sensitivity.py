import matplotlib.pyplot as plt

# data = {4: (1, 2, 3, 1, 2, 3), 8: (1, 2, 3, 1, 2, 3), 16: (1, 2, 3, 1, 2, 3)}

data = {2: (2182202570.3, 2414151423.1, 2285356977.9, 8388608.0, 16777216.0, 8388672.0),
        4: (9446486885.3, 10332666758.0, 10303837240.5, 33554432.0, 67108864.0, 33554496.0),
        8: (38675243736.5, 40349603641.1, 36467243717.1, 134217728.0, 268435456.0, 134217792.0),
        16: (157882553005.4, 170195137187.5, 159170102931.3, 536870912.0, 1073741824.0, 536870976.0)}

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
plt.xlabel('Input Size')
plt.ylabel('Time (s)')
plt.xticks(sorted(data.keys()), ['2x', '4x', '8x', '16x'])
plt.tight_layout()
plt.savefig("time.png")
plt.close("all")

for i in range(3):
    plt.plot(sizes, comms[i], label=names[i], linestyle=linestyles[i], color=colors[i])
plt.legend(loc='upper left')
plt.ylabel('Communicated Data (MB)')
plt.xlabel('Input Size')
plt.xticks(sorted(data.keys()), ['2x', '4x', '8x', '16x'])
plt.tight_layout()
plt.savefig("memory.png")
