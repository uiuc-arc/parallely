package main

import (
	"../diesel"
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
	var DynMap [4]diesel.DynMultiInterval
	
	var vel float64
	var interval diesel.ProbInterval = diesel.ProbInterval{0, 0}
	DynMap[0] = diesel.NewMultiInterval(vel,interval)
	//DynMap[0].SetProbIntervalI(interval,0)
	//DynMap[0].PrintValues()

	var acc float64
	DynMap[1] = diesel.NewMultiInterval(acc,interval)
	//DynMap[1].PrintValues()

	var vf float64
	DynMap[2] = diesel.NewMultiInterval(vf,interval)
	//DynMap[2].PrintValues()

	var dist float64
	//DynMap[3] = diesel.NewMultiInterval(vf,interval)

	//read from sensor
	vel = getVelocity()
	DynMap[0].SetValueI(vel,0)
	DynMap[0].SetReliabilityI(0.0001,0) 
	DynMap[0].SetDeltaI(0.8,0) //important!
	//DynMap[0].PrintValues()

	//read from sensor
	acc = getAcceleration()
	DynMap[1].SetValueI(acc,0)
	DynMap[1].SetReliabilityI(0.0001,0) 
	DynMap[1].SetDeltaI(0.05,0)


	var min_vel,max_vel, t1, t2 float64
	min_vel,max_vel = diesel.MinMaxMultiInterval(DynMap[0])


	//slow down
	if (min_vel > 10.0) {
		acc = acc - 5.0
		DynMap[1].SetValueI(acc,0)	//set dynamic tracking for this value
	} else if (max_vel < 10.0) {
		acc = acc + 5.0
		DynMap[1].SetValueI(acc+5.0,0) //set dynamic tracking for this value
	} else {
		t1 = acc-5.0
		t2 = acc+5.0
		DynMap[1].SetValueI(t1,0) //set dynamic tracking for both values
		
		DynMap[1].AddOneMore(t2,DynMap[1].GetIntervalI(0))
		//DynMap[1].SetValueI(acc+5.0,1) //set dynamic tracking for both values

		//set the concrete value
		if (vel > 10.0) {
			acc = t1
		} else {
			acc = t2
		}
	}
	//fmt.Println("acceleration:")
	//DynMap[1].PrintValues()

	//read from sensor
	vf = getFinalVelocity()
	DynMap[2].SetValueI(vf,0)
	DynMap[2].SetReliabilityI(0.0001,0) 
	DynMap[2].SetDeltaI(0.05,0)

	
	DynMap[3] = diesel.Divide(diesel.Subtract(diesel.Multiply(DynMap[2],DynMap[2]),diesel.Multiply(DynMap[0],DynMap[0])),DynMap[1])
	//DynMap[3].PrintValues()
	dist = DynMap[3].GetValueI(0)
	_ = dist


	var check = diesel.CheckFloat64MultiIntervalLessThan(DynMap[3],100.0,0.01)
	var check2 = diesel.CheckFloat64MultiIntervalGreaterThan(DynMap[3],-100.0,0.01)
	//fmt.Println(check)
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
