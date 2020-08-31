import matplotlib
matplotlib.use('Agg')

import matplotlib.pyplot as plt
import numpy as np

benchmarkName = 'BFS'

data = {8144: ([49430168.0, 48174500.0, 39782502.0, 42309669.0, 37452947.0, 39373950.0, 43499912.0, 46344150.0, 42096921.0, 42881427.0, 55013325.0, 41907622.0, 44601670.0, 41903039.0, 47331856.0, 46037848.0, 41184590.0, 42514926.0, 44671551.0, 46984456.0], [74929472.0, 64194076.0, 64078090.0, 67196516.0, 79541812.0, 62434830.0, 69451461.0, 66096713.0, 64811359.0, 68548254.0, 66357142.0, 71292696.0, 61578885.0, 64440976.0, 72565265.0, 66227527.0, 69618357.0, 64440627.0, 65323998.0, 71376973.0], [48617475.0, 44148723.0, 46766673.0, 41947475.0, 54139481.0, 42246354.0, 40457428.0, 48128491.0, 45794602.0, 40519698.0, 46691846.0, 65359204.0, 48151849.0, 49114580.0, 47565793.0, 49868548.0, 40303347.0, 44761404.0, 44311102.0, 44066594.0]), 22687: ([72406344.0, 71795031.0, 72183427.0, 75413490.0, 85002102.0, 83007928.0, 65370604.0, 67992465.0, 79696400.0, 70765858.0, 76087904.0, 78735321.0, 77455274.0, 83534917.0, 77113288.0, 80025549.0, 73282588.0, 83159347.0, 64693237.0, 84415025.0], [161844663.0, 144653870.0, 132030246.0, 140951907.0, 141649075.0, 147876860.0, 127531557.0, 134760058.0, 152173359.0, 137916522.0, 151786869.0, 140571159.0, 149967874.0, 142552395.0, 146976594.0, 150427771.0, 131799982.0, 155097531.0, 139542584.0, 144591588.0], [88850475.0, 74733017.0, 75227290.0, 69554813.0, 78969443.0, 77265040.0, 80523767.0, 83281856.0, 81606263.0, 86988009.0, 83728282.0, 90747668.0, 75245883.0, 84459861.0, 90411661.0, 81212713.0, 81771606.0, 71797187.0, 81855304.0, 80733155.0]), 36682: ([95736989.0, 95846929.0, 103977175.0, 94981735.0, 110399498.0, 91411316.0, 97565222.0, 108007954.0, 91664186.0, 96212651.0, 108416875.0, 86799931.0, 102646159.0, 97366294.0, 95027263.0, 96605705.0, 109211877.0, 93416532.0, 93199285.0, 90017019.0], [210603298.0, 199680616.0, 211357578.0, 206776329.0, 202515003.0, 193894307.0, 217207786.0, 214861480.0, 210606711.0, 200875569.0, 199957562.0, 205024734.0, 203309937.0, 195840394.0, 203889674.0, 221615115.0, 221377354.0, 215579961.0, 213712988.0, 202271373.0], [116993436.0, 103729867.0, 118821204.0, 111171548.0, 109011006.0, 120252876.0, 124691779.0, 108286251.0, 109298965.0, 109163899.0, 127939166.0, 112554887.0, 100270474.0, 95369690.0, 103272645.0, 111781827.0, 116299195.0, 109461085.0, 98730560.0, 116261711.0]), 62586: ([145861300.0, 148340992.0, 154798570.0, 141514967.0, 147009730.0, 152937034.0, 140555826.0, 161732485.0, 150567602.0, 133818882.0, 153032468.0, 142827917.0, 140568822.0, 118683937.0, 137872655.0, 142967421.0, 145697289.0, 130168963.0, 131679279.0, 127917054.0], [326642125.0, 321854492.0, 319836203.0, 357347166.0, 321090222.0, 319274837.0, 320354290.0, 324415427.0, 315174301.0, 329630067.0, 339347562.0, 329965807.0, 324218384.0, 321638282.0, 325436815.0, 337679560.0, 315602304.0, 314776576.0, 320547566.0, 318906211.0], [175904643.0, 146535274.0, 176731985.0, 173847956.0, 154522998.0, 173376373.0, 153768588.0, 162781134.0, 159109726.0, 165678576.0, 161643006.0, 166298670.0, 162055133.0, 146460424.0, 153785176.0, 158455503.0, 168602609.0, 177563700.0, 162192143.0, 166741210.0])}

names = ['Unoptimized', 'Optimized']
linestyles = ['--', ':']
colors = ['red', 'green']

sizes = []
times = ([], [])
comms = ([], [])

for datum in sorted(data.keys()):
  sizes.append(datum)
  baseline = np.mean(data[datum][0])/1e9
  # times[0].append(baseline/baseline)
  times[0].append(((np.mean(data[datum][1])/1e9)-baseline)/baseline*100)
  times[1].append(((np.mean(data[datum][2])/1e9)-baseline)/baseline*100)

# plt.figure(figsize=(2,2))
plt.title(benchmarkName, fontsize=20)
# plt.ticklabel_format(style='sci', axis='y', scilimits=(0,0))
for i in range(2):
  plt.plot(sizes, times[i],label=names[i],linestyle=linestyles[i],color=colors[i])
plt.legend(loc='upper left')
plt.xlabel('Input Size', fontsize=20)
plt.ylabel('Overhead%', fontsize=20)
plt.xticks(fontsize=20, rotation=90)
plt.yticks(fontsize=20)
plt.tight_layout()
plt.savefig('times-{}.png'.format(benchmarkName))
plt.close('all')

# # plt.figure(figsize=(2,2))
# plt.title(benchmarkName, fontsize=20)
# # plt.ticklabel_format(style='sci', axis='y', scilimits=(0,0))
# for i in range(3):
#   plt.plot(sizes,comms[i],label=names[i],linestyle=linestyles[i],color=colors[i])
# # plt.legend(loc='upper left')
# plt.xlabel('Input Size', fontsize=20)
# plt.ylabel('Data (MB)', fontsize=20)
# plt.xticks(fontsize=20, rotation=90)
# plt.yticks(fontsize=20)
# plt.tight_layout()
# plt.savefig('comms-{}.png'.format(benchmarkName))
# plt.close('all')
