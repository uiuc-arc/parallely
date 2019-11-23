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

names = ['Unoptimized','Optimized']
linestyles = ['--',':']
colors = ['red','green']

data = [
[
2.383020305529726e-07,
1.9040851260143656e-07,
1.4091436542074124e-07,
1.2906154996628063e-07,
1.12960607958712e-07,
1.0544567546702414e-07,
9.776101402277961e-08,
9.29408845761998e-08,
8.848317971356439e-08,
8.522242969802463e-08,
8.231030273328729e-08,
],[
2.383020305529726e-07,
2.383110961901025e-07,
2.383110961901025e-07,
2.383110961901025e-07,
2.383110961901025e-07,
2.383110961901025e-07,
2.383110961901025e-07,
2.383110961901025e-07,
2.383110961901025e-07,
2.383110961901025e-07,
2.383110961901025e-07,
]]

plt.title(benchmarkName, fontsize=25, fontweight='bold')
for i in range(2):
  plt.plot(range(11),[x*1e7 for x in data[i]],label=names[i],linestyle=linestyles[i],color=colors[i],linewidth=3)
plt.xlabel("Iteration")
plt.ylabel(r"Max Error ($\times 10^{-7}$)")
plt.ylim(0.0,2.5)

plt.tight_layout()
plt.savefig('error-{}.png'.format(benchmarkName), pad_inches = 0)
plt.close('all')
