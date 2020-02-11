import math
import random

# <0.01,0.01,0.95R(0.001>=d(x),0.001>=d(r))>
def approxP2C(x,r):
  x2 = x*x

  # sinx = ((((-1/7!)x2+(1/5!))x2+(-1/3!))x2+1)x
  # cosx = ((((1/8!)x2+(-1/6!))x2+(1/4!))x2+(-1/2!))x2+1
  sinxApprox = (((((1.0/362880.0)*x2+(-1.0/5040.0))*x2+(1.0/120.0))*x2+(-1.0/6.0))*x2+1.0)*x
  cosxApprox = ((((1.0/40320.0)*x2+(-1.0/720.0))*x2+(1.0/24.0))*x2+(-1.0/2.0))*x2+1.0

  return (sinxApprox*r,cosxApprox*r)

# <0.01,0.01,R(0.001>=d(x),0.001>=d(r))>
def exactP2C(x,r):
  return (math.sin(x)*r,math.cos(x)*r)

random.seed()
histogram = {x:0 for x in range(10)}

# guarantee
# 0.005+d(x)+d(r)>=d(sx), 0.005+d(x)+d(r)>=d(cx)
guaranteeFails = 0

for i in range(10000000):
  x = random.uniform(-math.pi,math.pi)
  xd = random.uniform(-0.001,0.001)
  r = random.random()
  rd = random.uniform(-0.001,0.001)
  sinxApprox,cosxApprox = approxP2C(x+xd,r+rd)
  sinx,cosx = exactP2C(x,r)
  sinxDel = abs(sinx-sinxApprox)
  cosxDel = abs(cosx-cosxApprox)
  maxDel = max(sinxDel,cosxDel)
  for j in range(10):
    if maxDel < 10**(-j):
      histogram[j] += 1
  if 0.005+xd+rd<maxDel:
    guaranteeFails += 1
  elif maxDel>0.01:
    print('BAD')

print(histogram)
print(guaranteeFails)

# post
# 1<=R(0.01>=d(sx),0.01>=d(cx))
# pre recover
# 1<=R(0.001>=d(x),0.001>=d(r))
# pre try unchecked
# 1<=0.95R(0.001>=d(x),0.001>=d(r))
# pre try checked
# 1<=R(0.01>=0.005+d(x)+d(r))
# pre combined unchecked
# 1<=R(0.001>=d(x),0.001>=d(r)) AND 1<=0.95R(0.001>=d(x),0.001>=d(r))
# pre combined checked
# 1<=R(0.001>=d(x),0.001>=d(r)) AND 1<=R(0.01>=0.005+d(x)+d(r))
