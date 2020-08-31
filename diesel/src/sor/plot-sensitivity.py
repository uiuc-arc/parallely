import matplotlib
matplotlib.use('Agg')
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
2.9791682565338817e-08,
2.4539767129461156e-08,
1.9990545214731985e-08,
1.7813827790647000e-08,
1.6929082730909877e-08,
1.6168948033539193e-08,
1.5645981299713622e-08,
1.5197750728243837e-08,
1.4796327994284757e-08,
1.4467376113626572e-08,
1.4160833686261824e-08,
],[
2.9791682565338817e-08,
2.9801069278967418e-08,
2.9801069278967418e-08,
2.9801069278967418e-08,
2.9801069278967418e-08,
2.9801069278967418e-08,
2.9801069278967418e-08,
2.9801069278967418e-08,
2.9801069278967418e-08,
2.9801069278967418e-08,
2.9801069278967418e-08,
]]

plt.title(benchmarkName, fontsize=25, fontweight='bold')
# plt.plot([0,10],[4,4],label='Bound',linestyle='-.',color='blue',linewidth=5)
for i in range(2):
  plt.plot(range(11),[x*1e8 for x in data[i]],label=names[i],linestyle=linestyles[i],color=colors[i],linewidth=3)
plt.xlabel("Iteration")
plt.ylabel(r"Max Error ($\times 10^{-8}$)")
plt.ylim(0.0,4.2)
plt.xticks(range(11),[str(s) for s in range(11)])

plt.tight_layout()
plt.savefig('error-{}.png'.format(benchmarkName), pad_inches = 0)
plt.close('all')
