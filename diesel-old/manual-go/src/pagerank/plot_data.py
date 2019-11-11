import matplotlib
# matplotlib.use("Agg")

import json
import numpy as np
import matplotlib.pyplot as plt
from matplotlib.ticker import StrMethodFormatter

plt.gca().yaxis.set_major_formatter(StrMethodFormatter('{x:,.5f}'))

# ax.yaxis.set_major_formatter(FormatStrFormatter('%.2f'))

with open("result_file_aggregated.txt", "r") as outfile:
    real_rel = json.load(outfile)
with open("result_file_aggregated_tracked.txt", "r") as outfile:
    tracked_rel = json.load(outfile)
with open("result_file_aggregated_decaf.txt", "r") as outfile:
    decaf_rel = json.load(outfile)

# print real_rel
temp_dict_real = {}
temp_dict_tracked = {}
temp_dict_decaf = {}

for key in real_rel:
    temp_dict_real[int(key)] = real_rel[key]
    # print tracked_rel[key]
    temp_dict_tracked[int(key)] = tracked_rel[key]
    temp_dict_decaf[int(key)] = decaf_rel[key]

# print temp_dict_decaf

data_x = sorted([int(x) for x in real_rel.keys()])

labels = {0: "min", 1: "avg", 2: "max"}
# print temp_dict_tracked[0]

data_y_real = [float(temp_dict_real[x]) / 100 for x in data_x]

# data_y_tracked_mean = [temp_dict_tracked[x][0][1] for x in data_x]
# data_y_tracked_min = [temp_dict_tracked[x][0][0] for x in data_x]
# data_y_tracked_max = [temp_dict_tracked[x][0][2] for x in data_x]

# print data_y_tracked_mean

data_y_tracked_mean = [np.mean([data[1] for data in temp_dict_tracked[x]]) for x in data_x]
data_y_tracked_min = [np.min([data[0] for data in temp_dict_tracked[x]]) for x in data_x]
data_y_tracked_max = [np.max([data[2] for data in temp_dict_tracked[x]]) for x in data_x]

# plt.plot(data_x, data_y_tracked_mean, label="tracked")
plt.plot(data_x, data_y_tracked_min, label="tracked-min")
# plt.plot(data_x, data_y_tracked_max, label="tracked-max")

# plt.savefig("tracking.png")
# plt.close("all")

data_y_decaf_mean = [np.mean([data[1] for data in temp_dict_decaf[x]]) for x in data_x]
data_y_decaf_min = [np.min([data[0] for data in temp_dict_decaf[x]]) for x in data_x]
data_y_decaf_max = [np.max([data[2] for data in temp_dict_decaf[x]]) for x in data_x]

# data_y_decaf_mean = [temp_dict_decaf[x][0][1] for x in data_x]
# data_y_decaf_min = [temp_dict_decaf[x][0][0] for x in data_x]
# data_y_decaf_max = [temp_dict_decaf[x][0][2] for x in data_x]

# plt.plot(data_x, data_y_decaf_mean, label="decaf")
plt.plot(data_x, data_y_decaf_min, label="decaf-min")
# plt.plot(data_x, data_y_decaf_max, label="decaf-max")

# plt.savefig("decaf.png")
# plt.close("all")

# plt.plot(data_x, data_y_real, label="real-{}".format(labels[i]))
# plt.plot(data_x, data_y_tracked, label="tracked-{}".format(labels[i]))
# plt.plot(data_x, data_y_decaf, label="decaf-{}".format(labels[i]))

# # data_y_tracked = [np.mean([data[i] for data in [temp_dict_tracked[x] for x in data_x]])]

# for i in labels:
#     data_y_tracked = [np.mean([data[i] for data in [temp_dict_tracked[x] for x in data_x]])]
#     data_y_decaf = [np.mean([data[i] for data in [temp_dict_decaf[x] for x in data_x]])]

#     print data_y_tracked

#     plt.plot(data_x, data_y_real, label="real-{}".format(labels[i]))
#     plt.plot(data_x, data_y_tracked, label="tracked-{}".format(labels[i]))
#     plt.plpot(data_x, data_y_decaf, label="decaf-{}".format(labels[i]))

plt.legend()

plt.show()
# plt.savefig("tracking.png")
