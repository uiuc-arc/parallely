import os

times = []
for i in range(500):
    os.system("./unoptimized > time.txt")
    f = open("time.txt")
    lines = f.readlines()
    time = lines[0].replace("ms","").replace("\n","")
    time = float(time)
    times.append(time)

print str(sum(times)/float(len(times))) + " ms"
    
