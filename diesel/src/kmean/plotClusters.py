import matplotlib.pyplot as plt

centers = []

datafile = open("output.txt", "r")
data = [float(i) for i in datafile.read()[1:-2].split(' ')]

centerfile = open("centers.txt", "r")
center = [float(i) for i in centerfile.read()[1:-2].split(' ')]

icenterfile = open("init-centers.txt", "r")
icenter = [float(i) for i in icenterfile.read()[1:-2].split(' ')]

xs = []
ys = []
for i in range(len(data)/2):
    xs.append(data[2*i])
    ys.append(data[2*i+1])
 
cxs = []
cys = []
for i in range(len(center)/2):
    cxs.append(center[2*i])
    cys.append(center[2*i+1])

icxs = []
icys = []
for i in range(len(icenter)/2):
    icxs.append(icenter[2*i])
    icys.append(icenter[2*i+1])           

plt.scatter(xs, ys, color='grey', alpha=0.3)
plt.scatter(cxs, cys, color='r')
plt.scatter(icxs, icys, color='g')
plt.tight_layout()
plt.savefig('clusters.png')
plt.close('all')
    
print len(data)
