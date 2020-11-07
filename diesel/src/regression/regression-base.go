package main

import (
<<<<<<< HEAD
	"diesel"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const numWorkers = 10

=======
  "fmt"
  "math"
  "math/rand"
  "os"
  "strconv"
  "time"
  "diesel"
)

const numWorkers = 10
>>>>>>> 23ae881cc9f84113e4fcd4c6bce92cae87a80620
var WorkPerThread, totalWork int
var X, Y []float64

var Alpha, Beta float64

const (
<<<<<<< HEAD
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
	defer diesel.Wg.Done()
	var DynMap [0]diesel.ProbInterval
	var my_chan_index int
	_ = my_chan_index
	_ = DynMap
	var workerAlpha float64
	var workerBeta float64
	var workerSamples int
	var alpha float64
	var beta float64
	var totalSamples int
	var tempF float64
	var tempDF0 float64
	var tempDF1 float64
	totalSamples = 0
	alpha = 0.0
	beta = 0.0
	for _, q := range Q {
		diesel.ReceiveFloat64(&workerAlpha, 0, q)
		diesel.ReceiveFloat64(&workerBeta, 0, q)
		diesel.ReceiveInt(&workerSamples, 0, q)
		tempF = convertToFloat(workerSamples)
		tempDF0 = workerAlpha * tempF
		tempDF1 = alpha + tempDF0
		alpha = tempDF1
		tempDF0 = workerBeta * tempF
		tempDF1 = beta + tempDF0
		beta = tempDF1
		totalSamples = totalSamples + workerSamples
	}
	tempF = convertToFloat(totalSamples)
	tempDF0 = alpha / tempF
	alpha = tempDF0
	tempDF0 = beta / tempF
	beta = tempDF0
	Alpha = alpha
	Beta = beta

	fmt.Println("Ending thread : ", 0)
}

func func_Q(tid int) {
	defer diesel.Wg.Done()
	var DynMap [0]diesel.ProbInterval
	var my_chan_index int
	_ = my_chan_index
	_ = DynMap
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
	for __temp_0 := 0; __temp_0 < WorkPerThread; __temp_0++ {
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
	diesel.SendFloat64(beta, tid, 0)
	diesel.SendInt(count, tid, 0)

	fmt.Println("Ending thread : ", q)
}

func main() {
	// rand.Seed(time.Now().UTC().UnixNano())
	seed := int64(12345)
	rand.Seed(seed) // deterministic seed for reproducibility

	WorkPerThread, _ = strconv.Atoi(os.Args[1])
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

	diesel.InitChannels(11)

	startTime := time.Now()
	go func_0()
	for _, index := range Q {
		go func_Q(index)
	}

	fmt.Println("Main thread waiting for others to finish")
	diesel.Wg.Wait()

	end := time.Now()
	elapsed := end.Sub(startTime)
	fmt.Println("Elapsed time :", elapsed.Nanoseconds())

	fmt.Println("Alpha: actual", alpha, "estimate", Alpha)
	fmt.Println("Beta : actual", beta, "estimate", Beta)
=======
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

var Q = []int {1,2,3,4,5,6,7,8,9,10};


func func_0() {
  defer diesel.Wg.Done();
  var DynMap [0]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var workerAlpha float64;
var workerBeta float64;
var workerSamples int;
var alpha float64;
var beta float64;
var totalSamples int;
var tempF float64;
var tempDF0 float64;
var tempDF1 float64;
totalSamples = 0;
alpha = 0.0;
beta = 0.0;
for _, q := range(Q) {
 diesel.ReceiveFloat64(&workerAlpha, 0, q);
diesel.ReceiveFloat64(&workerBeta, 0, q);
diesel.ReceiveInt(&workerSamples, 0, q);
tempF=convertToFloat(workerSamples);
tempDF0 = workerAlpha*tempF;
tempDF1 = alpha+tempDF0;
alpha = tempDF1;
tempDF0 = workerBeta*tempF;
tempDF1 = beta+tempDF0;
beta = tempDF1;
totalSamples = totalSamples+workerSamples;
 }
tempF=convertToFloat(totalSamples);
tempDF0 = alpha/tempF;
alpha = tempDF0;
tempDF0 = beta/tempF;
beta = tempDF0;
Alpha = alpha;
Beta = beta;


  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  defer diesel.Wg.Done();
  var DynMap [0]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
var x float64;
var y float64;
var mX float64;
var mY float64;
var ssXY float64;
var ssXX float64;
var alpha float64;
var beta float64;
var count int;
var idx int;
var tempF float64;
var tempDF0 float64;
var tempDF1 float64;
mX = 0.0;
mY = 0.0;
ssXY = 0.0;
ssXX = 0.0;
count = 0;
idx = 0;
for __temp_0 := 0; __temp_0 < WorkPerThread; __temp_0++ {
 _temp_index_1 := ((q-1)*WorkPerThread)+idx;
tempF=X[_temp_index_1];
x=tempF;
_temp_index_2 := ((q-1)*WorkPerThread)+idx;
tempF=Y[_temp_index_2];
y=tempF;
tempDF0 = mX+x;
mX = tempDF0;
tempDF0 = mY+y;
mY = tempDF0;
tempDF0 = x*y;
tempDF1 = ssXY+tempDF0;
ssXY = tempDF1;
tempDF0 = x*x;
tempDF1 = ssXX+tempDF0;
ssXX = tempDF1;
count = count+1;
idx = idx+1;
 }
tempF=convertToFloat(count);
tempDF0 = mX/tempF;
mX = tempDF0;
tempDF0 = mY/tempF;
mY = tempDF0;
tempDF0 = mX*mY;
tempDF1 = tempDF0*tempF;
tempDF0 = ssXY-tempDF1;
ssXY = tempDF0;
tempDF0 = mX*mX;
tempDF1 = tempDF0*tempF;
tempDF0 = ssXX-tempDF1;
ssXX = tempDF0;
beta = ssXY/ssXX;
tempDF0 = beta*mX;
alpha = mY-tempDF0;
diesel.SendFloat64(alpha, tid, 0);
diesel.SendFloat64(beta, tid, 0);
diesel.SendInt(count, tid, 0);

  fmt.Println("Ending thread : ", q);
}

func main() {
  // rand.Seed(time.Now().UTC().UnixNano())
  seed := int64(12345)
  rand.Seed(seed) // deterministic seed for reproducibility

  WorkPerThread, _ = strconv.Atoi(os.Args[1])
  totalWork = WorkPerThread*numWorkers
  X = make([]float64, totalWork)
  Y = make([]float64, totalWork)

  fmt.Println("Generating",totalWork,"points using random seed",seed)

  alpha := rand.NormFloat64()
  beta  := rand.NormFloat64()

  for i := 0; i < totalWork; i++ {
    X[i] = rand.NormFloat64()*math.Abs(100.0) // always use math library to satisfy Go
    Y[i] = alpha + beta*(X[i]+rand.NormFloat64()) // add some error
  }

  fmt.Println("Starting program");
	
  diesel.InitChannels(11);

  startTime := time.Now()
  go func_0();
for _, index := range Q {
go func_Q(index);
}


  fmt.Println("Main thread waiting for others to finish");  
  diesel.Wg.Wait()

  end := time.Now()
  elapsed := end.Sub(startTime)
  fmt.Println("Elapsed time :", elapsed.Nanoseconds())

  fmt.Println("Alpha: actual", alpha, "estimate", Alpha)
  fmt.Println("Beta : actual", beta , "estimate", Beta )
>>>>>>> 23ae881cc9f84113e4fcd4c6bce92cae87a80620
}
