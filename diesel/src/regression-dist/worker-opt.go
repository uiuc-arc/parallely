package main

import (
	"dieseldistrel"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
)

const NumWorkers = 10
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

func func_0() {
	dieseldistrel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
	dieseldistrel.WaitForWorkers(Num_threads)
	var DynMap [2005]float32
	var my_chan_index int
	_ = my_chan_index
	_ = DynMap
	var alpha float64
	DynMap[0] = 1.0
	var beta float64
	DynMap[1] = 1.0
	var idx0 int
	var idx1 int
	var xslice [1000]float64
	dieseldistrel.InitDynArray(2, 1000, DynMap[:])
	var yslice [1000]float64
	dieseldistrel.InitDynArray(1002, 1000, DynMap[:])
	var workerAlpha float64
	DynMap[2002] = 1.0
	var workerBeta float64
	DynMap[2003] = 1.0
	var tempF float64
	var tempDF float64
	DynMap[2004] = 1.0
	DynMap[0] = 1.0
	alpha = 0.0
	DynMap[1] = 1.0
	beta = 0.0
	dieseldistrel.StartTiming()
	idx0 = 0
	for _, q := range Q {
		idx1 = 0
		for __temp_0 := 0; __temp_0 < WorkPerThread; __temp_0++ {
			_temp_index_1 := idx0
			tempF = X[_temp_index_1]
			tempDF = tempF
			DynMap[2004] = 0.9999
			_temp_index_2 := idx1
			xslice[_temp_index_2] = tempDF
			DynMap[2+_temp_index_2] = DynMap[2004]
			_temp_index_3 := idx0
			tempF = Y[_temp_index_3]
			tempDF = tempF
			DynMap[2004] = 0.9999
			_temp_index_4 := idx1
			yslice[_temp_index_4] = tempDF
			DynMap[1002+_temp_index_4] = DynMap[2004]
			idx0 = idx0 + 1
			idx1 = idx1 + 1
		}
		dieseldistrel.SendFloat64Array(xslice[:], 0, q)
		dieseldistrel.SendFloat64Array(yslice[:], 0, q)
	}
	for _, q := range Q {
		dieseldistrel.ReceiveFloat64(&workerAlpha, 0, q)
		DynMap[2002] = float32(math.Pow(0.9999999, 2*1000))
		dieseldistrel.ReceiveFloat64(&workerBeta, 0, q)
		DynMap[2003] = float32(math.Pow(0.9999999, 2*1000))
		DynMap[2004] = DynMap[0] + DynMap[2002] - 1.0
		tempDF = alpha + workerAlpha
		DynMap[0] = DynMap[2004]
		alpha = tempDF
		DynMap[2004] = DynMap[1] + DynMap[2003] - 1.0
		tempDF = beta + workerBeta
		DynMap[1] = DynMap[2004]
		beta = tempDF
	}
	tempF = convertToFloat(NumWorkers)
	DynMap[2004] = DynMap[0]
	tempDF = alpha / tempF
	DynMap[0] = DynMap[2004]
	alpha = tempDF
	DynMap[2004] = DynMap[1]
	tempDF = beta / tempF
	DynMap[1] = DynMap[2004]
	beta = tempDF
	dieseldistrel.EndTiming()
	Alpha = alpha
	Beta = beta

	dieseldistrel.CleanupMain()
	fmt.Println("Ending thread : ", 0)
}
func func_Q(tid int) {
	dieseldistrel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
	dieseldistrel.PingMain(tid)
	var my_chan_index int
	_ = my_chan_index
	q := tid
	var xslice [1000]float64
	var yslice [1000]float64
	var x float64
	var y float64
	var mX float64
	var mY float64
	var ssXY float64
	var ssXX float64
	var alpha float64
	var beta float64
	var idx int
	var tempF float64
	var tempDF0 float64
	var tempDF1 float64
	dieseldistrel.ReceiveFloat64Array(xslice[:], tid, 0)
	dieseldistrel.ReceiveFloat64Array(yslice[:], tid, 0)
	mX = 0.0
	mY = 0.0
	ssXY = 0.0
	ssXX = 0.0
	idx = 0
	for __temp_3 := 0; __temp_3 < WorkPerThread; __temp_3++ {
		_temp_index_1 := idx
		x = xslice[_temp_index_1]
		_temp_index_2 := idx
		y = yslice[_temp_index_2]
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
		idx = idx + 1
	}
	tempF = convertToFloat(WorkPerThread)
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
	dieseldistrel.SendFloat64(alpha, tid, 0)
	dieseldistrel.SendFloat64(beta, tid, 0)
	fmt.Println("Ending thread : ", q)
}

func main() {
	tid, _ := strconv.Atoi(os.Args[1])
	fmt.Println("Starting worker thread: ", tid)

	totalWork = WorkPerThread * NumWorkers
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
