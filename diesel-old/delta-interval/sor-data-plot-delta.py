import math
import matplotlib.pyplot as plt

multiplier_power = 6
multiplier = math.pow(10,multiplier_power)

delta = []
num = []
dataFile = open('sor-data-delta.txt')
for line in dataFile:
  raw = line.strip().split()
  num.append(float(raw[1]))
  delta.append(float(raw[2]))
dataFile.close()

graph_lower = []
graph_upper = []
graph_num = []

for i in range(100):
  graph_lower.append(-delta[i]*multiplier)
  graph_upper.append(delta[i]*multiplier)
  graph_num.append(0)

plt.plot(range(1,101),graph_lower,label='Lower Bound')
plt.plot(range(1,101),graph_num,label='Calculated Value')
plt.plot(range(1,101),graph_upper,label='Upper Bound')
plt.xlabel('Iteration')
plt.ylabel('Relative Value (×10^'+str(multiplier_power)+')')
plt.legend()
plt.tight_layout()
plt.show()
