import subprocess
import re
import numpy as np

ArrayDimList = [100, 200, 300, 400]

results = {}

for ArrayDim in ArrayDimList:
    ArraySz = ArrayDim*ArrayDim
    SliceSz = ArraySz/10

    template_str = open("sor.tmpl", 'r').readlines()
    with open("sor.par", "w") as fout:
        for line in template_str:
            newline = line.replace('__ARRAYDIM__', str(ArrayDim))
            newline = newline.replace('__ARRAYSZ__', str(ArraySz))
            newline = newline.replace('__SLICESZ__', str(SliceSz))
            fout.write(newline)

    template_str = open("__basic_go_main.tmpl", 'r').readlines()
    with open("__basic_go_main.txt", "w") as fout:
        for line in template_str:
            newline = line.replace('__ARRAYDIM__', str(ArrayDim))
            newline = newline.replace('__ARRAYSZ__', str(ArraySz))
            newline = newline.replace('__SLICESZ__', str(SliceSz))
            fout.write(newline)

    template_str = open("__basic_go_worker.tmpl", 'r').readlines()
    with open("__basic_go_worker.txt", "w") as fout:
        for line in template_str:
            newline = line.replace('__ARRAYDIM__', str(ArrayDim))
            newline = newline.replace('__ARRAYSZ__', str(ArraySz))
            newline = newline.replace('__SLICESZ__', str(SliceSz))
            fout.write(newline)



    commstr = """python ../../../parser/crosscompiler-diesel-dist-acc.py -f sor.par -tm __basic_go_main.txt -tw __basic_go_worker.txt -o out.go -i"""
    result_test = subprocess.check_output(commstr, shell=True)
    print(result_test)

    result_test = subprocess.check_output("python getSpeedUp.py", shell=True)
    matches = re.findall("Overhead : .*\n", result_test)
    orig_overhead = float(matches[0].split(' : ')[-1]) / 1000000

    matches = re.findall("Overhead After Optimization : .*\n", result_test)
    opt_overhead = float(matches[0].split(' : ')[-1]) / 1000000

    results[ArrayDim] = (orig_overhead, opt_overhead)

    print(results)
