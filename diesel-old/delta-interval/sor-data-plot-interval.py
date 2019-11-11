import math
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

graph_lower = []
graph_upper = []
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

plt.plot(range(1,101),graph_lower,label='Lower Bound')
plt.plot(range(1,101),graph_num,label='Calculated Value')
plt.plot(range(1,101),graph_upper,label='Upper Bound')
plt.xlabel('Iteration')
plt.ylabel('Relative Value (×10^'+str(multiplier_power)+')')
plt.legend()
plt.tight_layout()
plt.show()
