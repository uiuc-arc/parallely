import subprocess
import numpy as np
import json

commstr = """./sssp ../../../benchmarks/inputs/p2p-Gnutella31.txt 62586 out-exact.txt"""

real_reliability = {}
tracked_reliability = {}
tracked_reliability_decaf = {}
for j in range(20):
    real_reliability[j] = 0
    tracked_reliability[j] = []
    tracked_reliability_decaf[j] = []

biggest_diff = 0
biggest_id = 0

for i in range(20):
    print "Running Iteration : ", i
    result_test = subprocess.check_output(commstr, shell=True)

    for j in range(20):
        calc_rel = []
        calc_rel_decaf = []
        with open("_iter_{}.txt".format(j)) as result_file:
            # print result_file.readline()[:-1]
            # print result_file.readline()[:-1], (result_file.readline()[:-1][0] == "t")
            if result_file.readline()[:-1][0] != "t":
                real_reliability[j] += 1

            for line in result_file.readlines():
                # print line[1:-2], line[1:-2].split(' ')
                _, tracked, decaf = line[1:-2].split(' ')
                calc_rel.append(float(tracked))
                calc_rel_decaf.append(float(decaf))

                # if float(tracked)-biggest_diff

        tracked_reliability[j].append([np.min(calc_rel), np.mean(calc_rel), np.max(calc_rel)])
        tracked_reliability_decaf[j].append([np.min(calc_rel_decaf), np.mean(calc_rel_decaf),
                                             np.max(calc_rel_decaf)])

        # tracked_reliability[j].append([calc_rel[0], calc_rel[0], calc_rel[0]])
        # tracked_reliability_decaf[j].append([calc_rel_decaf[0], calc_rel_decaf[1], calc_rel_decaf[2]])

with open("result_file_aggregated.txt", "w") as outfile:
    json.dump(real_reliability, outfile)
with open("result_file_aggregated_tracked.txt", "w") as outfile:
    json.dump(tracked_reliability, outfile)
with open("result_file_aggregated_decaf.txt", "w") as outfile:
    json.dump(tracked_reliability_decaf, outfile)
