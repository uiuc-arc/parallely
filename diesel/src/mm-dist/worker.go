package main

import (
  "fmt"
  "math"
  "dieseldistacc"
  "os"
  "strconv"
)

const ArrayDim = 100
const ArraySize = 10000
const SliceSize = 1000
const Iterations = 10
const RowsPerThread = ArrayDim/10
var Num_threads int

var A,B []float64

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
  var DynMap [52002]float64;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var mA32 [10000]float32;
dieseldistacc.InitDynArray(0, 10000, DynMap[:]);
var mB32 [10000]float32;
dieseldistacc.InitDynArray(10000, 10000, DynMap[:]);
var slice32 [1000]float32;
dieseldistacc.InitDynArray(20000, 1000, DynMap[:]);
var mA64 [10000]float64;
dieseldistacc.InitDynArray(21000, 10000, DynMap[:]);
var mB64 [10000]float64;
dieseldistacc.InitDynArray(31000, 10000, DynMap[:]);
var mC64 [10000]float64;
dieseldistacc.InitDynArray(41000, 10000, DynMap[:]);
var slice64 [1000]float64;
dieseldistacc.InitDynArray(51000, 1000, DynMap[:]);
var idx0 int;
var idx1 int;
var itemCheck int;
var arrayCheck int;
var tempF64 float64;
var tempDF64 float64;
DynMap[52000] = 0.0;
var tempDF32 float32;
DynMap[52001] = 0.0;
 dieseldistacc.StartTiming() ;
arrayCheck = 1;
idx0 = 0;
for __temp_0 := 0; __temp_0 < ArraySize; __temp_0++ {
 _temp_index_1 := idx0;
tempF64=A[_temp_index_1];
tempDF64=tempF64;
DynMap[52000] = 0.0;
_temp_index_2 := idx0;
mA64[_temp_index_2]=tempDF64;
DynMap[21000 + _temp_index_2] = DynMap[52000];
tempDF32 = float32(tempDF64);
DynMap[52001] = dieseldistacc.GetCastingError64to32(tempDF64, tempDF32);
itemCheck = dieseldistacc.Check(DynMap[52001], 1.0, 0.000001);
arrayCheck = arrayCheck*itemCheck;
_temp_index_3 := idx0;
mA32[_temp_index_3]=tempDF32;
DynMap[0 + _temp_index_3] = DynMap[52001];
_temp_index_4 := idx0;
tempF64=B[_temp_index_4];
tempDF64=tempF64;
DynMap[52000] = 0.0;
_temp_index_5 := idx0;
mB64[_temp_index_5]=tempDF64;
DynMap[31000 + _temp_index_5] = DynMap[52000];
tempDF32 = float32(tempDF64);
DynMap[52001] = dieseldistacc.GetCastingError64to32(tempDF64, tempDF32);
itemCheck = dieseldistacc.Check(DynMap[52001], 1.0, 0.000001);
arrayCheck = arrayCheck*itemCheck;
_temp_index_6 := idx0;
mB32[_temp_index_6]=tempDF32;
DynMap[10000 + _temp_index_6] = DynMap[52001];
idx0 = idx0+1;
 }
 fmt.Println("Master array send using FP32:",arrayCheck) ;
if arrayCheck != 0 {
 for _, q := range(Q) {
 dieseldistacc.SendInt(arrayCheck, 0, q);
dieseldistacc.SendDynFloat32ArrayO1(mA32[:], 0, q, DynMap[:], 0);
dieseldistacc.SendDynFloat32ArrayO1(mB32[:], 0, q, DynMap[:], 10000);
 }
 } else {
 for _, q := range(Q) {
 dieseldistacc.SendInt(arrayCheck, 0, q);
dieseldistacc.SendDynFloat64ArrayO1(mA64[:], 0, q, DynMap[:], 21000);
dieseldistacc.SendDynFloat64ArrayO1(mB64[:], 0, q, DynMap[:], 31000);
 }
 }
idx0 = 0;
for _, q := range(Q) {
 dieseldistacc.ReceiveInt(&arrayCheck, 0, q);
if arrayCheck != 0 {
 dieseldistacc.ReceiveDynFloat32ArrayO1(slice32[:], 0, q, DynMap[:], 20000);
idx1 = 0;
for __temp_1 := 0; __temp_1 < SliceSize; __temp_1++ {
 _temp_index_7 := idx1;
tempDF32=slice32[_temp_index_7];
DynMap[52001] = DynMap[20000 + _temp_index_7];
tempDF64 = float64(tempDF32);
DynMap[52000] = DynMap[52001];
_temp_index_8 := idx0;
mC64[_temp_index_8]=tempDF64;
DynMap[41000 + _temp_index_8] = DynMap[52000];
idx0 = idx0+1;
idx1 = idx1+1;
 }
 } else {
 dieseldistacc.ReceiveDynFloat64ArrayO1(slice64[:], 0, q, DynMap[:], 51000);
idx1 = 0;
for __temp_2 := 0; __temp_2 < SliceSize; __temp_2++ {
 _temp_index_9 := idx1;
tempDF64=slice64[_temp_index_9];
DynMap[52000] = DynMap[51000 + _temp_index_9];
_temp_index_10 := idx0;
mC64[_temp_index_10]=tempDF64;
DynMap[41000 + _temp_index_10] = DynMap[52000];
idx0 = idx0+1;
idx1 = idx1+1;
 }
 }
 }
 dieseldistacc.EndTiming() ;


  // dieseldistacc.CleanupMain()
  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) { 
  dieseldistacc.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistacc.PingMain(tid)
  var DynMap [42005]float64;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
var mA32 [10000]float32;
dieseldistacc.InitDynArray(0, 10000, DynMap[:]);
var mB32 [10000]float32;
dieseldistacc.InitDynArray(10000, 10000, DynMap[:]);
var slice32 [1000]float32;
dieseldistacc.InitDynArray(20000, 1000, DynMap[:]);
var mA64 [10000]float64;
dieseldistacc.InitDynArray(21000, 10000, DynMap[:]);
var mB64 [10000]float64;
dieseldistacc.InitDynArray(31000, 10000, DynMap[:]);
var slice64 [1000]float64;
dieseldistacc.InitDynArray(41000, 1000, DynMap[:]);
var myStartRow int;
var rowIdx int;
var colIdx int;
var innerIdx int;
var outIdx int;
var itemCheck int;
var arrayCheck int;
var sum float64;
DynMap[42000] = 0.0;
var tempDF0 float64;
DynMap[42001] = 0.0;
var tempDF1 float64;
DynMap[42002] = 0.0;
var tempDF2 float64;
DynMap[42003] = 0.0;
var tempDF32 float32;
DynMap[42004] = 0.0;
myStartRow = (q-1)*RowsPerThread;
dieseldistacc.ReceiveInt(&arrayCheck, tid, 0);
if arrayCheck != 0 {
 dieseldistacc.ReceiveDynFloat32ArrayO1(mA32[:], tid, 0, DynMap[:], 0);
dieseldistacc.ReceiveDynFloat32ArrayO1(mB32[:], tid, 0, DynMap[:], 10000);
outIdx = 0;
for __temp_3 := 0; __temp_3 < ArraySize; __temp_3++ {
 _temp_index_1 := outIdx;
tempDF32=mA32[_temp_index_1];
DynMap[42004] = DynMap[0 + _temp_index_1];
tempDF0 = float64(tempDF32);
DynMap[42001] = DynMap[42004];
_temp_index_2 := outIdx;
mA64[_temp_index_2]=tempDF0;
DynMap[21000 + _temp_index_2] = DynMap[42001];
_temp_index_3 := outIdx;
tempDF32=mB32[_temp_index_3];
DynMap[42004] = DynMap[10000 + _temp_index_3];
tempDF0 = float64(tempDF32);
DynMap[42001] = DynMap[42004];
_temp_index_4 := outIdx;
mB64[_temp_index_4]=tempDF0;
DynMap[31000 + _temp_index_4] = DynMap[42001];
outIdx = outIdx+1;
 }
 } else {
 dieseldistacc.ReceiveDynFloat64ArrayO1(mA64[:], tid, 0, DynMap[:], 21000);
dieseldistacc.ReceiveDynFloat64ArrayO1(mB64[:], tid, 0, DynMap[:], 31000);
 }
outIdx = 0;
rowIdx = myStartRow;
for __temp_4 := 0; __temp_4 < RowsPerThread; __temp_4++ {
 colIdx = 0;
for __temp_5 := 0; __temp_5 < ArrayDim; __temp_5++ {
 DynMap[42000] = 0.0;
sum = 0.0;
innerIdx = 0;
for __temp_6 := 0; __temp_6 < ArrayDim; __temp_6++ {
 _temp_index_5 := (rowIdx*ArrayDim)+innerIdx;
tempDF0=mA64[_temp_index_5];
DynMap[42001] = DynMap[21000 + _temp_index_5];
_temp_index_6 := (innerIdx*ArrayDim)+colIdx;
tempDF1=mB64[_temp_index_6];
DynMap[42002] = DynMap[31000 + _temp_index_6];
DynMap[42003] = math.Abs(float64(tempDF0)) * DynMap[42001] + math.Abs(float64(tempDF1)) * DynMap[42002] + DynMap[42001]*DynMap[42002];
tempDF2 = tempDF0*tempDF1;
DynMap[42001] = DynMap[42000] + DynMap[42003];
tempDF0 = sum+tempDF2;
DynMap[42000] = DynMap[42001];
DynMap[42000] = DynMap[42001];
sum = tempDF0;
innerIdx = innerIdx+1;
 }
_temp_index_7 := outIdx;
slice64[_temp_index_7]=sum;
DynMap[41000 + _temp_index_7] = DynMap[42000];
colIdx = colIdx+1;
outIdx = outIdx+1;
 }
rowIdx = rowIdx+1;
 }
arrayCheck = 1;
outIdx = 0;
for __temp_7 := 0; __temp_7 < SliceSize; __temp_7++ {
 _temp_index_8 := outIdx;
tempDF0=slice64[_temp_index_8];
DynMap[42001] = DynMap[41000 + _temp_index_8];
tempDF32 = float32(tempDF0);
DynMap[42004] = dieseldistacc.GetCastingError64to32(tempDF0, tempDF32);
itemCheck = dieseldistacc.Check(DynMap[42004], 1.0, 0.0000015);
arrayCheck = arrayCheck*itemCheck;
_temp_index_9 := outIdx;
slice32[_temp_index_9]=tempDF32;
DynMap[20000 + _temp_index_9] = DynMap[42004];
outIdx = outIdx+1;
 }
 fmt.Println("Worker",q,"slice send using FP32:",arrayCheck) ;
dieseldistacc.SendInt(arrayCheck, tid, 0);
if arrayCheck != 0 {
 dieseldistacc.SendDynFloat32ArrayO1(slice32[:], tid, 0, DynMap[:], 20000);
 } else {
 dieseldistacc.SendDynFloat64ArrayO1(slice64[:], tid, 0, DynMap[:], 41000);
 }

  fmt.Println("Ending thread : ", q);
}

func main() {
  // rand.Seed(time.Now().UTC().UnixNano())
  tid, _ := strconv.Atoi(os.Args[1])
  _ = math.Inf
    
  Num_threads = 11;
	
  func_Q(tid)

}
