package main

import (
  "fmt"
  "math/rand"
  "time"
  "diesel"
)

const numWorkers = 10
const WorkPerThread = 100
const totalWork = numWorkers*WorkPerThread
var X [totalWork]float64
var Y [totalWork]float64

var Alpha float64
var Beta  float64

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
var tempDF float64;
totalSamples = 0;
alpha = 0.0;
beta = 0.0;
for _, q := range(Q) {
 diesel.ReceiveFloat64(&workerAlpha, 0, q);
diesel.ReceiveFloat64(&workerBeta, 0, q);
diesel.ReceiveInt(&workerSamples, 0, q);
tempF=convertToFloat(workerSamples);
tempDF = workerAlpha*tempF;
alpha = alpha+tempDF;
tempDF = workerBeta*tempF;
beta = beta+tempDF;
totalSamples = totalSamples+workerSamples;
 }
tempF=convertToFloat(totalSamples);
alpha = alpha/tempF;
beta = beta/tempF;
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
var tempDF float64;
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
mX = mX+x;
mY = mY+y;
tempDF = x*y;
ssXY = ssXY+tempDF;
tempDF = x*x;
ssXX = ssXX+tempDF;
count = count+1;
idx = idx+1;
 }
tempF=convertToFloat(count);
mX = mX/tempF;
mY = mY/tempF;
tempDF = mX*mY;
tempDF = tempDF*tempF;
ssXY = ssXY-tempDF;
tempDF = mX*mX;
tempDF = tempDF*tempF;
ssXX = ssXX-tempDF;
beta = ssXY/ssXX;
tempDF = beta*mX;
alpha = mY-tempDF;
diesel.SendFloat64(alpha, tid, 0);
diesel.SendFloat64(beta, tid, 0);
diesel.SendInt(count, tid, 0);

  fmt.Println("Ending thread : ", q);
}

func main() {
  // rand.Seed(time.Now().UTC().UnixNano())
  seed := int64(12345)
  rand.Seed(seed) // deterministic seed for reproducibility

  fmt.Println("Generating data using seed",seed)

  alpha := rand.NormFloat64()
  beta  := rand.NormFloat64()

  for i := 0; i < totalWork; i++ {
    X[i] = rand.NormFloat64()*100
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
}
