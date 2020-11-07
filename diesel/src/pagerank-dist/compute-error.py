import numpy as np

orig_out = open("output.txt", "r").readlines()
orig_prs = []

for line in orig_out:
    orig_prs.append(float(line))

tmp = np.array(orig_prs)
orig_top_10 = tmp.argsort()[-10:][::-1]
print(orig_top_10)

max_error = []
avg_error = []
count_error = []

for i in range(100):
    new_prs = []
    new_out = open("outputs/output-{}.txt".format(i), "r")
    j = 0
    errors = []
    for line in new_out:
        new_prs.append(float(line))
        diff = abs(float(line) - orig_prs[j])
        if diff > 0:
            errors.append(diff)
        j += 1
    if len(errors) > 0:
        tmp2 = np.array(new_prs)
        new_top_10 = tmp2.argsort()[-10:][::-1]
        print(orig_top_10)

        
        max_error.append(max(errors))
        avg_error.append(sum(errors)/len(errors))
        count_error.append(len(errors))

print(max_error)
print(avg_error)
print(count_error)
