package main

import (
	"diesel"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
)

const numWorkers = 10
const WorkPerThread = 100

var totalWork int
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

func func_Q(tid int) {
	diesel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
	diesel.PingMain(tid)

	var my_chan_index int
	_ = my_chan_index

	q := tid
	var x float64
	var y float64
	var mX float64
	var mY float64
	var ssXY float64
	var ssXX float64
	var alpha float64
	var beta float64
	var count int
	var idx int
	var tempF float64
	var tempDF0 float64
	var tempDF1 float64
	mX = 0.0
	mY = 0.0
	ssXY = 0.0
	ssXX = 0.0
	count = 0
	idx = 0
	for __temp_2 := 0; __temp_2 < WorkPerThread; __temp_2++ {
		_temp_index_1 := ((q - 1) * WorkPerThread) + idx
		tempF = X[_temp_index_1]
		x = tempF
		_temp_index_2 := ((q - 1) * WorkPerThread) + idx
		tempF = Y[_temp_index_2]
		y = tempF
		tempDF0 = mX + x

		mX = tempDF0
		tempDF0 = mY + y
		mY = tempDF0
		tempDF0 = x * y
		tempDF1 = ssXY + tempDF0
		ssXY = tempDF1
		tempDF0 = x * x
		tempDF1 = ssXX + tempDF0
		ssXX = tempDF1
		count = count + 1
		idx = idx + 1
	}
	tempF = convertToFloat(count)
	tempDF0 = mX / tempF
	mX = tempDF0
	tempDF0 = mY / tempF
	mY = tempDF0
	tempDF0 = mX * mY
	tempDF1 = tempDF0 * tempF
	tempDF0 = ssXY - tempDF1
	ssXY = tempDF0
	tempDF0 = mX * mX
	tempDF1 = tempDF0 * tempF
	tempDF0 = ssXX - tempDF1
	ssXX = tempDF0
	beta = ssXY / ssXX
	tempDF0 = beta * mX

	alpha = mY - tempDF0
	diesel.SendFloat64(alpha, tid, 0)
	// diesel.SendDynVal(DynMap[6], tid, 0)
	diesel.SendFloat64(beta, tid, 0)
	// diesel.SendDynVal(DynMap[7], tid, 0)
	diesel.SendInt(count, tid, 0)

	diesel.CleanupMain()
	fmt.Println("Ending thread : ", q)
}

func main() {
	tid, _ := strconv.Atoi(os.Args[1])
	fmt.Println("Starting worker thread: ", tid)

	totalWork = WorkPerThread * numWorkers
	seed := 0
	X = make([]float64, totalWork)
	Y = make([]float64, totalWork)

	fmt.Println("Generating", totalWork, "points using random seed", seed)

	alpha := rand.NormFloat64()
	beta := rand.NormFloat64()

	for i := 0; i < totalWork; i++ {
		X[i] = rand.NormFloat64() * math.Abs(100.0)   // always use math library to satisfy Go
		Y[i] = alpha + beta*(X[i]+rand.NormFloat64()) // add some error
	}

	Num_threads = 11

	func_Q(tid)
}
