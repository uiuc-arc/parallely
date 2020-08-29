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

	
	var vel float64


	var acc float64



	var vf float64

	//DynMap[2].PrintValues()

	var dist float64
	//DynMap[3] = diesel.NewMultiInterval(vf,interval)

	//read from sensor
	vel = getVelocity()



	//read from sensor
	acc = getAcceleration()


	//set the concrete value
	if (vel > 10.0) {
		acc = acc-5.0
	} else {
		acc = acc+5.0
	}

	//read from sensor
	vf = getFinalVelocity()


	dist = (vf*vf)-(vel*vel)/acc
	_ = dist
	/*if (dist < 100 && dist >-100){
		
	}*/
	//var check = diesel.CheckFloat64MultiIntervalLessThan(DynMap[3],100.0,0.01)

	//var check2 = diesel.CheckFloat64MultiIntervalGreaterThan(DynMap[3],-100.0,0.01)
	//fmt.Println(check)
	//fmt.Println(check2)
	//fmt.Println(check,check2)
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
