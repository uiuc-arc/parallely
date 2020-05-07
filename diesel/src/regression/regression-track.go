package main

import (
  "fmt"
  "math"
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
  var DynMap [6]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var workerAlpha float64;
DynMap[0] = diesel.ProbInterval{1, 0};
var workerBeta float64;
DynMap[1] = diesel.ProbInterval{1, 0};
var workerSamples int;
var alpha float64;
DynMap[2] = diesel.ProbInterval{1, 0};
var beta float64;
DynMap[3] = diesel.ProbInterval{1, 0};
var totalSamples int;
var tempF float64;
var tempDF0 float64;
DynMap[4] = diesel.ProbInterval{1, 0};
var tempDF1 float64;
DynMap[5] = diesel.ProbInterval{1, 0};
totalSamples = 0;
DynMap[2] = diesel.ProbInterval{1, 0};
DynMap[2].Delta = 0;
alpha = 0.0;
DynMap[3] = diesel.ProbInterval{1, 0};
DynMap[3].Delta = 0;
beta = 0.0;
for _, q := range(Q) {
 diesel.ReceiveFloat64(&workerAlpha, 0, q);
my_chan_index = q * diesel.Numprocesses + 0;
__temp_rec_val_1 := <- diesel.DynamicChannelMap[my_chan_index];
DynMap[0] = __temp_rec_val_1;
diesel.ReceiveFloat64(&workerBeta, 0, q);
my_chan_index = q * diesel.Numprocesses + 0;
__temp_rec_val_2 := <- diesel.DynamicChannelMap[my_chan_index];
DynMap[1] = __temp_rec_val_2;
diesel.ReceiveInt(&workerSamples, 0, q);
tempF=convertToFloat(workerSamples);
DynMap[4].Reliability = DynMap[0].Reliability;
DynMap[4].Delta = math.Abs(workerAlpha) * DynMap[0].Delta;
tempDF0 = workerAlpha*tempF;
DynMap[5].Reliability = DynMap[2].Reliability + DynMap[4].Reliability - 1.0;
DynMap[5].Delta = DynMap[2].Delta + DynMap[4].Delta;
tempDF1 = alpha+tempDF0;
DynMap[2].Reliability = DynMap[5].Reliability;
alpha = tempDF1;
DynMap[4].Reliability = DynMap[1].Reliability;
DynMap[4].Delta = math.Abs(workerBeta) * DynMap[1].Delta;
tempDF0 = workerBeta*tempF;
DynMap[5].Reliability = DynMap[3].Reliability + DynMap[4].Reliability - 1.0;
DynMap[5].Delta = DynMap[3].Delta + DynMap[4].Delta;
tempDF1 = beta+tempDF0;
DynMap[3].Reliability = DynMap[5].Reliability;
beta = tempDF1;
totalSamples = totalSamples+workerSamples;
 }
tempF=convertToFloat(totalSamples);
DynMap[4].Reliability = DynMap[2].Reliability;
DynMap[4].Delta =  DynMap[2].Delta / math.Abs(tempF);
tempDF0 = alpha/tempF;
DynMap[2].Reliability = DynMap[4].Reliability;
alpha = tempDF0;
DynMap[4].Reliability = DynMap[3].Reliability;
DynMap[4].Delta =  DynMap[3].Delta / math.Abs(tempF);
tempDF0 = beta/tempF;
DynMap[3].Reliability = DynMap[4].Reliability;
beta = tempDF0;
Alpha = alpha;
Beta = beta;


  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  defer diesel.Wg.Done();
  var DynMap [10]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
var x float64;
DynMap[0] = diesel.ProbInterval{1, 0};
var y float64;
DynMap[1] = diesel.ProbInterval{1, 0};
var mX float64;
DynMap[2] = diesel.ProbInterval{1, 0};
var mY float64;
DynMap[3] = diesel.ProbInterval{1, 0};
var ssXY float64;
DynMap[4] = diesel.ProbInterval{1, 0};
var ssXX float64;
DynMap[5] = diesel.ProbInterval{1, 0};
var alpha float64;
DynMap[6] = diesel.ProbInterval{1, 0};
var beta float64;
DynMap[7] = diesel.ProbInterval{1, 0};
var count int;
var idx int;
var tempF float64;
var tempDF0 float64;
DynMap[8] = diesel.ProbInterval{1, 0};
var tempDF1 float64;
DynMap[9] = diesel.ProbInterval{1, 0};
DynMap[2] = diesel.ProbInterval{1, 0};
DynMap[2].Delta = 0;
mX = 0.0;
DynMap[3] = diesel.ProbInterval{1, 0};
DynMap[3].Delta = 0;
mY = 0.0;
DynMap[4] = diesel.ProbInterval{1, 0};
DynMap[4].Delta = 0;
ssXY = 0.0;
DynMap[5] = diesel.ProbInterval{1, 0};
DynMap[5].Delta = 0;
ssXX = 0.0;
count = 0;
idx = 0;
for __temp_2 := 0; __temp_2 < WorkPerThread; __temp_2++ {
 _temp_index_1 := ((q-1)*WorkPerThread)+idx;
tempF=X[_temp_index_1];
x=tempF;
DynMap[0] = diesel.ProbInterval{0.99, 0.001};
_temp_index_2 := ((q-1)*WorkPerThread)+idx;
tempF=Y[_temp_index_2];
y=tempF;
DynMap[1] = diesel.ProbInterval{0.99, 0.001};
DynMap[8].Reliability = DynMap[0].Reliability + DynMap[2].Reliability - 1.0;
DynMap[8].Delta = DynMap[2].Delta + DynMap[0].Delta;
tempDF0 = mX+x;
DynMap[2].Reliability = DynMap[8].Reliability;
mX = tempDF0;
DynMap[8].Reliability = DynMap[1].Reliability + DynMap[3].Reliability - 1.0;
DynMap[8].Delta = DynMap[3].Delta + DynMap[1].Delta;
tempDF0 = mY+y;
DynMap[3].Reliability = DynMap[8].Reliability;
mY = tempDF0;
DynMap[8].Reliability = DynMap[1].Reliability + DynMap[0].Reliability - 1.0;
DynMap[8].Delta = math.Abs(x) * DynMap[0].Delta + math.Abs(y) * DynMap[1].Delta + DynMap[0].Delta*DynMap[1].Delta;
tempDF0 = x*y;
DynMap[9].Reliability = DynMap[8].Reliability + DynMap[4].Reliability - 1.0;
DynMap[9].Delta = DynMap[4].Delta + DynMap[8].Delta;
tempDF1 = ssXY+tempDF0;
DynMap[4].Reliability = DynMap[9].Reliability;
ssXY = tempDF1;
DynMap[8].Reliability = DynMap[0].Reliability;
DynMap[8].Delta = math.Abs(x) * DynMap[0].Delta + math.Abs(x) * DynMap[0].Delta + DynMap[0].Delta*DynMap[0].Delta;
tempDF0 = x*x;
DynMap[9].Reliability = DynMap[5].Reliability + DynMap[8].Reliability - 1.0;
DynMap[9].Delta = DynMap[5].Delta + DynMap[8].Delta;
tempDF1 = ssXX+tempDF0;
DynMap[5].Reliability = DynMap[9].Reliability;
ssXX = tempDF1;
count = count+1;
idx = idx+1;
 }
tempF=convertToFloat(count);
DynMap[8].Reliability = DynMap[2].Reliability;
DynMap[8].Delta =  DynMap[2].Delta / math.Abs(tempF);
tempDF0 = mX/tempF;
DynMap[2].Reliability = DynMap[8].Reliability;
mX = tempDF0;
DynMap[8].Reliability = DynMap[3].Reliability;
DynMap[8].Delta =  DynMap[3].Delta / math.Abs(tempF);
tempDF0 = mY/tempF;
DynMap[3].Reliability = DynMap[8].Reliability;
mY = tempDF0;
DynMap[8].Reliability = DynMap[3].Reliability + DynMap[2].Reliability - 1.0;
DynMap[8].Delta = math.Abs(mX) * DynMap[2].Delta + math.Abs(mY) * DynMap[3].Delta + DynMap[2].Delta*DynMap[3].Delta;
tempDF0 = mX*mY;
DynMap[9].Reliability = DynMap[8].Reliability;
DynMap[9].Delta = math.Abs(tempDF0) * DynMap[8].Delta;
tempDF1 = tempDF0*tempF;
DynMap[8].Reliability = DynMap[9].Reliability + DynMap[4].Reliability - 1.0;
DynMap[8].Delta = DynMap[4].Delta + DynMap[9].Delta;
tempDF0 = ssXY-tempDF1;
DynMap[4].Reliability = DynMap[8].Reliability;
ssXY = tempDF0;
DynMap[8].Reliability = DynMap[2].Reliability;
DynMap[8].Delta = math.Abs(mX) * DynMap[2].Delta + math.Abs(mX) * DynMap[2].Delta + DynMap[2].Delta*DynMap[2].Delta;
tempDF0 = mX*mX;
DynMap[9].Reliability = DynMap[8].Reliability;
DynMap[9].Delta = math.Abs(tempDF0) * DynMap[8].Delta;
tempDF1 = tempDF0*tempF;
DynMap[8].Reliability = DynMap[9].Reliability + DynMap[5].Reliability - 1.0;
DynMap[8].Delta = DynMap[5].Delta + DynMap[9].Delta;
tempDF0 = ssXX-tempDF1;
DynMap[5].Reliability = DynMap[8].Reliability;
ssXX = tempDF0;
DynMap[7].Reliability = DynMap[5].Reliability + DynMap[4].Reliability - 1.0;
DynMap[7].Delta = math.Abs(ssXY) * DynMap[4].Delta + math.Abs(ssXX) * DynMap[5].Delta / (math.Abs(ssXX) * (math.Abs(ssXY)-DynMap[5].Delta));
beta = ssXY/ssXX;
DynMap[8].Reliability = DynMap[7].Reliability + DynMap[2].Reliability - 1.0;
DynMap[8].Delta = math.Abs(beta) * DynMap[7].Delta + math.Abs(mX) * DynMap[2].Delta + DynMap[7].Delta*DynMap[2].Delta;
tempDF0 = beta*mX;
DynMap[6].Reliability = DynMap[3].Reliability + DynMap[8].Reliability - 1.0;
DynMap[6].Delta = DynMap[3].Delta + DynMap[8].Delta;
alpha = mY-tempDF0;
diesel.SendFloat64(alpha, tid, 0);
diesel.SendDynVal(DynMap[6], tid, 0);
diesel.SendFloat64(beta, tid, 0);
diesel.SendDynVal(DynMap[7], tid, 0);
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
