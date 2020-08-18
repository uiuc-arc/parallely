package main

import (
  "fmt"
  "math"
  "math/rand"
  "dieseldistacc"
)

const ArrayDim = 100
const InnerDim = ArrayDim-2
const ArraySize = 10000
const SliceSize = 1000
const Iterations = 10
const RowsPerThread = ArrayDim/10
const Num_threads = 11

var Input [10000]float64

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

var Q = []int {1,2,3,4,5,6,7,8,9,10};


func func_0() {
  dieseldistacc.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistacc.WaitForWorkers(Num_threads)
  var DynMap [21002]dieseldistacc.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var array [10000]float64;
dieseldistacc.InitDynArray(0, 10000, DynMap[:]);
var array32 [10000]float32;
dieseldistacc.InitDynArray(10000, 10000, DynMap[:]);
var slice [1000]float32;
dieseldistacc.InitDynArray(20000, 1000, DynMap[:]);
var idx0 int;
var idx1 int;
var tempDF64 float64;
DynMap[21000] = dieseldistacc.ProbInterval{0};
var tempDF32 float32;
DynMap[21001] = dieseldistacc.ProbInterval{0};
array=Input;
dieseldistacc.InitDynArray(0, 10000, DynMap[:]);
 dieseldistacc.StartTiming() ;
idx0 = 0;
for __temp_0 := 0; __temp_0 < ArraySize; __temp_0++ {
 _temp_index_1 := idx0;
tempDF64=array[_temp_index_1];
DynMap[21000] = DynMap[0 + _temp_index_1];
tempDF32 = float32(tempDF64);
DynMap[21001].Delta = dieseldistacc.GetCastingError64to32(tempDF64, tempDF32);
_temp_index_2 := idx0;
array32[_temp_index_2]=tempDF32;
DynMap[10000 + _temp_index_2] = DynMap[21001];
idx0 = idx0+1;
 }
for __temp_1 := 0; __temp_1 < Iterations; __temp_1++ {
 for _, q := range(Q) {
 dieseldistacc.SendDynFloat32ArrayO1(array32[:], 0, q, DynMap[:], 10000);
 }
idx0 = 0;
for _, q := range(Q) {
 dieseldistacc.ReceiveDynFloat32ArrayO1(slice[:], 0, q, DynMap[:], 20000);
idx1 = 0;
for __temp_2 := 0; __temp_2 < SliceSize; __temp_2++ {
 _temp_index_3 := idx1;
tempDF32=slice[_temp_index_3];
DynMap[21001] = DynMap[20000 + _temp_index_3];
_temp_index_4 := idx0;
array32[_temp_index_4]=tempDF32;
DynMap[10000 + _temp_index_4] = DynMap[21001];
idx0 = idx0+1;
idx1 = idx1+1;
 }
 }
 }
idx0 = 0;
for __temp_3 := 0; __temp_3 < ArraySize; __temp_3++ {
 _temp_index_5 := idx0;
tempDF32=array32[_temp_index_5];
DynMap[21001] = DynMap[10000 + _temp_index_5];
tempDF64 = float64(tempDF32);
DynMap[21000] = DynMap[21001];
_temp_index_6 := idx0;
array[_temp_index_6]=tempDF64;
DynMap[0 + _temp_index_6] = DynMap[21000];
idx0 = idx0+1;
 }
 dieseldistacc.EndTiming() ;


  dieseldistacc.CleanupMain()
  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  dieseldistacc.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistacc.PingMain(tid)
  var DynMap [22005]dieseldistacc.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
var array [10000]float64;
dieseldistacc.InitDynArray(0, 10000, DynMap[:]);
var slice [1000]float64;
dieseldistacc.InitDynArray(10000, 1000, DynMap[:]);
var array32 [10000]float32;
dieseldistacc.InitDynArray(11000, 10000, DynMap[:]);
var slice32 [1000]float32;
dieseldistacc.InitDynArray(21000, 1000, DynMap[:]);
var idx0 int;
var tempDF64 float64;
DynMap[22000] = dieseldistacc.ProbInterval{0};
var tempDF32 float32;
DynMap[22001] = dieseldistacc.ProbInterval{0};
var myStartRow int;
var rowIdx int;
var colIdx int;
var outIdx int;
var firstRow int;
var lastRow int;
var edgeRow int;
var tempDF0 float64;
DynMap[22002] = dieseldistacc.ProbInterval{0};
var tempDF1 float64;
DynMap[22003] = dieseldistacc.ProbInterval{0};
var tempDF2 float64;
DynMap[22004] = dieseldistacc.ProbInterval{0};
myStartRow = (q-1)*RowsPerThread;
for __temp_4 := 0; __temp_4 < Iterations; __temp_4++ {
 dieseldistacc.ReceiveDynFloat32ArrayO1(array32[:], tid, 0, DynMap[:], 11000);
idx0 = 0;
for __temp_5 := 0; __temp_5 < ArraySize; __temp_5++ {
 _temp_index_1 := idx0;
tempDF32=array32[_temp_index_1];
DynMap[22001] = DynMap[11000 + _temp_index_1];
tempDF64 = float64(tempDF32);
DynMap[22000] = DynMap[22001];
_temp_index_2 := idx0;
array[_temp_index_2]=tempDF64;
DynMap[0 + _temp_index_2] = DynMap[22000];
idx0 = idx0+1;
 }
rowIdx = myStartRow;
outIdx = 0;
for __temp_6 := 0; __temp_6 < RowsPerThread; __temp_6++ {
 colIdx = 0;
firstRow = dieseldistacc.ConvBool(rowIdx==0);
lastRow = dieseldistacc.ConvBool(rowIdx==(ArrayDim-1));
edgeRow = dieseldistacc.ConvBool(firstRow==1 || lastRow==1);
if edgeRow != 0 {
 for __temp_7 := 0; __temp_7 < ArrayDim; __temp_7++ {
 _temp_index_3 := (rowIdx*ArrayDim)+colIdx;
tempDF0=array[_temp_index_3];
DynMap[22002] = DynMap[0 + _temp_index_3];
_temp_index_4 := outIdx;
slice[_temp_index_4]=tempDF0;
DynMap[10000 + _temp_index_4] = DynMap[22002];
colIdx = colIdx+1;
outIdx = outIdx+1;
 }
 } else {
 _temp_index_5 := rowIdx*ArrayDim;
tempDF0=array[_temp_index_5];
DynMap[22002] = DynMap[0 + _temp_index_5];
_temp_index_6 := outIdx;
slice[_temp_index_6]=tempDF0;
DynMap[10000 + _temp_index_6] = DynMap[22002];
outIdx = outIdx+1;
colIdx = 1;
for __temp_8 := 0; __temp_8 < InnerDim; __temp_8++ {
 _temp_index_7 := ((rowIdx-1)*ArrayDim)+colIdx;
tempDF0=array[_temp_index_7];
DynMap[22002] = DynMap[0 + _temp_index_7];
_temp_index_8 := ((rowIdx+1)*ArrayDim)+colIdx;
tempDF1=array[_temp_index_8];
DynMap[22003] = DynMap[0 + _temp_index_8];
DynMap[22004].Delta = DynMap[22002].Delta + DynMap[22003].Delta;
tempDF2 = tempDF0+tempDF1;
_temp_index_9 := (rowIdx*ArrayDim)+colIdx-1;
tempDF0=array[_temp_index_9];
DynMap[22002] = DynMap[0 + _temp_index_9];
DynMap[22003].Delta = DynMap[22004].Delta + DynMap[22002].Delta;
tempDF1 = tempDF2+tempDF0;
_temp_index_10 := (rowIdx*ArrayDim)+colIdx+1;
tempDF2=array[_temp_index_10];
DynMap[22004] = DynMap[0 + _temp_index_10];
DynMap[22002].Delta = DynMap[22003].Delta + DynMap[22004].Delta;
tempDF0 = tempDF1+tempDF2;
_temp_index_11 := (rowIdx*ArrayDim)+colIdx;
tempDF1=array[_temp_index_11];
DynMap[22003] = DynMap[0 + _temp_index_11];
DynMap[22004].Delta = DynMap[22002].Delta + DynMap[22003].Delta;
tempDF2 = tempDF0+tempDF1;
DynMap[22002].Delta =  DynMap[22004].Delta / math.Abs(5.0);
tempDF0 = tempDF2/5.0;
_temp_index_12 := outIdx;
slice[_temp_index_12]=tempDF0;
DynMap[10000 + _temp_index_12] = DynMap[22002];
colIdx = colIdx+1;
outIdx = outIdx+1;
 }
_temp_index_13 := (rowIdx*ArrayDim)+colIdx;
tempDF0=array[_temp_index_13];
DynMap[22002] = DynMap[0 + _temp_index_13];
_temp_index_14 := outIdx;
slice[_temp_index_14]=tempDF0;
DynMap[10000 + _temp_index_14] = DynMap[22002];
outIdx = outIdx+1;
 }
rowIdx = rowIdx+1;
 }
idx0 = 0;
for __temp_9 := 0; __temp_9 < SliceSize; __temp_9++ {
 _temp_index_15 := idx0;
tempDF64=slice[_temp_index_15];
DynMap[22000] = DynMap[10000 + _temp_index_15];
tempDF32 = float32(tempDF64);
DynMap[22001].Delta = dieseldistacc.GetCastingError64to32(tempDF64, tempDF32);
_temp_index_16 := idx0;
slice32[_temp_index_16]=tempDF32;
DynMap[21000 + _temp_index_16] = DynMap[22001];
idx0 = idx0+1;
 }
dieseldistacc.SendDynFloat32ArrayO1(slice32[:], tid, 0, DynMap[:], 21000);
 }

  fmt.Println("Ending thread : ", q);
}

func main() {
  // rand.Seed(time.Now().UTC().UnixNano())
  seed := int64(12345)
  rand.Seed(seed) // deterministic seed for reproducibility

  fmt.Println("Generating array of size",ArraySize,"using random seed",seed)

  for i := 0; i < ArraySize; i++ {
    Input[i] = rand.NormFloat64()*math.Abs(1.0)
  }

  fmt.Println("Starting program");

  func_0();


}
