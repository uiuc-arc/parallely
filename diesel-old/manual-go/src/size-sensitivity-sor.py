import matplotlib.pyplot as plt

rawData = [(100,0.107,0.163,0.109),(200,0.407,0.643,0.433),(300,0.850,1.256,0.891),(400,1.433,2.072,1.486),(500,2.156,3.303,2.274)]

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
  comms[0].append(440*Dim2)
  comms[1].append(1320*Dim2)
  comms[2].append(440*Dim2+1600)

fig, axs = plt.subplots(2,1)
for i in range(3):
  axs[0].plot(sizes,times[i],label=names[i],linestyle=linestyles[i],color=colors[i])
  axs[1].plot(sizes,comms[i],label=names[i],linestyle=linestyles[i],color=colors[i])
for i in range(2):
  axs[i].legend(loc='upper left')
  axs[i].set_xlabel('Input Size')
axs[0].set_ylabel('Time (s)')
axs[1].set_ylabel('Comm. Data (Bytes)')
fig.set_size_inches(4,6)
plt.tight_layout()
plt.show()
