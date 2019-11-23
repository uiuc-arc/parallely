import matplotlib
import matplotlib.pyplot as plt

benchmarkName = 'SOR'

plt.rcParams['axes.labelsize'] = 25
plt.rcParams['axes.labelweight'] = 'bold'
plt.rcParams["font.size"] = 25
plt.rcParams["font.weight"] = "bold"
plt.rcParams["axes.labelweight"] = "bold"
plt.rcParams['axes.titlepad'] = 20
matplotlib.rcParams.update({'legend.fontsize': 20})
matplotlib.rcParams.update({'font.weight': 'bold'})
matplotlib.rcParams.update({'axes.linewidth': 4})

rawData = [(100,107.0,163.0,109.0),(200,407.0,643.0,433.0),(300,850.0,1256.,891.0),(400,1433.,2072.,1486.),(500,2156.,3303.,2274.)]

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
  comms[0].append(440*Dim2/1e8)
  comms[1].append(1320*Dim2/1e8)
  comms[2].append((440*Dim2+1600)/1e8)

plt.title(benchmarkName, fontsize=25, fontweight='bold')
for i in range(3):
  plt.plot(sizes,comms[i],label=names[i],linestyle=linestyles[i],color=colors[i],linewidth=3)
# plt.legend(loc='upper left')
plt.xlabel('Input Size')
plt.ylabel('Data (×100MB)')
plt.xticks(range(100,600,100), ['100²', '200²', '300²', '400²', '500²'])
plt.yticks(range(4), ['0', '1', '2', '3'])
plt.tight_layout()
plt.savefig('comms-{}-paper.png'.format(benchmarkName), pad_inches = 0)
plt.close('all')
