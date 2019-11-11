import math
import matplotlib
import matplotlib.pyplot as plt

multiplier_power = 6
multiplier = math.pow(10,multiplier_power)

lower = []
upper = []
num = []
dataFile = open('sor-data-interval.txt')
for line in dataFile:
  raw = line.strip().split()
  lower.append(float(raw[1]))
  num.append(float(raw[2]))
  upper.append(float(raw[3]))
dataFile.close()

delta = []
dataFile = open('sor-data-delta.txt')
for line in dataFile:
  raw = line.strip().split()
  delta.append(float(raw[2]))
dataFile.close()

graph_lower = []
graph_upper = []
graph_lower2 = []
graph_upper2 = []
graph_num = []

#relative to bounds
#for i in range(100):
#  size = upper[i]-lower[i]
#  graph_lower.append(-size/2.0*multiplier)
#  graph_upper.append(size/2.0*multiplier)
#  pos = (num[i]-lower[i])-(size/2.0)
#  graph_num.append(pos*multiplier)

#relative to value
for i in range(100):
  graph_lower.append((lower[i]-num[i])*multiplier)
  graph_upper.append((upper[i]-num[i])*multiplier)
  graph_num.append(0)

for i in range(100):
  graph_lower2.append(graph_num[i]-delta[i]*multiplier)
  graph_upper2.append(graph_num[i]+delta[i]*multiplier)

font = {'size'   : 14}

matplotlib.rc('font', **font)

plt.plot(range(1,101),graph_lower2,label='Lower Bound (Delta)')
plt.plot(range(1,101),graph_lower,label='Lower Bound (Interval)')
plt.plot(range(1,101),graph_num,label='Calculated Value')
plt.plot(range(1,101),graph_upper,label='Upper Bound (Interval)')
plt.plot(range(1,101),graph_upper2,label='Upper Bound (Delta)')
plt.xlabel('Iteration')
plt.ylabel('Relative Value (×10^'+str(multiplier_power)+')')
plt.legend()
plt.tight_layout()
plt.show()
