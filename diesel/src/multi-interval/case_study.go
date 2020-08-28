package main

import (
	"../diesel"
	"fmt"
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
	return 1.0
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
	var interval diesel.ProbInterval = diesel.ProbInterval{1, 0}
	DynMap[0] = diesel.NewMultiInterval(vel,interval)

	var acc float64
	DynMap[1] = diesel.NewMultiInterval(acc,interval)


	var vf float64
	DynMap[2] = diesel.NewMultiInterval(vf,interval)


	var dist float64
	//DynMap[3] = diesel.NewMultiInterval(vf,interval)

	//read from sensor
	vel = getVelocity()
	DynMap[0].SetValueI(vel,0)

	//read from sensor
	acc = getAcceleration()
	DynMap[1].SetValueI(acc,0)

	var min_vel,max_vel float64
	max_vel,min_vel = diesel.MinMaxMultiInterval(DynMap[0])

	//slow down
	if (min_vel > 10.0) {
		acc = acc - 5.0
		DynMap[1].SetValueI(acc,0)	//set dynamic tracking for this value
	} else if (max_vel < 10.0) {
		acc = acc + 5.0
		DynMap[1].SetValueI(acc+5.0,0) //set dynamic tracking for this value
	} else {
		DynMap[1].SetValueI(acc-5.0,0) //set dynamic tracking for both values
		DynMap[1].SetValueI(acc+5.0,1) //set dynamic tracking for both values
		//set the concrete value
		if (vel > 10.0) {
			acc = acc - 5.0
		} else {
			acc = acc + 5.0
		}
	}


	//read from sensor
	vf = getVelocity()
	DynMap[2].SetValueI(vf,0)

	
	DynMap[3] = diesel.Divide(diesel.Subtract(diesel.Multiply(DynMap[2],DynMap[2]),diesel.Multiply(DynMap[0],DynMap[0])),DynMap[1])
	dist = DynMap[3].GetValueI(0)
	_ = dist


	var check = diesel.CheckFloat64MultiIntervalLessThan(DynMap[3],10.0,0.1)
	fmt.Println(check)


	fmt.Println("Ending thread : ", 0)
}

func main() {
	fmt.Println("Starting main thread")
	diesel.InitChannels(1)

	startTime := time.Now()

	go func_0()

	fmt.Println("Main thread waiting for others to finish")
	diesel.Wg.Wait()

	elapsed := time.Since(startTime)

	fmt.Println("Done!")
	fmt.Println("Elapsed time : ", elapsed.Nanoseconds())
}
