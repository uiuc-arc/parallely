#! /usr/bin/env python3

import sys
import random

# 1st argument - size of array
arraySize = int(sys.argv[1])
# 2nd argument - output file
outFileName = sys.argv[2]

random.seed()
outFile = open(outFileName, 'w')

print('Generating array of size',arraySize,'and writing to',outFileName)

for i in range(arraySize):
    outFile.write(str(random.randint(-2**31,2**31-1))+'\n')

outFile.close()

print('Done')
