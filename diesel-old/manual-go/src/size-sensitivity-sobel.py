import matplotlib.pyplot as plt

rawData = [(100,13.9,23.2,13.9),(200,47.5,82.8,52.6),(300,101.5,150.3,102.8),(400,167.8,319.4,170.2),(500,282.7,402.8,264.2)]

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
  comms[0].append(44*Dim2/1e6)
  comms[1].append(132*Dim2/1e6)
  comms[2].append((44*Dim2+160)/1e6)

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
