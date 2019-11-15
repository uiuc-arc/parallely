package dynfloats

import "math"


type DynFairnessFloat struct {
  val float64
  epsilon float64
  delta float64
}


func AddFloatFairness(x,y DynFairnessFloat)(ret DynFairnessFloat){
     ret.val = x.val + y.val
     ret.epsilon = x.epsilon + y.epsilon
     ret.delta = x.delta + y.delta
     return
}

func NegFloatFairness(x DynFairnessFloat)(ret DynFairnessFloat){
     ret.val = -x.val
     ret.epsilon = x.epsilon
     ret.delta = x.delta
     return
}


func MulFloatFairness(x,y DynFairnessFloat)(ret DynFairnessFloat){
     ret.val = x.val * y.val
     ret.epsilon = (math.Abs(ret.val)* y.epsilon)+(math.Abs(y.val)*x.epsilon)+(x.epsilon*y.epsilon)
     ret.delta = x.delta + y.delta
     return
}

func InvFloatFairness(x DynFairnessFloat)(ret DynFairnessFloat){
     if ((x.val-x.epsilon<0)&&(x.val+x.epsilon>0)){
     	ret.epsilon = math.Inf(1) //possible division by zero
		ret.val = x.val
		ret.delta = x.delta
     } else {
     	ret.val = 1/(x.val)
		ret.epsilon = x.epsilon/(math.Abs(x.val)*(math.Abs(x.val)-x.epsilon))
		ret.delta = x.delta
     }
     return
}

func DivFloatFairness(x,y DynFairnessFloat)(ret DynFairnessFloat){
     y_inv := InvFloatFairness(y)
     ret = MulFloatFairness(x,y_inv)
     return
}

func checkIneq(x DynFairnessFloat, eps,del float64)(ret bool){
     ret = (x.epsilon <= eps && x.delta <= del) 
	 return
}



