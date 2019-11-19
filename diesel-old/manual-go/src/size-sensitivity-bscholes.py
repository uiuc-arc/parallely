import matplotlib.pyplot as plt

rawData = [(2000,1.74,4.30,2.04),(4000,3.73,7.16,4.08),(6000,5.36,10.09,5.82),(8000,6.55,12.65,7.00),(10000,8.76,17.76,11.46)]

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
  Dim2 = datum[0]
  comms[0].append(40*Dim2/1e6)
  comms[1].append(120*Dim2/1e6)
  comms[2].append((40*Dim2+160)/1e6)

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
