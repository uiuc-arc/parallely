package dynfloat_fairness

import "math"


type DynFairnessFloat struct {
  Val float64
  Epsilon float64
  Delta float64
}



func AddFloatFairness(x,y DynFairnessFloat)(ret DynFairnessFloat){
     ret.Val = x.Val + y.Val
     ret.Epsilon = x.Epsilon + y.Epsilon
     ret.Delta = x.Delta + y.Delta
     return
}

func NegFloatFairness(x DynFairnessFloat)(ret DynFairnessFloat){
     ret.Val = -x.Val
     ret.Epsilon = x.Epsilon
     ret.Delta = x.Delta
     return
}


func MulFloatFairness(x,y DynFairnessFloat)(ret DynFairnessFloat){
     ret.Val = x.Val * y.Val
     ret.Epsilon = (math.Abs(ret.Val)* y.Epsilon)+(math.Abs(y.Val)*x.Epsilon)+(x.Epsilon*y.Epsilon)
     ret.Delta = x.Delta + y.Delta
     return
}


func ConstMulFloatFairness(c float64,x DynFairnessFloat)(ret DynFairnessFloat){
	C := DynFairnessFloat{val: c,epsilon:0;delta:0}
	ret = MulFloatFairness(C,x)
	return
}

func InvFloatFairness(x DynFairnessFloat)(ret DynFairnessFloat){
     if ((x.Val-x.Epsilon<0)&&(x.Val+x.Epsilon>0)){
     	ret.Epsilon = math.Inf(1) //possible division by zero
		ret.Val = x.Val
		ret.Delta = x.Delta
     } else {
     	ret.Val = 1/(x.Val)
		ret.Epsilon = x.Epsilon/(math.Abs(x.Val)*(math.Abs(x.Val)-x.Epsilon))
		ret.Delta = x.Delta
     }
     return
}

func DivFloatFairness(x,y DynFairnessFloat)(ret DynFairnessFloat){
     y_inv := InvFloatFairness(y)
     ret = MulFloatFairness(x,y_inv)
     return
}

func checkIneq(x DynFairnessFloat, eps,del float64)(ret bool){
     ret = (x.Epsilon <= eps && x.Delta <= del) 
     return
}



