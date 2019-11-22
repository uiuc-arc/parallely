import os

cwd = os.getcwd()
data = dict()

for dirname in ["uninstrumented","unoptimized"]:
   os.chdir(dirname)

   os.system("python getTime.py > out.txt")
   file = open("out.txt")
   lines = file.readlines()
   line = lines[0].replace(" ","").replace("ms","").replace("\n","")
   data["20000"]=str(line)

   
   os.system("sed -i 's/20000/50000/g' " + dirname + ".go")
   os.system("go build")
   os.system("python getTime.py > out.txt")
   file = open("out.txt")
   lines = file.readlines()
   line = lines[0].replace(" ","").replace("ms","").replace("\n","")
   data["50000"]=str(line)


   os.system("sed -i 's/50000/100000/g' " + dirname + ".go")
   os.system("go build")
   os.system("python getTime.py > out.txt")
   file = open("out.txt")
   lines = file.readlines()
   line = lines[0].replace(" ","").replace("ms","").replace("\n","")
   data["100000"]=str(line)
   
   os.system("sed -i 's/100000/150000/g' " + dirname + ".go")
   os.system("go build")
   os.system("python getTime.py > out.txt")
   file = open("out.txt")
   lines = file.readlines()
   line = lines[0].replace(" ","").replace("ms","").replace("\n","")
   data["150000"]=str(line)


   
   print data
   #reset back
   os.system("sed -i 's/150000/20000/g' " + dirname + ".go")
   os.system("go build")
   print("done with " + dirname)
   os.chdir(cwd)
