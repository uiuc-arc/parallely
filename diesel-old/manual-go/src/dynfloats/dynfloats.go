package dynfloats

import "math"
//import "sort"
import "math/rand"

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

type DynRelyInt struct {
	Value int
	Reliability float32
	// ReliabilityDecaf float32
}

type DynRelyFloat struct {
	Value float64
	Reliability float32
	// ReliabilityDecaf float32
}

func GreaterIntRely(x, y DynRelyInt) (ret bool) {
  return (x.Value > y.Value)
}

func EqualconstIntRely(x DynRelyInt, y int) (ret bool) {
  return (x.Value == y)
}

func AddConstRely(x DynRelyInt, y int) (ret DynRelyInt) {
  ret.Value = x.Value + y
	// ret.ReliabilityDecaf = x.ReliabilityDecaf
  ret.Reliability = x.Reliability
  return
}

func MultConstRelyIF(x DynRelyInt, y float64) (ret DynRelyInt) {
  ret.Value = int(float64(x.Value) * y)
	// ret.ReliabilityDecaf = x.ReliabilityDecaf
  ret.Reliability = x.Reliability
  return
}

func AddIntRely(x, y DynRelyInt) (ret DynRelyInt) {
  ret.Value = x.Value + y.Value
	// ret.ReliabilityDecaf = x.ReliabilityDecaf * y.ReliabilityDecaf
  ret.Reliability = max(x.Reliability + y.Reliability - 1, 0)
  return
}

func MinusIntRely(x, y DynRelyInt) (ret DynRelyInt) {
  ret.Value = x.Value - y.Value
	// ret.ReliabilityDecaf = x.ReliabilityDecaf * y.ReliabilityDecaf
  ret.Reliability = max(x.Reliability + y.Reliability - 1, 0)
  return
}

func MultIntRely(x, y DynRelyInt) (ret DynRelyInt) {
  ret.Value = x.Value * y.Value
	// ret.ReliabilityDecaf = x.ReliabilityDecaf * y.ReliabilityDecaf
  ret.Reliability = max(x.Reliability + y.Reliability - 1, 0)
  return
}

func DivIntRely(x, y DynRelyInt) (ret DynRelyInt) {
  ret.Value = x.Value * y.Value
	// ret.ReliabilityDecaf = x.ReliabilityDecaf * y.ReliabilityDecaf
  ret.Reliability = max(x.Reliability + y.Reliability - 1, 0)
  return
}

func DynSendDynInt(signal chan bool, regular chan DynRelyInt, value DynRelyInt, failrate float32) (failed bool) {
	failure := rand.Float32()

	sendVal := DynRelyInt{value.Value, value.Reliability * (1 -
		failrate)}

	if failure >= failrate {		
		signal <- true
		regular <- sendVal

		return true
	} else {
		signal <- false
		return false
	}
	// regular <- value.Valpue
	// dynamic <- value.Reliability
}

func DynSendDynIntArray(signal chan bool, regular chan []DynRelyInt, value []DynRelyInt, failrate float32) (failed bool) {
	failure := rand.Float32()
		
	for i := range(value) {
		value[i].Reliability = value[i].Reliability * (1 - failrate)
	}

	if failure >= failrate {
		// fmt.Println("Pass")
		signal <- true
		regular <- value
		return true
	} else {
		// fmt.Println("Fail")
		signal <- false
		return false
	}
}

func DynSendIntArray(signal chan bool, regular chan []int, value []int, failrate float32) (failed bool) {
	failure := rand.Float32()
		
	// for i := range(value) {
	// 	value[i].Reliability = value[i].Reliability * (1 - failrate)
	// }

	if failure >= failrate {
		// fmt.Println("Pass")
		signal <- true
		regular <- value
		return true
	} else {
		// fmt.Println("Fail")
		signal <- false
		return false
	}
}

func DynSendInt(signal chan bool, regular chan int, value int, failrate float32) (failed bool) {
	failure := rand.Float32()

	sendVal := value

	if failure >= failrate {		
		signal <- true
		regular <- sendVal

		return true
	} else {
		signal <- false
		return false
	}
	// regular <- value.Valpue
	// dynamic <- value.Reliability
}

func AddFloatRely(x, y DynRelyFloat) (ret DynRelyFloat) {
  ret.Value = x.Value + y.Value
	// ret.ReliabilityDecaf = x.ReliabilityDecaf * y.ReliabilityDecaf
  ret.Reliability = max(x.Reliability + y.Reliability - 1, 0)
  return
}

func MinusFloatRely(x, y DynRelyFloat) (ret DynRelyFloat) {
  ret.Value = x.Value - y.Value
	// ret.ReliabilityDecaf = x.ReliabilityDecaf * y.ReliabilityDecaf
  ret.Reliability = max(x.Reliability + y.Reliability - 1, 0)
  return
}

func MultFloatRely(x, y DynRelyFloat) (ret DynRelyFloat) {
  ret.Value = x.Value * y.Value
	// ret.ReliabilityDecaf = x.ReliabilityDecaf * y.ReliabilityDecaf
  ret.Reliability = max(x.Reliability + y.Reliability - 1, 0)
  return
}

func DivFloatRely(x, y DynRelyFloat) (ret DynRelyFloat) {
  ret.Value = x.Value / y.Value
	// ret.ReliabilityDecaf = x.ReliabilityDecaf * y.ReliabilityDecaf
  ret.Reliability = max(x.Reliability + y.Reliability - 1, 0)
  return
}


func DynSendDynFloat(signal chan bool, regular chan DynRelyFloat, value DynRelyFloat, failrate float32) (failed bool) {
	failure := rand.Float32()

	// sendval := DynRelyFloat{value.Value,
	// 	value.Reliability * (1 - failrate), value.ReliabilityDecaf * (1 - failrate)}

	sendval := DynRelyFloat{value.Value, value.Reliability * (1 - failrate)}
	
	// value.Reliability = value.Reliability * (1 - failrate)
	// value.ReliabilityDecaf = value.ReliabilityDecaf * (1 - failrate)

	if failure >= failrate {		
		signal <- true
		regular <- sendval
		return true
	} else {
		signal <- false
		return false
	}
	// regular <- value.Valpue
	// dynamic <- value.Reliability
}

func DynSendFloat(signal chan bool, regular chan float64, value float64, failrate float32) (failed bool) {
	failure := rand.Float32()

	sendval := value

	if failure >= failrate {		
		signal <- true
		regular <- sendval
		return true
	} else {
		signal <- false
		return false
	}
}

func AddFloatInterval(x, y, d1, d2 float64) float64 {
  return d1+d2
}

func SubFloatInterval(x, y, d1, d2 float64) float64 {
  return d1+d2
}

func MulFloatInterval(x, y, d1, d2 float64) float64 {
/*
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
*/
  return math.Abs(y*d1) + math.Abs(x*d2) + d1*d2
}

func DivFloatInterval(x, y, d1, d2 float64) float64 {
/*
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
*/
  if y-d2 < 0 && y+d2 > 0 {
    return math.Inf(0)
  }
  absy := math.Abs(y)
  return (absy*d1 + math.Abs(x)*d2)/absy/(absy-d2)
/*
  temp := (math.Abs(y*d1) + math.Abs(x*d2))/y
  if y > 0 {
    return temp / (y - d2)
  } else {
    return temp / (y + d2)
  }
*/
}

func AddRoundingError(x, d, epsilon float64) float64 {
  xlow := x-d
  xhig := x+d
  xlow -= math.Abs(xlow)*epsilon
  xhig += math.Abs(xhig)*epsilon
  return math.Max(xhig-x,x-xlow)
}

type DynFloat32 struct {
  Num float32
  Delta float64
//  Ops uint64
}

const Float32Epsilon = 1.0/16777216.0

func MakeDynFloat32(x float32) (ret DynFloat32) {
  ret.Num = x
  ret.Delta = 0
//  ret.Ops = 0
  return
}

func AddDynFloat32(x, y DynFloat32) (ret DynFloat32) {
  ret.Num = x.Num + y.Num
  ret.Delta = AddFloatInterval(float64(x.Num), float64(y.Num), x.Delta, y.Delta)
//  ret.Delta = AddRoundingError(float64(ret.Num), ret.Delta, Float32Epsilon)
//  ret.Ops = x.Ops + y.Ops + 1
  return
}

func SubDynFloat32(x, y DynFloat32) (ret DynFloat32) {
  ret.Num = x.Num - y.Num
  ret.Delta = SubFloatInterval(float64(x.Num), float64(y.Num), x.Delta, y.Delta)
//  ret.Delta = AddRoundingError(float64(ret.Num), ret.Delta, Float32Epsilon)
//  ret.Ops = x.Ops + y.Ops + 1
  return
}

func MulDynFloat32(x, y DynFloat32) (ret DynFloat32) {
  ret.Num = x.Num * y.Num
  ret.Delta = MulFloatInterval(float64(x.Num), float64(y.Num), x.Delta, y.Delta)
//  ret.Delta = AddRoundingError(float64(ret.Num), ret.Delta, Float32Epsilon)
//  ret.Ops = x.Ops + y.Ops + 1
  return
}

func DivDynFloat32(x, y DynFloat32) (ret DynFloat32) {
  ret.Num = x.Num / y.Num
  ret.Delta = DivFloatInterval(float64(x.Num), float64(y.Num), x.Delta, y.Delta)
//  ret.Delta = AddRoundingError(float64(ret.Num), ret.Delta, Float32Epsilon)
//  ret.Ops = x.Ops + y.Ops + 1
  return
}

func SqrtDynFloat32(x DynFloat32) (ret DynFloat32) {
  val := math.Sqrt(float64(x.Num))
  max := math.Sqrt(float64(x.Num)+x.Delta)
  min := math.Sqrt(float64(x.Num)-x.Delta)
  ret.Delta = math.Max(max-val,val-min)
//  ret.Delta = AddRoundingError(val, ret.Delta, Float32Epsilon)
  ret.Num = float32(val)
//  ret.Ops = x.Ops + 1
  return
}

func ExpDynFloat32(x DynFloat32) (ret DynFloat32) {
  val := math.Exp(float64(x.Num))
  max := math.Exp(float64(x.Num)+x.Delta)
  min := math.Exp(float64(x.Num)-x.Delta)
  ret.Delta = math.Max(max-val,val-min)
//  ret.Delta = AddRoundingError(val, ret.Delta, Float32Epsilon)
  ret.Num = float32(val)
//  ret.Ops = x.Ops + 1
  return
}

func LogDynFloat32(x DynFloat32) (ret DynFloat32) {
  val := math.Log(float64(x.Num))
  max := math.Log(float64(x.Num)+x.Delta)
  min := math.Log(float64(x.Num)-x.Delta)
  ret.Delta = math.Max(max-val,val-min)
//  ret.Delta = AddRoundingError(val, ret.Delta, Float32Epsilon)
  ret.Num = float32(val)
//  ret.Ops = x.Ops + 1
  return
}

type DynFloat64 struct {
  Num float64
  Delta float64
//  Ops uint64
}

const Float64Epsilon = 1.0/9007199254740992.0

func MakeDynFloat64(x float64) (ret DynFloat64) {
  ret.Num = x
  ret.Delta = 0
//  ret.Ops = 0
  return
}

func AddDynFloat64(x, y DynFloat64) (ret DynFloat64) {
  ret.Num = x.Num + y.Num
  ret.Delta = AddFloatInterval(x.Num, y.Num, x.Delta, y.Delta)
//  ret.Delta = AddRoundingError(ret.Num, ret.Delta, Float64Epsilon)
//  ret.Ops = x.Ops + y.Ops + 1
  return
}

func SubDynFloat64(x, y DynFloat64) (ret DynFloat64) {
  ret.Num = x.Num - y.Num
  ret.Delta = SubFloatInterval(x.Num, y.Num, x.Delta, y.Delta)
//  ret.Delta = AddRoundingError(ret.Num, ret.Delta, Float64Epsilon)
//  ret.Ops = x.Ops + y.Ops + 1
  return
}

func MulDynFloat64(x, y DynFloat64) (ret DynFloat64) {
  ret.Num = x.Num * y.Num
  ret.Delta = MulFloatInterval(x.Num, y.Num, x.Delta, y.Delta)
//  ret.Delta = AddRoundingError(ret.Num, ret.Delta, Float64Epsilon)
//  ret.Ops = x.Ops + y.Ops + 1
  return
}

func DivDynFloat64(x, y DynFloat64) (ret DynFloat64) {
  ret.Num = x.Num / y.Num
  ret.Delta = DivFloatInterval(x.Num, y.Num, x.Delta, y.Delta)
//  ret.Delta = AddRoundingError(ret.Num, ret.Delta, Float64Epsilon)
//  ret.Ops = x.Ops + y.Ops + 1
  return
}

func DynFloat32To64(x DynFloat32) (ret DynFloat64) {
  ret.Num = float64(x.Num)
  ret.Delta = x.Delta
  return
}

func DynFloat64To32(x DynFloat64) (ret DynFloat32) {
  var temp32 float32
  var temp64 float64
  ret.Num = float32(x.Num)
  temp32 = float32(x.Num-x.Delta); temp64 = float64(temp32)
  xlow := math.Min(temp64,x.Num-x.Delta)
  temp32 = float32(x.Num+x.Delta); temp64 = float64(temp32)
  xhig := math.Min(temp64,x.Num+x.Delta)
  ret.Delta = xhig-xlow
//  ret.Ops = x.Ops + 1
  return
}

