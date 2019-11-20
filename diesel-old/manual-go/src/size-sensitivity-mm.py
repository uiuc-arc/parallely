import matplotlib.pyplot as plt

benchmarkName = 'mm'

rawData = [(50,7.7,12.2,7.8),(100,24.5,41.7,25.3),(150,52.2,88.8,55.6),(200,88.1,142.8,94.7),(250,135.2,219.0,143.2)]

names = ['Baseline','Unoptimized','Optimized']
linestyles = ['-','--',':']
colors = ['orange','red','green']

sizes = []
times = ([],[],[])
comms = ([],[],[])
for datum in rawData:
  sizes.append(datum[0])
  for i in range(1,4):
    times[i-1].append(datum[i])
  Dim2 = datum[0]**2
  comms[0].append(84*Dim2/1e6)
  comms[1].append(252*Dim2/1e6)
  comms[2].append((84*Dim2+240)/1e6)

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
