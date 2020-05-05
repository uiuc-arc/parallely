import matplotlib.pyplot as plt

benchmarkName = 'Scale'

# data = {2: (1582244716.5, 1822829189.15, 1624374035.65, 8388608.0, 16777216.0, 8388672.0),
#         4: (8381588392.0, 8289309085.75, 8820298747.1, 33554432.0, 67108864.0, 33554496.0),
#         8: (37375237429.65, 37388729992.6, 37423974429.7, 134217728.0, 268435456.0, 134217792.0),
#         16: (144266163663.5, 152786712083.8, 155771628638.75, 536870912.0, 1073741824.0, 536870976.0)}

# data = {8: (1713570389.8, 1728842023.0, 1748779914.85, 134217728.0, 268435456.0, 134217792.0),
#         16: (8305078369.25, 8442914958.65, 8350695066.35, 536870912.0, 1073741824.0, 536870976.0),
#         2: (94419846.55, 94371646.45, 83042953.5, 8388608.0, 16777216.0, 8388672.0),
#         4: (336615840.55, 368953822.05, 380192136.05, 33554432.0, 67108864.0, 33554496.0)}

data = {8: (218859706.4, 199388810.2, 199422034.2, 134217728.0, 268435456.0, 134217792.0), 16: (803198708.0, 797988940.0, 973658599.0, 536870912.0, 1073741824.0, 536870976.0), 2: (22923054.0, 22246132.4, 21935514.4, 8388608.0, 16777216.0, 8388672.0), 4: (74712630.6, 61715476.4, 88939365.8, 33554432.0, 67108864.0, 33554496.0)}

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
plt.legend(loc='upper left')
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
