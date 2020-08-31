import matplotlib.pyplot as plt

benchmarkName = 'BFS'

data = {8114: (265492600.425, 221900730.1, 178767290.475, 7291200.0,
        14582400.0, 7292800.0), 22687: (418209894.025, 660249377.825,
                                        422191235.475, 20549600.0, 41099200.0, 20551200.0), 36682:
 (738363828.325, 1020799504.55, 804386645.95, 33345600.0, 66691200.0,
  33347200.0), 62586: (1233804882.225, 1932927366.675, 1407868098.75,
                       58068800.0, 116137600.0, 58070400.0)}

# {8114: (169022414.375, 216057207.8, 153246880.425, 7291200.0, 14582400.0, 7292800.0),
#         22687: (391415299.15, 606179433.2, 414396086.475, 20549600.0, 41099200.0, 20551200.0),
#         36682: (832640887.725, 1011389912.475, 726658919.275, 33345600.0, 66691200.0, 33347200.0),
#         62586: (1196522856.7, 1922401998.05, 1366976605.9, 58068800.0, 116137600.0, 58070400.0)}

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
