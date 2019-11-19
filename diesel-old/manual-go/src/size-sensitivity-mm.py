import matplotlib.pyplot as plt

rawData = [(50,7.7,12.2,7.8),(100,24.5,41.7,25.3),(150,52.2,88.8,55.6),(200,88.1,142.8,94.7),(250,135.2,219.0,143.2)]

names = ['Baseline','Unoptimized','Optimized']
linestyles = ['-','--','-.']
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

fig, axs = plt.subplots(2,1)
for i in range(3):
  axs[0].plot(sizes,times[i],label=names[i],linestyle=linestyles[i],color=colors[i])
  axs[1].plot(sizes,comms[i],label=names[i],linestyle=linestyles[i],color=colors[i])
for i in range(2):
  axs[i].legend(loc='upper left')
  axs[i].set_xlabel('Input Size')
axs[0].set_ylabel('Time (ms)')
axs[1].set_ylabel('Communicated Data (MB)')
fig.set_size_inches(4,6)
plt.tight_layout()
plt.show()
