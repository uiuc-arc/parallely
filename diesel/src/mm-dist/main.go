package main

import "fmt"
import "math/rand"
import "math"
import "dieseldistacc"

const ArrayDim = 100
const ArraySize = 10000
const SliceSize = 1000
const Iterations = 10
const RowsPerThread = ArrayDim/10
var Num_threads int

var A,B []float64

var Q = []int {1,2,3,4,5,6,7,8,9,10};


func func_0() {
  dieseldistacc.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistacc.WaitForWorkers(Num_threads)
  var DynMap [0]float64;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var mC [10000]float64;
var mA [10000]float32;
var mB [10000]float32;
var slice [1000]float32;
var idx0 int;
var idx1 int;
var tempF64 float64;
var tempDF64 float64;
var tempDF32 float32;
 dieseldistacc.StartTiming() ;
idx0 = 0;
for __temp_0 := 0; __temp_0 < ArraySize; __temp_0++ {
 _temp_index_1 := idx0;
tempF64=A[_temp_index_1];
tempDF64=tempF64;
tempDF32 = float32(tempDF64);
_temp_index_2 := idx0;
mA[_temp_index_2]=tempDF32;
_temp_index_3 := idx0;
tempF64=B[_temp_index_3];
tempDF64=tempF64;
tempDF32 = float32(tempDF64);
_temp_index_4 := idx0;
mB[_temp_index_4]=tempDF32;
idx0 = idx0+1;
 }
for _, q := range(Q) {
 dieseldistacc.SendFloat32Array(mA[:], 0, q);
dieseldistacc.SendFloat32Array(mB[:], 0, q);
 }
idx0 = 0;
for _, q := range(Q) {
 dieseldistacc.ReceiveFloat32Array(slice[:], 0, q);
idx1 = 0;
for __temp_1 := 0; __temp_1 < SliceSize; __temp_1++ {
 _temp_index_5 := idx1;
tempDF32=slice[_temp_index_5];
tempDF64 = float64(tempDF32);
_temp_index_6 := idx0;
mC[_temp_index_6]=tempDF64;
idx0 = idx0+1;
idx1 = idx1+1;
 }
 }
 dieseldistacc.EndTiming() ;


  dieseldistacc.CleanupMain()
  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  dieseldistacc.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistacc.PingMain(tid)
  var DynMap [0]float64;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
var mA [10000]float64;
var mB [10000]float64;
var slice [1000]float64;
var mA32 [10000]float32;
var mB32 [10000]float32;
var slice32 [1000]float32;
var myStartRow int;
var rowIdx int;
var colIdx int;
var innerIdx int;
var outIdx int;
var sum float64;
var tempDF0 float64;
var tempDF1 float64;
var tempDF2 float64;
var tempDF32 float32;
myStartRow = (q-1)*RowsPerThread;
dieseldistacc.ReceiveFloat32Array(mA32[:], tid, 0);
dieseldistacc.ReceiveFloat32Array(mB32[:], tid, 0);
outIdx = 0;
for __temp_2 := 0; __temp_2 < ArraySize; __temp_2++ {
 _temp_index_1 := outIdx;
tempDF32=mA32[_temp_index_1];
tempDF0 = float64(tempDF32);
_temp_index_2 := outIdx;
mA[_temp_index_2]=tempDF0;
_temp_index_3 := outIdx;
tempDF32=mB32[_temp_index_3];
tempDF0 = float64(tempDF32);
_temp_index_4 := outIdx;
mB[_temp_index_4]=tempDF0;
outIdx = outIdx+1;
 }
outIdx = 0;
rowIdx = myStartRow;
for __temp_3 := 0; __temp_3 < RowsPerThread; __temp_3++ {
 colIdx = 0;
for __temp_4 := 0; __temp_4 < ArrayDim; __temp_4++ {
 sum = 0.0;
innerIdx = 0;
for __temp_5 := 0; __temp_5 < ArrayDim; __temp_5++ {
 _temp_index_5 := (rowIdx*ArrayDim)+innerIdx;
tempDF0=mA[_temp_index_5];
_temp_index_6 := (innerIdx*ArrayDim)+colIdx;
tempDF1=mB[_temp_index_6];
tempDF2 = tempDF0*tempDF1;
tempDF0 = sum+tempDF2;
sum = tempDF0;
innerIdx = innerIdx+1;
 }
_temp_index_7 := outIdx;
slice[_temp_index_7]=sum;
colIdx = colIdx+1;
outIdx = outIdx+1;
 }
rowIdx = rowIdx+1;
 }
outIdx = 0;
for __temp_6 := 0; __temp_6 < SliceSize; __temp_6++ {
 _temp_index_8 := outIdx;
tempDF0=slice[_temp_index_8];
tempDF32 = float32(tempDF0);
_temp_index_9 := outIdx;
slice32[_temp_index_9]=tempDF32;
outIdx = outIdx+1;
 }
dieseldistacc.SendFloat32Array(slice32[:], tid, 0);

  fmt.Println("Ending thread : ", q);
}

func main() {
	fmt.Println("Starting main thread");

  Num_threads = 11;
	
  // rand.Seed(time.Now().UTC().UnixNano())
  seed := int64(12345)
  rand.Seed(seed) // deterministic seed for reproducibility

  A = make([]float64, ArraySize)
  B = make([]float64, ArraySize)

  fmt.Println("Generating matrices of size",ArraySize,"using random seed",seed)

  for i := 0; i < ArraySize; i++ {
    A[i] = rand.NormFloat64()*math.Abs(1.0)
    B[i] = rand.NormFloat64()*math.Abs(1.0)
  }

  func_0();


}
