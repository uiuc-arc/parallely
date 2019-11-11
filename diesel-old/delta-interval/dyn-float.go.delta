package main

import "math"
import "sort"

func AddFloatInterval(x, y, d1, d2 float64) float64 {
  return d1+d2
}

func SubFloatInterval(x, y, d1, d2 float64) float64 {
  return d1+d2
}

func MulFloatInterval(x, y, d1, d2 float64) float64 {
  xlow := x-d1
  xhig := x+d1
  ylow := y-d2
  yhig := y+d2
  limits := []float64{0,0,0,0}
  limits[0] = xlow*ylow
  limits[1] = xlow*yhig
  limits[2] = xhig*ylow
  limits[3] = xhig*yhig
  sort.Float64s(limits)
  return math.Max(limits[3]-x*y,x*y-limits[0])
}

func DivFloatInterval(x, y, d1, d2 float64) float64 {
  ylow := y-d2
  yhig := y+d2
  inf := math.Inf(0)
  if (ylow<0 && yhig>0) || (ylow==0 && yhig==0) {
    ylow = -inf; yhig = inf
  } else if ylow==0 {
    ylow = 1.0/yhig; yhig = inf
  } else if yhig==0 {
    yhig = 1.0/ylow; ylow = -inf
  } else {
    temp := 1.0/yhig; yhig = 1.0/ylow; ylow = temp
  }
  return MulFloatInterval(x, y, d1, math.Max(yhig-1.0/y,1.0/y-ylow))
}

func AddRoundingError(x, d, epsilon float64) float64 {
  xlow := x-d
  xhig := x+d
  xlow -= math.Abs(xlow)*epsilon
  xhig += math.Abs(xhig)*epsilon
  return math.Max(xhig-x,x-xlow)
}

type DynFloat32 struct {
  num float32
  delta float64
}

const Float32Epsilon = 1.0/16777216.0

func MakeDynFloat32(x float32) (ret DynFloat32) {
  ret.num = x
  ret.delta = 0
  return
}

func AddDynFloat32(x, y DynFloat32) (ret DynFloat32) {
  ret.num = x.num + y.num
  ret.delta = AddFloatInterval(float64(x.num), float64(y.num), x.delta, y.delta)
  ret.delta = AddRoundingError(float64(ret.num), ret.delta, Float32Epsilon)
  return
}

func SubDynFloat32(x, y DynFloat32) (ret DynFloat32) {
  ret.num = x.num - y.num
  ret.delta = SubFloatInterval(float64(x.num), float64(y.num), x.delta, y.delta)
  ret.delta = AddRoundingError(float64(ret.num), ret.delta, Float32Epsilon)
  return
}

func MulDynFloat32(x, y DynFloat32) (ret DynFloat32) {
  ret.num = x.num * y.num
  ret.delta = MulFloatInterval(float64(x.num), float64(y.num), x.delta, y.delta)
  ret.delta = AddRoundingError(float64(ret.num), ret.delta, Float32Epsilon)
  return
}

func DivDynFloat32(x, y DynFloat32) (ret DynFloat32) {
  ret.num = x.num / y.num
  ret.delta = DivFloatInterval(float64(x.num), float64(y.num), x.delta, y.delta)
  ret.delta = AddRoundingError(float64(ret.num), ret.delta, Float32Epsilon)
  return
}

type DynFloat64 struct {
  num float64
  delta float64
}

const Float64Epsilon = 1.0/9007199254740992.0

func MakeDynFloat64(x float64) (ret DynFloat64) {
  ret.num = x
  ret.delta = 0
  return
}

func AddDynFloat64(x, y DynFloat64) (ret DynFloat64) {
  ret.num = x.num + y.num
  ret.delta = AddFloatInterval(x.num, y.num, x.delta, y.delta)
  ret.delta = AddRoundingError(ret.num, ret.delta, Float64Epsilon)
  return
}

func SubDynFloat64(x, y DynFloat64) (ret DynFloat64) {
  ret.num = x.num - y.num
  ret.delta = SubFloatInterval(x.num, y.num, x.delta, y.delta)
  ret.delta = AddRoundingError(ret.num, ret.delta, Float64Epsilon)
  return
}

func MulDynFloat64(x, y DynFloat64) (ret DynFloat64) {
  ret.num = x.num * y.num
  ret.delta = MulFloatInterval(x.num, y.num, x.delta, y.delta)
  ret.delta = AddRoundingError(ret.num, ret.delta, Float64Epsilon)
  return
}

func DivDynFloat64(x, y DynFloat64) (ret DynFloat64) {
  ret.num = x.num / y.num
  ret.delta = DivFloatInterval(x.num, y.num, x.delta, y.delta)
  ret.delta = AddRoundingError(ret.num, ret.delta, Float64Epsilon)
  return
}

func DynFloat32To64(x DynFloat32) (ret DynFloat64) {
  ret.num = float64(x.num)
  ret.delta = x.delta
  return
}

func DynFloat64To32(x DynFloat64) (ret DynFloat32) {
  var temp32 float32
  var temp64 float64
  ret.num = float32(x.num)
  temp32 = float32(x.num-x.delta); temp64 = float64(temp32)
  xlow := math.Min(temp64,x.num-x.delta)
  temp32 = float32(x.num+x.delta); temp64 = float64(temp32)
  xhig := math.Min(temp64,x.num+x.delta)
  ret.delta = xhig-xlow
  return
}

