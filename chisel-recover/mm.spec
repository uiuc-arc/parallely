1.0 <= R(0.01 >= d(mr))

m1 in [0.0, 1.0]
m2 in [0.0, 1.0]
mr in [0.0, 1.0]

approxMult in < 2, 1, 1.00, 0.000001, [0.0, 1.0], 0.000001, [0.0, 1.0], 0.00001,  [0.0, 1.0] >
exactMult  in < 2, 1, 1.00, 0.000001, [0.0, 1.0], 0.000001, [0.0, 1.0], 0.000002, [0.0, 1.0] >

chk ensures < 4*m1e+4*m2e >= prod >