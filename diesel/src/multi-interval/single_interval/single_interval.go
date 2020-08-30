package main

import (
	"../../diesel"
	"fmt"
	"math/rand"
	"time"
)

var Num_threads int

var Num_nodes int
var Num_edges int
var NodesPerThread int



func convertToFloat(x int) float64 {
	return float64(x)
}



func getVelocity() float64{
	return 10
	return 8+(5*rand.Float64())
	//return 1.0
}


func getFinalVelocity() float64{
	return 20
	return 20+(3*rand.Float64())
	//return 1.0
}



func getAcceleration() float64{
	return 1.0
}



func func_0() {
	defer diesel.Wg.Done()
	//var DynMap [10]diesel.ProbInterval
	//Declare variables
	var DynMap [4]diesel.ProbInterval
	
	var vel float64
	DynMap[0] = diesel.ProbInterval{0, 0}

	var acc float64
	DynMap[1] = diesel.ProbInterval{0, 0}


	var vf float64
	DynMap[2] = diesel.ProbInterval{0, 0}
	//DynMap[2].PrintValues()

	var dist float64
	//DynMap[3] = diesel.NewMultiInterval(vf,interval)

	//read from sensor
	vel = getVelocity()
	DynMap[0].Reliability = 0.0001
	DynMap[0].Delta = 0.8



	//read from sensor
	acc = getAcceleration()
	DynMap[1].Reliability = 0.0001
	DynMap[1].Delta = 0.8


	//slow down
	if (vel - DynMap[0].Delta > 10.0) {
		acc = acc - 5.0


	} else if (vel + DynMap[0].Delta < 10.0) {
		acc = acc + 5.0



	} else {
		var t1,t2 float64
		t1 = acc-5.0
		t2 = acc+5.0
		var diff = (t2+DynMap[1].Delta)-(t1-DynMap[1].Delta)
		DynMap[1].Delta = diff

		//set the concrete value
		if (vel > 10.0) {
			acc = t1
		} else {
			acc = t2
		}
	}


	//read from sensor
	vf = getFinalVelocity()
	DynMap[2].Reliability  = 0.0001 

	var tmp1 float64
	var int1 diesel.ProbInterval 
	tmp1, int1 = diesel.MulProbInterval(vel,vel,DynMap[2],DynMap[2])

	var tmp2 float64
	var int2 diesel.ProbInterval
	tmp2, int2 = diesel.MulProbInterval(vf,vf,DynMap[0],DynMap[0]) 


	var tmp3 float64
	var int3 diesel.ProbInterval
	tmp3, int3 = diesel.SubProbInterval(tmp2,tmp1,int2,int1)

	//var tmp4 float64
	var int4 diesel.ProbInterval
	dist, int4 = diesel.DivProbInterval(tmp3,acc,int3,DynMap[1])
	DynMap[3] = int4

	//var check = diesel.CheckFloat64MultiIntervalLessThan(DynMap[3],100.0,0.01)
	var check = diesel.CheckFloat64(dist,int4,0.001,100-dist)
	var check2 = diesel.CheckFloat64(dist,int4,0.001,100-dist)
	//var check2 = diesel.CheckFloat64MultiIntervalGreaterThan(DynMap[3],-100.0,0.01)
	//fmt.Println(check)
	//fmt.Println(check2)
	fmt.Println(check,check2)
	fmt.Println("Ending thread : ", 0)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	fmt.Println("Starting main thread")
	diesel.InitChannels(1)

	startTime := time.Now()

	//go func_0()
	func_0()

	fmt.Println("Main thread waiting for others to finish")
	diesel.Wg.Wait()

	elapsed := time.Since(startTime)

	fmt.Println("Done!")
	fmt.Println("Elapsed time : ", elapsed.Nanoseconds())
}
