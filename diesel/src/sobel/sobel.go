package main

import "fmt"
import "math/rand"
import "time"
import "math"
import "dieselacc"

var Input_array [10000]float64
var Num_threads int

var Q = []int {1,2,3,4,5,6,7,8,9,10};


func func_0() {
  defer dieselacc.Wg.Done();
  var DynMap [21003]float64;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var output [10000]float64;
dieselacc.InitDynArray(0, 10000, DynMap[:]);
var i int;
var idx int;
var size int;
var perthread int;
var temp1 float64;
var temp2 float32;
DynMap[10000] = 0.0;
var temp3 float32;
DynMap[10001] = 0.0;
var input32 [10000]float32;
dieselacc.InitDynArray(10002, 10000, DynMap[:]);
var slice [1000]float32;
dieselacc.InitDynArray(20002, 1000, DynMap[:]);
var elem float64;
DynMap[21002] = 0.0;
i = 0;
size = 10000;
perthread = 1000;
for __temp_0 := 0; __temp_0 < size; __temp_0++ {
 _temp_index_1 := i;
temp1=Input_array[_temp_index_1];
temp2 = float32(temp1);
DynMap[10000] = dieselacc.GetCastingError64to32(temp1, temp2);
_temp_index_2 := i;
input32[_temp_index_2]=temp2;
DynMap[10002 + _temp_index_2] = DynMap[10000];
i = i+1;
 }
for _, q := range(Q) {
 dieselacc.SendDynFloat32Array(input32[:], 0, q, DynMap[:], 10002);
 }
for _, q := range(Q) {
 dieselacc.ReceiveDynFloat32Array(slice[:], 0, q, DynMap[:], 20002);
i = 0;
for __temp_1 := 0; __temp_1 < perthread; __temp_1++ {
 idx = (q-1)*100+i;
_temp_index_3 := i;
temp3=slice[_temp_index_3];
DynMap[10001] = DynMap[20002 + _temp_index_3];
elem = float64(temp3);
DynMap[21002] = DynMap[10001];
_temp_index_4 := idx;
output[_temp_index_4]=elem;
DynMap[0 + _temp_index_4] = DynMap[21002];
i = i+1;
 }
 }


  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  defer dieselacc.Wg.Done();
  var DynMap [22010]float64;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
var array [10000]float64;
dieselacc.InitDynArray(0, 10000, DynMap[:]);
var slice [1000]float64;
dieselacc.InitDynArray(10000, 1000, DynMap[:]);
var array32 [10000]float32;
dieselacc.InitDynArray(11000, 10000, DynMap[:]);
var slice32 [1000]float32;
dieselacc.InitDynArray(21000, 1000, DynMap[:]);
var i int;
var j int;
var c1 int;
var c2 int;
var k int;
var cond1 int;
var cond2 int;
var conditional int;
var point float64;
DynMap[22000] = 0.0;
var temp1 float64;
DynMap[22001] = 0.0;
var temp2 float64;
DynMap[22002] = 0.0;
var temp3 float64;
DynMap[22003] = 0.0;
var temp4 float64;
DynMap[22004] = 0.0;
var temp5 float64;
DynMap[22005] = 0.0;
var temp6 float64;
DynMap[22006] = 0.0;
var temp7 float64;
DynMap[22007] = 0.0;
var temp8 float64;
DynMap[22008] = 0.0;
var temp32 float32;
DynMap[22009] = 0.0;
dieselacc.ReceiveDynFloat32Array(array32[:], tid, 0, DynMap[:], 11000);
i = 0;
for __temp_2 := 0; __temp_2 < 10000; __temp_2++ {
 _temp_index_1 := i;
temp32=array32[_temp_index_1];
DynMap[22009] = DynMap[11000 + _temp_index_1];
temp1 = float64(temp32);
DynMap[22001] = DynMap[22009];
_temp_index_2 := i;
array[_temp_index_2]=temp1;
DynMap[0 + _temp_index_2] = DynMap[22001];
i = i+1;
 }
i = (q-1)*10;
k = 0;
c1 = 10;
c2 = 98;
for __temp_3 := 0; __temp_3 < c1; __temp_3++ {
 j = 1;
cond1 = dieselacc.ConvBool(i<99);
cond2 = dieselacc.ConvBool(i>0);
conditional = dieselacc.ConvBool(cond1==1 && cond2==1);
if conditional != 0 {
 for __temp_4 := 0; __temp_4 < c2; __temp_4++ {
 _temp_index_3 := i*100+j-101;
temp1=array[_temp_index_3];
DynMap[22001] = DynMap[0 + _temp_index_3];
_temp_index_4 := i*100+j-100;
temp2=array[_temp_index_4];
DynMap[22002] = DynMap[0 + _temp_index_4];
_temp_index_5 := i*100+j-99;
temp3=array[_temp_index_5];
DynMap[22003] = DynMap[0 + _temp_index_5];
_temp_index_6 := i*100+j+99;
temp4=array[_temp_index_6];
DynMap[22004] = DynMap[0 + _temp_index_6];
_temp_index_7 := i*100+j+100;
temp5=array[_temp_index_7];
DynMap[22005] = DynMap[0 + _temp_index_7];
_temp_index_8 := i*100+j+101;
temp6=array[_temp_index_8];
DynMap[22006] = DynMap[0 + _temp_index_8];
DynMap[22007] = DynMap[22002] + DynMap[22002];
temp7 = temp2+temp2;
DynMap[22008] = DynMap[22005] + DynMap[22005];
temp8 = temp5+temp5;
DynMap[22000] = DynMap[22001] + DynMap[22007];
point = temp1+temp7;
DynMap[22000] = DynMap[22000] + DynMap[22003];
point = point+temp3;
DynMap[22000] = DynMap[22000] + DynMap[22004];
point = point-temp4;
DynMap[22000] = DynMap[22000] + DynMap[22008];
point = point-temp8;
DynMap[22000] = DynMap[22000] + DynMap[22006];
point = point-temp6;
_temp_index_9 := k;
slice[_temp_index_9]=point;
DynMap[10000 + _temp_index_9] = DynMap[22000];
j = j+1;
k = k+1;
 }
 } else {
 for __temp_5 := 0; __temp_5 < c2; __temp_5++ {
 _temp_index_10 := i*100+j;
temp7=array[_temp_index_10];
DynMap[22007] = DynMap[0 + _temp_index_10];
_temp_index_11 := k;
slice[_temp_index_11]=temp7;
DynMap[10000 + _temp_index_11] = DynMap[22007];
j = j+1;
k = k+1;
 }
 }
i = i+1;
 }
i = 0;
for __temp_6 := 0; __temp_6 < 1000; __temp_6++ {
 _temp_index_12 := i;
temp1=slice[_temp_index_12];
DynMap[22001] = DynMap[10000 + _temp_index_12];
temp32 = float32(temp1);
DynMap[22009] = dieselacc.GetCastingError64to32(temp1, temp32);
_temp_index_13 := i;
slice32[_temp_index_13]=temp32;
DynMap[21000 + _temp_index_13] = DynMap[22009];
i = i+1;
 }
dieselacc.SendDynFloat32Array(slice32[:], tid, 0, DynMap[:], 21000);

  fmt.Println("Ending thread : ", q);
}

func main() {
	fmt.Println("Starting main thread");

  Num_threads = 11;
	
	dieselacc.InitChannels(11);

  // Using math becasue GoLang is crazy
  randSource := rand.NewSource(int64(math.Abs(1)))
  randGen := rand.New(randSource)

  for i:=0; i<10000; i++ {
    Input_array[i] = randGen.Float64()
  }

  startTime := time.Now()

  go func_0();
for _, index := range Q {
go func_Q(index);
}


  fmt.Println("Main thread waiting for others to finish");  
	dieselacc.Wg.Wait()
  elapsed := time.Since(startTime)
  
  fmt.Println("Elapsed time : ", elapsed.Nanoseconds())
}