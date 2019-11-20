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

plt.rcParams.update({'font.size': 18})

for i in range(3):
  plt.plot(sizes,times[i],label=names[i],linestyle=linestyles[i],color=colors[i])
plt.legend(loc='upper left')
plt.xlabel('Input Size')
plt.ylabel('Time (ms)')
plt.tight_layout()
plt.savefig('times-sobel.png')
plt.close()

for i in range(3):
  plt.plot(sizes,comms[i],label=names[i],linestyle=linestyles[i],color=colors[i])
plt.legend(loc='upper left')
plt.xlabel('Input Size')
plt.ylabel('Communicated Data (MB)')
plt.tight_layout()
plt.savefig('comms-sobel.png')
