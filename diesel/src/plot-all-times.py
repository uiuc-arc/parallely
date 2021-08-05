import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import numpy as np
from matplotlib.lines import Line2D

def geoMean(iterable):
  arr = np.array(iterable)
  return arr.prod()**(1./len(arr))

benchmarks = ['PageRank',
              'SSSP',
              'BFS',
              #'Gaussian',
              'SOR',
              'Sobel',
              'MatrixMult',
              'KMeans',
              'Regression',
              #'Hiring',
              #'Income SVM',
              #'Income DT',
              #'Income NN',
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
        'Hiring': [(1e6,None,1.505),(2e6,None,2.248),(4e6,None,7.021),(8e6,None,1.710)],
        'Income SVM': [(1e6,None,0.582),(2e6,None,3.221),(4e6,None,5.746),(8e6,None,1.067)],
        'Income DT': [(1e6,None,1.164),(2e6,None,1.935),(4e6,None,2.957),(8e6,None,0.560)],
        'Income NN': [(1e6,None,2.856),(2e6,None,3.177),(4e6,None,1.154),(8e6,None,2.064)],
       }
        

names = ['baseline', 'Diamont']
linestyles = ['--', '-']
colors = ['red', 'green', 'blue']
markers = ['o', 'v', '^', 's', 'p', 'P', '*', 'X', 'D', 'd', '1', '2']

geomeanData = [None,([],[]),([],[]),([],[]),([],[]),([],[]),([],[]),([],[]),([],[]),([],[])]
geomeanRange = 9

plt.figure(figsize=(11,6))
# plt.title('Input Size vs. Overhead', fontsize=20)
for i, benchmark in enumerate(benchmarks):
  benchmarkData = data[benchmark]
  sizes = [datum[0] for datum in benchmarkData]
  relSizes = [size/sizes[0] for size in sizes]
  baseTimes = [datum[1] for datum in benchmarkData]
  dieselTimes = [datum[2] for datum in benchmarkData]
  hasBaseTimes = (baseTimes[0] != None)
  # interpolate geomean data
  for j in range(1, geomeanRange):
    interpBaseTime = None
    interpDieselTime = None
    for k, relSize in enumerate(relSizes):
      if relSize == j:
        # found exact
        if hasBaseTimes:
          interpBaseTime = baseTimes[k]
        interpDieselTime = dieselTimes[k]
        break
      if relSize > j or k==len(relSizes)-1:
        # need to interpolate with previous
        scaleFactor = 1./(relSizes[k]-relSizes[k-1])
        lowerFactor = (j-relSizes[k-1])*scaleFactor
        upperFactor = (relSizes[k]-j)*scaleFactor
        if hasBaseTimes:
          interpBaseTime = baseTimes[k-1]*lowerFactor + baseTimes[k]*upperFactor
        interpDieselTime = dieselTimes[k-1]*lowerFactor + dieselTimes[k]*upperFactor
        break
    if interpBaseTime:
      geomeanData[j][0].append(max(interpBaseTime+1.,1.))
    if interpDieselTime:
      geomeanData[j][1].append(max(interpDieselTime+1.,1.))
  dieselColor = colors[1]
  if hasBaseTimes:
    plt.plot(relSizes, baseTimes, label=benchmark+'-'+names[0], linestyle=linestyles[0], color=colors[0], marker=markers[i], markersize=10)
  else:
    dieselColor = colors[2]
  plt.plot(relSizes, dieselTimes, label=benchmark+'-'+names[1], linestyle=linestyles[1], color=dieselColor, marker=markers[i], markersize=10)

baseGeomeans = []
dieselGeomeans = []
for j in range(1,geomeanRange):
  baseGeomeans.append(geoMean(geomeanData[j][0])-1.)
  dieselGeomeans.append(geoMean(geomeanData[j][1])-1.)
print(baseGeomeans)
print(dieselGeomeans)
# plt.plot(range(1,geomeanRange), baseGeomeans, color='blue')
# plt.plot(range(1,geomeanRange), dieselGeomeans, color='yellow')

plt.plot([0,9],[25,25],linestyle=':',color='black')

plt.xlabel('Relative Input Size', fontsize=18)
plt.ylabel('Overhead%', fontsize=18)
plt.xticks(fontsize=18)
plt.yticks(fontsize=18)
plt.xlim((1.95,8.2))
plt.ylim((0, 205))

legend_line_elements = [Line2D([0], [0], color='r', linestyle="--", lw=4, label='Baseline'),
                        Line2D([0], [0], color='g', lw=4, label='Diamont')]
legend_point_elements = [Line2D([0], [0], marker=markers[i], label=benchmarks[i],
                                markerfacecolor='black', markersize=15, lw=0) for i in range(len(benchmarks))]
legend = legend_line_elements + legend_point_elements

legend1 = plt.legend(handles=legend, fontsize=15, loc='upper center', ncol=5)
# legend2 = plt.legend(handles= legend2_elements, fontsize=16, bbox_to_anchor=(1, 1), loc='upper right', ncol=4)
# plt.gca().add_artist(legend2)
# plt.gca().add_artist(legend1)
# plt.legend(handles= legend_elements, fontsize=16,bbox_to_anchor=(0.99,-0.07),loc='lower left')
plt.tight_layout()
# plt.subplots_adjust(right=0.1)

plt.savefig('times-all.png')
plt.close('all')
