import random
import sys

output = open(sys.argv[1],"w")

random.seed()
for i in range(int(sys.argv[2])):
  line = str(random.randint(1, 1000))
  output.write(line+"\n")

output.close()
