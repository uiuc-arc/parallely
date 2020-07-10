package main

import (
	"diesel"
	"fmt"
	"math"
	"math/rand"
)

const numWorkers = 10

var WorkPerThread, totalWork int
var X, Y []float64
var Num_threads int

var Alpha, Beta float64

const (
	// single whitespace character
	ws = "[ \n\r\t\v\f]"
	// isolated comment
	cmt = "#[^\n\r]*"
	// comment sub expression
	cmts = "(" + ws + "*" + cmt + "[\n\r])"
	// number with leading comments
	num = "(" + cmts + "+" + ws + "*|" + ws + "+)([0-9]+)"
)

func convertToFloat(x int) float64 {
	return float64(x)
}

var Q = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func func_0() {
	diesel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
	diesel.WaitForWorkers(Num_threads)
	var DynMap [6]diesel.ProbInterval
	var my_chan_index int
	_ = my_chan_index
	_ = DynMap
	var workerAlpha float64
	DynMap[0] = diesel.ProbInterval{1, 0}
	var workerBeta float64
	DynMap[1] = diesel.ProbInterval{1, 0}
	var workerSamples int
	var alpha float64
	DynMap[2] = diesel.ProbInterval{1, 0}
	var beta float64
	DynMap[3] = diesel.ProbInterval{1, 0}
	var totalSamples int
	var tempF float64
	var tempDF0 float64
	DynMap[4] = diesel.ProbInterval{1, 0}
	var tempDF1 float64
	DynMap[5] = diesel.ProbInterval{1, 0}
	totalSamples = 0
	DynMap[2] = diesel.ProbInterval{1, 0}
	alpha = 0.0
	DynMap[3] = diesel.ProbInterval{1, 0}
	beta = 0.0
	diesel.StartTiming()
	for _, q := range Q {
		diesel.ReceiveFloat64(&workerAlpha, 0, q)
		DynMap[0] = diesel.ProbInterval{1.0, 1e-7}
		diesel.ReceiveFloat64(&workerBeta, 0, q)
		DynMap[1] = diesel.ProbInterval{1.0, 1e-7}
		diesel.ReceiveInt(&workerSamples, 0, q)
		tempF = convertToFloat(workerSamples)
		DynMap[4].Reliability = DynMap[0].Reliability
		DynMap[4].Delta = math.Abs(float64(tempF)) * DynMap[0].Delta
		tempDF0 = workerAlpha * tempF
		DynMap[5].Reliability = DynMap[2].Reliability + DynMap[4].Reliability - 1.0
		DynMap[5].Delta = DynMap[2].Delta + DynMap[4].Delta
		tempDF1 = alpha + tempDF0
		DynMap[2].Reliability = DynMap[5].Reliability
		DynMap[2].Delta = DynMap[5].Delta
		alpha = tempDF1
		DynMap[4].Reliability = DynMap[1].Reliability
		DynMap[4].Delta = math.Abs(float64(tempF)) * DynMap[1].Delta
		tempDF0 = workerBeta * tempF
		DynMap[5].Reliability = DynMap[3].Reliability + DynMap[4].Reliability - 1.0
		DynMap[5].Delta = DynMap[3].Delta + DynMap[4].Delta
		tempDF1 = beta + tempDF0
		DynMap[3].Reliability = DynMap[5].Reliability
		DynMap[3].Delta = DynMap[5].Delta
		beta = tempDF1
		totalSamples = totalSamples + workerSamples
	}
	tempF = convertToFloat(totalSamples)
	DynMap[4].Reliability = DynMap[2].Reliability
	DynMap[4].Delta = DynMap[2].Delta / math.Abs(tempF)
	tempDF0 = alpha / tempF
	DynMap[2].Reliability = DynMap[4].Reliability
	DynMap[2].Delta = DynMap[4].Delta
	alpha = tempDF0
	DynMap[4].Reliability = DynMap[3].Reliability
	DynMap[4].Delta = DynMap[3].Delta / math.Abs(tempF)
	tempDF0 = beta / tempF
	DynMap[3].Reliability = DynMap[4].Reliability
	DynMap[3].Delta = DynMap[4].Delta
	beta = tempDF0
	diesel.EndTiming()
	Alpha = alpha
	Beta = beta

	diesel.CleanupMain()
	fmt.Println("Ending thread : ", 0)
}

func main() {
	// rand.Seed(time.Now().UTC().UnixNano())
	seed := int64(12345)
	rand.Seed(seed) // deterministic seed for reproducibility

	WorkPerThread = 100
	totalWork = WorkPerThread * numWorkers
	X = make([]float64, totalWork)
	Y = make([]float64, totalWork)

	fmt.Println("Generating", totalWork, "points using random seed", seed)

	alpha := rand.NormFloat64()
	beta := rand.NormFloat64()

	for i := 0; i < totalWork; i++ {
		X[i] = rand.NormFloat64() * math.Abs(100.0)   // always use math library to satisfy Go
		Y[i] = alpha + beta*(X[i]+rand.NormFloat64()) // add some error
	}

	fmt.Println("Starting program")

	Num_threads = 11

	func_0()

	// fmt.Println("Main thread waiting for others to finish");
	// diesel.Wg.Wait()

	// end := time.Now()
	// elapsed := end.Sub(startTime)
	// fmt.Println("Elapsed time :", elapsed.Nanoseconds())

	// fmt.Println("Alpha: actual", alpha, "estimate", Alpha)
	// fmt.Println("Beta : actual", beta , "estimate", Beta )
}
