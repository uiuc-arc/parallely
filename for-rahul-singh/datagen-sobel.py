#! /usr/bin/env python3

import sys
import random

# 1st argument - dimension of array
arrayDim = int(sys.argv[1])
# 2nd argument - output file
outFileName = sys.argv[2]

arraySize = arrayDim**2

random.seed()
outFile = open(outFileName, 'w')

print('Generating square 2D array of dimension',arrayDim,'and writing to',outFileName)

for i in range(arraySize):
    outFile.write(str(random.randint(0,255))+'\n')

outFile.close()

print('Done')
