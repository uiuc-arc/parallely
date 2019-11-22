import os

benchmarkName = 'SVM_Income'

cwd = os.getcwd()
data = dict()
data2 = dict()
for dirname in ["uninstrumented","unoptimized"]:
   os.chdir(dirname)

   os.system("python getTime.py > out.txt")
   file = open("out.txt")
   lines = file.readlines()
   line = lines[0].replace(" ","").replace("ns","").replace("\n","")
   data["20000"]=str(line)

   
   os.system("sed -i 's/20000/50000/g' " + dirname + ".go")
   os.system("go build")
   os.system("python getTime.py > out.txt")
   file = open("out.txt")
   lines = file.readlines()
   line = lines[0].replace(" ","").replace("ns","").replace("\n","")
   data["50000"]=str(line)


   os.system("sed -i 's/50000/100000/g' " + dirname + ".go")
   os.system("go build")
   os.system("python getTime.py > out.txt")
   file = open("out.txt")
   lines = file.readlines()
   line = lines[0].replace(" ","").replace("ns","").replace("\n","")
   data["100000"]=str(line)
   
   os.system("sed -i 's/100000/150000/g' " + dirname + ".go")
   os.system("go build")
   os.system("python getTime.py > out.txt")
   file = open("out.txt")
   lines = file.readlines()
   line = lines[0].replace(" ","").replace("ns","").replace("\n","")
   data["150000"]=str(line)

   if dirname == "uninstrumented":
      data2["20000"] = data["20000"]
      data2["50000"] = data["50000"]
      data2["100000"] = data["100000"]
      data2["150000"] = data["150000"]
      
   print data
   #reset back
   os.system("sed -i 's/150000/20000/g' " + dirname + ".go")
   os.system("go build")
   print("done with " + dirname)
   os.chdir(cwd)


import matplotlib.pyplot as plt

data1 = data

data = {20000: (float(data2["20000"]), float(data1["20000"])),
        50000: (float(data2["50000"]), float(data1["50000"])),
        100000: (float(data2["100000"]), float(data1["100000"])),
        150000: (float(data2["150000"]), float(data1["150000"]))}

names = ['Baseline', 'Unoptimized']
linestyles = ['-', '--']
colors = ['orange', 'red' ]

sizes = []
times = ([], [], [])
comms = ([], [], [])

for datum in sorted(data.keys()):
  sizes.append(datum)
  times[0].append(data[datum][0])
  times[1].append(data[datum][1])
  #times[2].append(data[datum][2]/1e9)
  '''
  comms[0].append(data[datum][3]/1e6)
  comms[1].append(data[datum][4]/1e6)
  comms[2].append(data[datum][5]/1e6)
  '''
plt.figure(figsize=(2,2))
plt.rcParams.update({'font.size': 12})
plt.ticklabel_format(style='sci', axis='y', scilimits=(0,0))
for i in range(2):
  plt.plot(sizes,times[i],label=names[i],linestyle=linestyles[i],color=colors[i])
# plt.legend(loc='upper left')
plt.xlabel('Input Size')
plt.ylabel('Time (ms)')
plt.tight_layout()
plt.savefig('times-{}.png'.format(benchmarkName))
plt.close()
