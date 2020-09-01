import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import numpy as np

def geoMean(iterable):
  arr = np.array(iterable)
  return arr.prod()**(1./len(arr))

benchmarks = ['PageRank',
                  'SSSP',
                  'BFS',
                  'Gaussian',
                  'SOR',
                  'Sobel',
                  'MatrixMult',
                  'KMeans',
                  'Regression',
                 ]

data = {
        'PageRank': [(8144, 39.942990838888562, 5.3038781010606133), (22687, 66.090542885120001, 5.4688009551664862), (36682, 66.785405031401893, 5.4745750546201215), (62586, 79.358979772256305, 10.146616359932214)],
        'SSSP': [(8114, 56.340225592889226, 6.7141361380467366), (22687, 95.725894751313064, 13.202103651880105), (36682, 113.68491543272962, 13.766150928003498), (62586, 137.07275651828729, 15.028044330073831)],
        'BFS': [(8144, 53.31178085942382, 5.5986196191271986), (22687, 88.859752152819766, 6.3612704582469686), (36682, 111.94461646504392, 13.523106623786212), (62586, 128.31724585597374, 14.656620596938938)],
        'Gaussian': [(147456, 2.2367055205358635, -0.9879276090274567), (262144, 4.756731508769932, 2.5616135069020665), (1048576, 3.2908452469406595, 0.7618934680070565)], # (65536, 0.4095003540466425, -1.1665634223183645), 
        'SOR': [(10000, 69.92116158573907, 3.9544192221090695), (40000, 110.42792950580116, 5.883356799603509), (90000, 124.95085391220152, 8.683904918742101)],#, (160000, 128.71759019120674, 8.03079398214932)],
        'Sobel': [(10000,92.9,13.1),(14400,103.,13.7),(19600,114.,15.0),(25600,117.,18.6),(32400,121.,19.7),(90000,126.,21.0)],
        'MatrixMult': [(10000,108.,18.5),(14400,110.,29.5),(19600,132.,30.5),(25600,130.,38.0),(32400,166.,65.6),(40000,168.,73.0),(90000,167.,73.5)],
        'KMeans': [(248, 37.97132382950256, 2.1114444212934282), (512, 38.08782047671841, 0.3284102709640509), (1024, 41.634931039260096, 1.6938769212795561), (2048, 46.490315535838185, 4.990946215861019)],
        'Regression': [(500, 37.677171760711595, 14.765906836545136), (1000, 35.07539131142273, 10.79295400994414), (1500, 36.644736922643325, 13.001743127020996), (2000, 44.577089225723384, 16.70492627710252), (2500, 46.97774337946356, 18.59035151775253), (3000, 48.001511642987026, 19.94895202234672), (3500, 49.336999492105434, 22.190399334523093), (4000, 63.59011048405577, 31.34188253822615)],
       }
        

names = ['Baseline', 'Diesel']
linestyles = ['--', '-']
colors = ['red', 'green']
markers = ['o', 'v', '^', 's', 'p', 'P', '*', 'X', 'D']

geomeanData = [None,([],[]),([],[]),([],[]),([],[]),([],[]),([],[]),([],[]),([],[]),([],[])]

plt.figure(figsize=(6,6))
# plt.title('Input Size vs. Overhead', fontsize=20)
for i, benchmark in enumerate(benchmarks):
  benchmarkData = data[benchmark]
  sizes = [datum[0] for datum in benchmarkData]
  relSizes = [size/sizes[0] for size in sizes]
  baseTimes = [datum[1] for datum in benchmarkData]
  dieselTimes = [datum[2] for datum in benchmarkData]
  # interpolate geomean data
  for j in range(1,10):
    interpBaseTime = None
    interpDieselTime = None
    for k, relSize in enumerate(relSizes):
      if relSize == j:
        # found exact
        interpBaseTime = baseTimes[k]
        interpDieselTime = dieselTimes[k]
        break
      if relSize > j:
        # need to interpolate with previous
        scaleFactor = 1./(relSizes[k]-relSizes[k-1])
        lowerFactor = (j-relSizes[k-1])*scaleFactor
        upperFactor = (relSizes[k]-j)*scaleFactor
        interpBaseTime = baseTimes[k-1]*lowerFactor + baseTimes[k]*upperFactor
        interpDieselTime = dieselTimes[k-1]*lowerFactor + dieselTimes[k]*upperFactor
        break
    if interpBaseTime and interpDieselTime:
      geomeanData[j][0].append(max(interpBaseTime+1.,1.))
      geomeanData[j][1].append(max(interpDieselTime+1.,1.))
  plt.plot(relSizes, baseTimes, label=benchmark+'-base', linestyle=linestyles[0], color=colors[0], marker=markers[i], markersize=10)
  plt.plot(relSizes, dieselTimes, label=benchmark+'-Diesel', linestyle=linestyles[1], color=colors[1], marker=markers[i], markersize=10)

baseGeomeans = []
dieselGeomeans = []
for j in range(1,10):
  baseGeomeans.append(geoMean(geomeanData[j][0])-1.)
  dieselGeomeans.append(geoMean(geomeanData[j][1])-1.)
plt.plot(range(1,10), baseGeomeans, color='blue')
plt.plot(range(1,10), dieselGeomeans, color='yellow')

plt.xlabel('Relative Input Size', fontsize=18)
plt.ylabel('Overhead%', fontsize=18)
plt.xticks(fontsize=18)#, rotation=90)
plt.yticks(fontsize=18)
plt.legend(fontsize=15,bbox_to_anchor=(0.99,-0.12),loc='lower left')
plt.tight_layout()
plt.savefig('times-all.png',bbox_inches='tight')
plt.close('all')
