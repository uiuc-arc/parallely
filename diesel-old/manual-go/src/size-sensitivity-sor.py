import matplotlib.pyplot as plt

rawData = [(100,107.0,163.0,109.0),(200,407.0,643.0,433.0),(300,850.0,1256.,891.0),(400,1433.,2072.,1486.),(500,2156.,3303.,2274.)]

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
  comms[0].append(440*Dim2/1e6)
  comms[1].append(1320*Dim2/1e6)
  comms[2].append((440*Dim2+1600)/1e6)

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
