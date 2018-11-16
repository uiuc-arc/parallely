#use python 3

import random
import sys

numNodes = int(sys.argv[1])
numEdges = int(sys.argv[2])

random.seed()

edges = random.sample(range(numNodes*numNodes), numEdges)

edges.sort()

for edge in edges:
  node1 = edge//numNodes
  node2 = edge%numNodes
  print(str(node1)+'\t'+str(node2))
