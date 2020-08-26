package main

import (
  "fmt"
  "math"
  "strconv"
  "os"
  "dieseldistacc"
)

const ArrayDim = 200
const InnerDim = ArrayDim-2
const ArraySize = 40000
const SliceSize = 4000
const Iterations = 10
const RowsPerThread = ArrayDim/10
const Num_threads = 11;

var Input [40000]float64

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
  var DynMap [84002]float64;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var array [40000]float64;
dieseldistacc.InitDynArray(0, 40000, DynMap[:]);
var array32 [40000]float32;
dieseldistacc.InitDynArray(40000, 40000, DynMap[:]);
var slice [4000]float32;
dieseldistacc.InitDynArray(80000, 4000, DynMap[:]);
var idx0 int;
var idx1 int;
var tempDF64 float64;
DynMap[84000] = 0;
var tempDF32 float32;
DynMap[84001] = 0;
array=Input;
dieseldistacc.InitDynArray(0, 40000, DynMap[:]);
 dieseldistacc.StartTiming() ;
idx0 = 0;
for __temp_0 := 0; __temp_0 < ArraySize; __temp_0++ {
 _temp_index_1 := idx0;
tempDF64=array[_temp_index_1];
DynMap[84000] = DynMap[0 + _temp_index_1];
tempDF32 = float32(tempDF64);
DynMap[84001] = dieseldistacc.GetCastingError64to32(tempDF64, tempDF32);
_temp_index_2 := idx0;
array32[_temp_index_2]=tempDF32;
DynMap[40000 + _temp_index_2] = DynMap[84001];
idx0 = idx0+1;
 }
for __temp_1 := 0; __temp_1 < Iterations; __temp_1++ {
 for _, q := range(Q) {
 dieseldistacc.SendDynFloat32ArrayO1(array32[:], 0, q, DynMap[:], 40000);
 }
idx0 = 0;
for _, q := range(Q) {
 dieseldistacc.ReceiveDynFloat32ArrayO1(slice[:], 0, q, DynMap[:], 80000);
idx1 = 0;
for __temp_2 := 0; __temp_2 < SliceSize; __temp_2++ {
 _temp_index_3 := idx1;
tempDF32=slice[_temp_index_3];
DynMap[84001] = DynMap[80000 + _temp_index_3];
_temp_index_4 := idx0;
array32[_temp_index_4]=tempDF32;
DynMap[40000 + _temp_index_4] = DynMap[84001];
idx0 = idx0+1;
idx1 = idx1+1;
 }
 }
 }
idx0 = 0;
for __temp_3 := 0; __temp_3 < ArraySize; __temp_3++ {
 _temp_index_5 := idx0;
tempDF32=array32[_temp_index_5];
DynMap[84001] = DynMap[40000 + _temp_index_5];
tempDF64 = float64(tempDF32);
DynMap[84000] = DynMap[84001];
_temp_index_6 := idx0;
array[_temp_index_6]=tempDF64;
DynMap[0 + _temp_index_6] = DynMap[84000];
idx0 = idx0+1;
 }
 dieseldistacc.EndTiming() ;


  dieseldistacc.CleanupMain()
  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  dieseldistacc.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistacc.PingMain(tid)
  var DynMap [88005]float64;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
var array [40000]float64;
dieseldistacc.InitDynArray(0, 40000, DynMap[:]);
var slice [4000]float64;
dieseldistacc.InitDynArray(40000, 4000, DynMap[:]);
var array32 [40000]float32;
dieseldistacc.InitDynArray(44000, 40000, DynMap[:]);
var slice32 [4000]float32;
dieseldistacc.InitDynArray(84000, 4000, DynMap[:]);
var idx0 int;
var tempDF64 float64;
DynMap[88000] = 0;
var tempDF32 float32;
DynMap[88001] = 0;
var myStartRow int;
var rowIdx int;
var colIdx int;
var outIdx int;
var firstRow int;
var lastRow int;
var edgeRow int;
var tempDF0 float64;
DynMap[88002] = 0;
var tempDF1 float64;
DynMap[88003] = 0;
var tempDF2 float64;
DynMap[88004] = 0;
myStartRow = (q-1)*RowsPerThread;
for __temp_4 := 0; __temp_4 < Iterations; __temp_4++ {
 dieseldistacc.ReceiveDynFloat32ArrayO1(array32[:], tid, 0, DynMap[:], 44000);
idx0 = 0;
for __temp_5 := 0; __temp_5 < ArraySize; __temp_5++ {
 _temp_index_1 := idx0;
tempDF32=array32[_temp_index_1];
DynMap[88001] = DynMap[44000 + _temp_index_1];
tempDF64 = float64(tempDF32);
DynMap[88000] = DynMap[88001];
_temp_index_2 := idx0;
array[_temp_index_2]=tempDF64;
DynMap[0 + _temp_index_2] = DynMap[88000];
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
DynMap[88002] = DynMap[0 + _temp_index_3];
_temp_index_4 := outIdx;
slice[_temp_index_4]=tempDF0;
DynMap[40000 + _temp_index_4] = DynMap[88002];
colIdx = colIdx+1;
outIdx = outIdx+1;
 }
 } else {
 _temp_index_5 := rowIdx*ArrayDim;
tempDF0=array[_temp_index_5];
DynMap[88002] = DynMap[0 + _temp_index_5];
_temp_index_6 := outIdx;
slice[_temp_index_6]=tempDF0;
DynMap[40000 + _temp_index_6] = DynMap[88002];
outIdx = outIdx+1;
colIdx = 1;
for __temp_8 := 0; __temp_8 < InnerDim; __temp_8++ {
 _temp_index_7 := ((rowIdx-1)*ArrayDim)+colIdx;
tempDF0=array[_temp_index_7];
DynMap[88002] = DynMap[0 + _temp_index_7];
_temp_index_8 := ((rowIdx+1)*ArrayDim)+colIdx;
tempDF1=array[_temp_index_8];
DynMap[88003] = DynMap[0 + _temp_index_8];
DynMap[88004] = DynMap[88002] + DynMap[88003];
tempDF2 = tempDF0+tempDF1;
_temp_index_9 := (rowIdx*ArrayDim)+colIdx-1;
tempDF0=array[_temp_index_9];
DynMap[88002] = DynMap[0 + _temp_index_9];
DynMap[88003] = DynMap[88004] + DynMap[88002];
tempDF1 = tempDF2+tempDF0;
_temp_index_10 := (rowIdx*ArrayDim)+colIdx+1;
tempDF2=array[_temp_index_10];
DynMap[88004] = DynMap[0 + _temp_index_10];
DynMap[88002] = DynMap[88003] + DynMap[88004];
tempDF0 = tempDF1+tempDF2;
_temp_index_11 := (rowIdx*ArrayDim)+colIdx;
tempDF1=array[_temp_index_11];
DynMap[88003] = DynMap[0 + _temp_index_11];
DynMap[88004] = DynMap[88002] + DynMap[88003];
tempDF2 = tempDF0+tempDF1;
DynMap[88002] =  DynMap[88004] / math.Abs(5.0);
tempDF0 = tempDF2/5.0;
_temp_index_12 := outIdx;
slice[_temp_index_12]=tempDF0;
DynMap[40000 + _temp_index_12] = DynMap[88002];
colIdx = colIdx+1;
outIdx = outIdx+1;
 }
_temp_index_13 := (rowIdx*ArrayDim)+colIdx;
tempDF0=array[_temp_index_13];
DynMap[88002] = DynMap[0 + _temp_index_13];
_temp_index_14 := outIdx;
slice[_temp_index_14]=tempDF0;
DynMap[40000 + _temp_index_14] = DynMap[88002];
outIdx = outIdx+1;
 }
rowIdx = rowIdx+1;
 }
idx0 = 0;
for __temp_9 := 0; __temp_9 < SliceSize; __temp_9++ {
 _temp_index_15 := idx0;
tempDF64=slice[_temp_index_15];
DynMap[88000] = DynMap[40000 + _temp_index_15];
tempDF32 = float32(tempDF64);
DynMap[88001] = dieseldistacc.GetCastingError64to32(tempDF64, tempDF32);
_temp_index_16 := idx0;
slice32[_temp_index_16]=tempDF32;
DynMap[84000 + _temp_index_16] = DynMap[88001];
idx0 = idx0+1;
 }
dieseldistacc.SendDynFloat32ArrayO1(slice32[:], tid, 0, DynMap[:], 84000);
 }

  fmt.Println("Ending thread : ", q);
}

func main() {
	tid, _ := strconv.Atoi(os.Args[1])

  fmt.Println("Starting program");
 
  // using this s#!t because go complains when math not used
  _ = math.Inf
  
  func_Q(tid)

}
