import subprocess
import re
import numpy as np

commstr = """./bfs ../../../benchmarks/inputs/p2p-Gnutella31.txt 62586 out-exact.txt"""

times = []

for i in range(10):
    print "Running Iteration : ", i
    result_test = subprocess.check_output(commstr, shell=True)

    matches = re.findall("Elapsed time : .*\n", result_test)
    time_spent = float(matches[0].split(' : ')[-1])
    print time_spent
    times.append(time_spent)


print np.mean(times)
