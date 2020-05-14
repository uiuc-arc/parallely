package main

import (
  "fmt"
  "math"
  "math/rand"
  "time"
  "diesel"
)

const NumWorkers = 10
const totalWork = NumWorkers*1000
var X, Y [totalWork]float64

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

var Q = []int {1,2,3,4,5,6,7,8,9,10};


func func_0() {
  defer diesel.Wg.Done();
  var DynMap [2009]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var alpha float64;
DynMap[0] = diesel.ProbInterval{1, 0};
var beta float64;
DynMap[1] = diesel.ProbInterval{1, 0};
var idx0 int;
var idx1 int;
var xslice [1000]float64;
diesel.InitDynArray(2, 1000, DynMap[:]);
var yslice [1000]float64;
diesel.InitDynArray(1002, 1000, DynMap[:]);
var workerAlpha float64;
DynMap[2002] = diesel.ProbInterval{1, 0};
var workerBeta float64;
DynMap[2003] = diesel.ProbInterval{1, 0};
var tempF float64;
var mX float64;
DynMap[2004] = diesel.ProbInterval{1, 0};
var mY float64;
DynMap[2005] = diesel.ProbInterval{1, 0};
var ssXX float64;
DynMap[2006] = diesel.ProbInterval{1, 0};
var ssXY float64;
DynMap[2007] = diesel.ProbInterval{1, 0};
var tempDF float64;
DynMap[2008] = diesel.ProbInterval{1, 0};
DynMap[0] = diesel.ProbInterval{1, 0};
alpha = 0.0;
DynMap[1] = diesel.ProbInterval{1, 0};
beta = 0.0;
idx0 = 0;
for _, q := range(Q) {
 idx1 = 0;
for __temp_0 := 0; __temp_0 < 1000; __temp_0++ {
 _temp_index_1 := idx0;
tempF=X[_temp_index_1];
tempDF=tempF;
DynMap[2008] = diesel.ProbInterval{1.0, 1e-7};
_temp_index_2 := idx1;
xslice[_temp_index_2]=tempDF;
DynMap[2 + _temp_index_2] = DynMap[2008];
_temp_index_3 := idx0;
tempF=Y[_temp_index_3];
tempDF=tempF;
DynMap[2008] = diesel.ProbInterval{1.0, 1e-7};
_temp_index_4 := idx1;
yslice[_temp_index_4]=tempDF;
DynMap[1002 + _temp_index_4] = DynMap[2008];
idx0 = idx0+1;
idx1 = idx1+1;
 }
diesel.SendDynFloat64Array(xslice[:], 0, q, DynMap[:], 2);
diesel.SendDynFloat64Array(yslice[:], 0, q, DynMap[:], 1002);
 }
for _, q := range(Q) {
 diesel.ReceiveFloat64(&mX, 0, q);
my_chan_index = q * diesel.Numprocesses + 0;
__temp_rec_val_2 := <- diesel.DynamicChannelMap[my_chan_index];
DynMap[2004] = __temp_rec_val_2;
DynMap[2004] = diesel.ProbInterval{1.0, 1e-7} // d(mX) <= d(xSlice)
diesel.ReceiveFloat64(&mY, 0, q);
my_chan_index = q * diesel.Numprocesses + 0;
__temp_rec_val_3 := <- diesel.DynamicChannelMap[my_chan_index];
DynMap[2005] = __temp_rec_val_3;
DynMap[2005] = diesel.ProbInterval{1.0, 1e-7} // d(mY) <= d(ySlice)
diesel.ReceiveFloat64(&ssXX, 0, q);
my_chan_index = q * diesel.Numprocesses + 0;
__temp_rec_val_4 := <- diesel.DynamicChannelMap[my_chan_index];
DynMap[2006] = __temp_rec_val_4;
DynMap[2006] = diesel.ProbInterval{1.0, 0.04} // d(ssXX) <= 4*PtsPerWorker*maxAV(xSlice)*d(xSlice)
diesel.ReceiveFloat64(&ssXY, 0, q);
my_chan_index = q * diesel.Numprocesses + 0;
__temp_rec_val_5 := <- diesel.DynamicChannelMap[my_chan_index];
DynMap[2007] = __temp_rec_val_5;
DynMap[2007] = diesel.ProbInterval{1.0, 0.04} // d(ssXY) <= 2*PtsPerWorker*(maxAV(ySlice)*d(xSlice)+maxAV(xSlice)*d(ySlice))
DynMap[2003].Reliability = DynMap[2006].Reliability + DynMap[2007].Reliability - 1.0;
DynMap[2003].Delta = math.Abs(ssXY) * DynMap[2007].Delta + math.Abs(ssXX) * DynMap[2006].Delta / (math.Abs(ssXX) * (math.Abs(ssXY)-DynMap[2006].Delta));
workerBeta = ssXY/ssXX;
DynMap[2008].Reliability = DynMap[2004].Reliability + DynMap[2003].Reliability - 1.0;
DynMap[2008].Delta = math.Abs(float64(workerBeta)) * DynMap[2003].Delta + math.Abs(float64(mX)) * DynMap[2004].Delta + DynMap[2003].Delta*DynMap[2004].Delta;
tempDF = workerBeta*mX;
DynMap[2002].Reliability = DynMap[2005].Reliability + DynMap[2008].Reliability - 1.0;
DynMap[2002].Delta = DynMap[2005].Delta + DynMap[2008].Delta;
workerAlpha = mY-tempDF;
DynMap[2008].Reliability = DynMap[0].Reliability + DynMap[2002].Reliability - 1.0;
DynMap[2008].Delta = DynMap[0].Delta + DynMap[2002].Delta;
tempDF = alpha+workerAlpha;
DynMap[0].Reliability = DynMap[2008].Reliability;
DynMap[0].Delta = DynMap[2008].Delta;
alpha = tempDF;
DynMap[2008].Reliability = DynMap[1].Reliability + DynMap[2003].Reliability - 1.0;
DynMap[2008].Delta = DynMap[1].Delta + DynMap[2003].Delta;
tempDF = beta+workerBeta;
DynMap[1].Reliability = DynMap[2008].Reliability;
DynMap[1].Delta = DynMap[2008].Delta;
beta = tempDF;
 }
tempF=convertToFloat(NumWorkers);
DynMap[2008].Reliability = DynMap[0].Reliability;
DynMap[2008].Delta =  DynMap[0].Delta / math.Abs(tempF);
tempDF = alpha/tempF;
DynMap[0].Reliability = DynMap[2008].Reliability;
DynMap[0].Delta = DynMap[2008].Delta;
alpha = tempDF;
DynMap[2008].Reliability = DynMap[1].Reliability;
DynMap[2008].Delta =  DynMap[1].Delta / math.Abs(tempF);
tempDF = beta/tempF;
DynMap[1].Reliability = DynMap[2008].Reliability;
DynMap[1].Delta = DynMap[2008].Delta;
beta = tempDF;
Alpha = alpha;
Beta = beta;
  fmt.Println(DynMap[0], DynMap[1])

  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  defer diesel.Wg.Done();
  var DynMap [2008]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
var xslice [1000]float64;
diesel.InitDynArray(0, 1000, DynMap[:]);
var yslice [1000]float64;
diesel.InitDynArray(1000, 1000, DynMap[:]);
var x float64;
DynMap[2000] = diesel.ProbInterval{1, 0};
var y float64;
DynMap[2001] = diesel.ProbInterval{1, 0};
var mX float64;
DynMap[2002] = diesel.ProbInterval{1, 0};
var mY float64;
DynMap[2003] = diesel.ProbInterval{1, 0};
var ssXY float64;
DynMap[2004] = diesel.ProbInterval{1, 0};
var ssXX float64;
DynMap[2005] = diesel.ProbInterval{1, 0};
var idx int;
var tempF float64;
var tempDF0 float64;
DynMap[2006] = diesel.ProbInterval{1, 0};
var tempDF1 float64;
DynMap[2007] = diesel.ProbInterval{1, 0};
diesel.ReceiveDynFloat64Array(xslice[:], tid, 0, DynMap[:], 0);
diesel.ReceiveDynFloat64Array(yslice[:], tid, 0, DynMap[:], 1000);
DynMap[2002] = diesel.ProbInterval{1, 0};
mX = 0.0;
DynMap[2003] = diesel.ProbInterval{1, 0};
mY = 0.0;
DynMap[2004] = diesel.ProbInterval{1, 0};
ssXY = 0.0;
DynMap[2005] = diesel.ProbInterval{1, 0};
ssXX = 0.0;
idx = 0;
for __temp_5 := 0; __temp_5 < 1000; __temp_5++ {
 _temp_index_1 := idx;
x=xslice[_temp_index_1];
DynMap[2000] = DynMap[0 + _temp_index_1];
_temp_index_2 := idx;
y=yslice[_temp_index_2];
DynMap[2001] = DynMap[1000 + _temp_index_2];
DynMap[2006].Reliability = DynMap[2000].Reliability + DynMap[2002].Reliability - 1.0;
DynMap[2006].Delta = DynMap[2002].Delta + DynMap[2000].Delta;
tempDF0 = mX+x;
DynMap[2002].Reliability = DynMap[2006].Reliability;
DynMap[2002].Delta = DynMap[2006].Delta;
mX = tempDF0;
DynMap[2006].Reliability = DynMap[2001].Reliability + DynMap[2003].Reliability - 1.0;
DynMap[2006].Delta = DynMap[2003].Delta + DynMap[2001].Delta;
tempDF0 = mY+y;
DynMap[2003].Reliability = DynMap[2006].Reliability;
DynMap[2003].Delta = DynMap[2006].Delta;
mY = tempDF0;
DynMap[2006].Reliability = DynMap[2001].Reliability + DynMap[2000].Reliability - 1.0;
DynMap[2006].Delta = math.Abs(float64(x)) * DynMap[2000].Delta + math.Abs(float64(y)) * DynMap[2001].Delta + DynMap[2000].Delta*DynMap[2001].Delta;
tempDF0 = x*y;
DynMap[2007].Reliability = DynMap[2006].Reliability + DynMap[2004].Reliability - 1.0;
DynMap[2007].Delta = DynMap[2004].Delta + DynMap[2006].Delta;
tempDF1 = ssXY+tempDF0;
DynMap[2004].Reliability = DynMap[2007].Reliability;
DynMap[2004].Delta = DynMap[2007].Delta;
ssXY = tempDF1;
DynMap[2006].Reliability = DynMap[2000].Reliability;
DynMap[2006].Delta = math.Abs(float64(x)) * DynMap[2000].Delta + math.Abs(float64(x)) * DynMap[2000].Delta + DynMap[2000].Delta*DynMap[2000].Delta;
tempDF0 = x*x;
DynMap[2007].Reliability = DynMap[2005].Reliability + DynMap[2006].Reliability - 1.0;
DynMap[2007].Delta = DynMap[2005].Delta + DynMap[2006].Delta;
tempDF1 = ssXX+tempDF0;
DynMap[2005].Reliability = DynMap[2007].Reliability;
DynMap[2005].Delta = DynMap[2007].Delta;
ssXX = tempDF1;
idx = idx+1;
 }
tempF=convertToFloat(1000);
DynMap[2006].Reliability = DynMap[2002].Reliability;
DynMap[2006].Delta =  DynMap[2002].Delta / math.Abs(tempF);
tempDF0 = mX/tempF;
DynMap[2002].Reliability = DynMap[2006].Reliability;
DynMap[2002].Delta = DynMap[2006].Delta;
mX = tempDF0;
DynMap[2006].Reliability = DynMap[2003].Reliability;
DynMap[2006].Delta =  DynMap[2003].Delta / math.Abs(tempF);
tempDF0 = mY/tempF;
DynMap[2003].Reliability = DynMap[2006].Reliability;
DynMap[2003].Delta = DynMap[2006].Delta;
mY = tempDF0;
DynMap[2006].Reliability = DynMap[2003].Reliability + DynMap[2002].Reliability - 1.0;
DynMap[2006].Delta = math.Abs(float64(mX)) * DynMap[2002].Delta + math.Abs(float64(mY)) * DynMap[2003].Delta + DynMap[2002].Delta*DynMap[2003].Delta;
tempDF0 = mX*mY;
DynMap[2007].Reliability = DynMap[2006].Reliability;
DynMap[2007].Delta = math.Abs(float64(tempF)) * DynMap[2006].Delta;
tempDF1 = tempDF0*tempF;
DynMap[2006].Reliability = DynMap[2007].Reliability + DynMap[2004].Reliability - 1.0;
DynMap[2006].Delta = DynMap[2004].Delta + DynMap[2007].Delta;
tempDF0 = ssXY-tempDF1;
DynMap[2004].Reliability = DynMap[2006].Reliability;
DynMap[2004].Delta = DynMap[2006].Delta;
ssXY = tempDF0;
DynMap[2006].Reliability = DynMap[2002].Reliability;
DynMap[2006].Delta = math.Abs(float64(mX)) * DynMap[2002].Delta + math.Abs(float64(mX)) * DynMap[2002].Delta + DynMap[2002].Delta*DynMap[2002].Delta;
tempDF0 = mX*mX;
DynMap[2007].Reliability = DynMap[2006].Reliability;
DynMap[2007].Delta = math.Abs(float64(tempF)) * DynMap[2006].Delta;
tempDF1 = tempDF0*tempF;
DynMap[2006].Reliability = DynMap[2007].Reliability + DynMap[2005].Reliability - 1.0;
DynMap[2006].Delta = DynMap[2005].Delta + DynMap[2007].Delta;
tempDF0 = ssXX-tempDF1;
DynMap[2005].Reliability = DynMap[2006].Reliability;
DynMap[2005].Delta = DynMap[2006].Delta;
ssXX = tempDF0;
diesel.SendFloat64(mX, tid, 0);
diesel.SendDynVal(DynMap[2002], tid, 0);
diesel.SendFloat64(mY, tid, 0);
diesel.SendDynVal(DynMap[2003], tid, 0);
diesel.SendFloat64(ssXX, tid, 0);
diesel.SendDynVal(DynMap[2005], tid, 0);
diesel.SendFloat64(ssXY, tid, 0);
diesel.SendDynVal(DynMap[2004], tid, 0);

  fmt.Println("Ending thread : ", q);
}

func main() {
  // rand.Seed(time.Now().UTC().UnixNano())
  seed := int64(12345)
  rand.Seed(seed) // deterministic seed for reproducibility

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
}
