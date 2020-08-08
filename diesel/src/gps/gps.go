package main

import (
	"diesel"
	"fmt"
	"math"
	"time"
)

var Num_threads int
var Edges [811400]int
var Inlinks [8114]int
var Outlinks [8114]int
var DistGlobal [8114]int
var Num_nodes int
var Num_edges int
var NodesPerThread int

var Out float64

func max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func convertToFloat(x int) float64 {
	return float64(x)
}

func getGPSX() float64 {
	return math.Sqrt(10)
}

func getGPSY() float64 {
	return 1.0
}

func func_0() {
	defer diesel.Wg.Done()
	var DynMap [10]diesel.ProbInterval
	var my_chan_index int
	_ = my_chan_index
	_ = DynMap
	var x1 float64
	DynMap[0] = diesel.ProbInterval{1, 0}
	var y1 float64
	DynMap[1] = diesel.ProbInterval{1, 0}
	var x2 float64
	DynMap[2] = diesel.ProbInterval{1, 0}
	var y2 float64
	DynMap[3] = diesel.ProbInterval{1, 0}
	var t1 float64
	DynMap[4] = diesel.ProbInterval{1, 0}
	var t2 float64
	DynMap[5] = diesel.ProbInterval{1, 0}
	var t3 float64
	DynMap[6] = diesel.ProbInterval{1, 0}
	var t4 float64
	DynMap[7] = diesel.ProbInterval{1, 0}
	var t5 float64
	DynMap[8] = diesel.ProbInterval{1, 0}
	var t6 float64
	DynMap[9] = diesel.ProbInterval{1, 0}
	x1 = getGPSX()
	y1 = getGPSY()
	DynMap[0].Reliability = 0.95
	DynMap[0].Delta = 7.8
	DynMap[1].Reliability = 0.95
	DynMap[1].Delta = 7.8

	for __temp_0 := 0; __temp_0 < 10; __temp_0++ {
		x2 = getGPSX()
		y2 = getGPSY()
		DynMap[4].Reliability = DynMap[2].Reliability + DynMap[0].Reliability - 1.0
		DynMap[4].Delta = DynMap[0].Delta + DynMap[2].Delta
		t1 = x1 - x2
		DynMap[5].Reliability = DynMap[4].Reliability
		DynMap[5].Delta = math.Abs(float64(t1))*DynMap[4].Delta + math.Abs(float64(t1))*DynMap[4].Delta + DynMap[4].Delta*DynMap[4].Delta
		t2 = t1 * t1
		DynMap[6].Reliability = DynMap[1].Reliability + DynMap[3].Reliability - 1.0
		DynMap[6].Delta = DynMap[1].Delta + DynMap[3].Delta
		t3 = y1 - y2
		DynMap[7].Reliability = DynMap[6].Reliability
		DynMap[7].Delta = math.Abs(float64(t3))*DynMap[6].Delta + math.Abs(float64(t3))*DynMap[6].Delta + DynMap[6].Delta*DynMap[6].Delta
		t4 = t3 * t3
		DynMap[8].Reliability = DynMap[7].Reliability + DynMap[5].Reliability - 1.0
		DynMap[8].Delta = DynMap[7].Delta + DynMap[5].Delta
		t5 = t4 + t2
		DynMap[9].Reliability = DynMap[8].Reliability
		DynMap[9].Delta = DynMap[8].Delta / math.Abs(2)
		t6 = t5 / 2

		fmt.Println(DynMap[9])
	}
	Out = t6

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
