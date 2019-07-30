import numpy as np

error_sums = []

for i in range(10):
    precisefile = open("__output__precise_{}.txt".format(i), "r")
    approxfile = open("__output__approx_{}.txt".format(i), "r")

    precise = [float(k) for k in precisefile.readlines()[0][1:-1].split(' ')[:10]]
    approx = [float(k) for k in approxfile.readlines()[0][1:-1].split(' ')[:10]]

    sum_sd = 0
    for i in range(len(precise)):
        sum_sd += ((precise[i] - approx[i]) * (precise[i] - approx[i]))

    error_sums.append(sum_sd / float(len(precise)))
    print error_sums

print np.mean(error_sums)
