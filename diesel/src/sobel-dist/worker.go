package main

import "fmt"
import "math"
import "dieseldist"
import "os"
import "strconv"

var Input_array [10000]float64
var Num_threads int

var Q = []int {1,2,3,4,5,6,7,8,9,10};


func func_0() {
  dieseldist.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldist.WaitForWorkers(Num_threads)
  var DynMap [0]dieseldist.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var output [10000]float64;
var i int;
var idx int;
var size int;
var perthread int;
var temp1 float64;
var temp2 float32;
var temp3 float32;
var input32 [10000]float32;
var slice [1000]float32;
var elem float64;
i = 0;
size = 10000;
perthread = 1000;
 dieseldist.StartTiming() ;
for __temp_0 := 0; __temp_0 < size; __temp_0++ {
 _temp_index_1 := i;
temp1=Input_array[_temp_index_1];
temp2 = float32(temp1);
_temp_index_2 := i;
input32[_temp_index_2]=temp2;
i = i+1;
 }
for _, q := range(Q) {
 dieseldist.SendFloat32Array(input32[:], 0, q);
 }
for _, q := range(Q) {
 dieseldist.ReceiveFloat32Array(slice[:], 0, q);
i = 0;
for __temp_1 := 0; __temp_1 < perthread; __temp_1++ {
 idx = (q-1)*100+i;
_temp_index_3 := i;
temp3=slice[_temp_index_3];
elem = float64(temp3);
_temp_index_4 := idx;
output[_temp_index_4]=elem;
i = i+1;
 }
 }
 dieseldist.EndTiming() ;


  dieseldist.CleanupMain()
  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  dieseldist.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldist.PingMain(tid)
  var DynMap [0]dieseldist.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
var array [10000]float64;
var slice [1000]float64;
var array32 [10000]float32;
var slice32 [1000]float32;
var i int;
var j int;
var c1 int;
var c2 int;
var k int;
var cond1 int;
var cond2 int;
var conditional int;
var point float64;
var temp1 float64;
var temp2 float64;
var temp3 float64;
var temp4 float64;
var temp5 float64;
var temp6 float64;
var temp7 float64;
var temp8 float64;
var temp32 float32;
dieseldist.ReceiveFloat32Array(array32[:], tid, 0);
i = 0;
for __temp_2 := 0; __temp_2 < 10000; __temp_2++ {
 _temp_index_1 := i;
temp32=array32[_temp_index_1];
temp1 = float64(temp32);
_temp_index_2 := i;
array[_temp_index_2]=temp1;
i = i+1;
 }
i = (q-1)*10;
k = 0;
c1 = 10;
c2 = 98;
for __temp_3 := 0; __temp_3 < c1; __temp_3++ {
 j = 1;
cond1 = dieseldist.ConvBool(i<99);
cond2 = dieseldist.ConvBool(i>0);
conditional = dieseldist.ConvBool(cond1==1 && cond2==1);
if conditional != 0 {
 for __temp_4 := 0; __temp_4 < c2; __temp_4++ {
 _temp_index_3 := i*100+j-101;
temp1=array[_temp_index_3];
_temp_index_4 := i*100+j-100;
temp2=array[_temp_index_4];
_temp_index_5 := i*100+j-99;
temp3=array[_temp_index_5];
_temp_index_6 := i*100+j+99;
temp4=array[_temp_index_6];
_temp_index_7 := i*100+j+100;
temp5=array[_temp_index_7];
_temp_index_8 := i*100+j+101;
temp6=array[_temp_index_8];
temp7 = temp2+temp2;
temp8 = temp5+temp5;
point = temp1+temp7;
point = point+temp3;
point = point-temp4;
point = point-temp8;
point = point-temp6;
_temp_index_9 := k;
slice[_temp_index_9]=point;
j = j+1;
k = k+1;
 }
 } else {
 for __temp_5 := 0; __temp_5 < c2; __temp_5++ {
 _temp_index_10 := i*100+j;
temp7=array[_temp_index_10];
_temp_index_11 := k;
slice[_temp_index_11]=temp7;
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
temp32 = float32(temp1);
_temp_index_13 := i;
slice32[_temp_index_13]=temp32;
i = i+1;
 }
dieseldist.SendFloat32Array(slice32[:], tid, 0);

  fmt.Println("Ending thread : ", q);
}

func main() {
	fmt.Println("Starting main thread");
  tid, _ := strconv.Atoi(os.Args[1])

  Num_threads = 11;
  // using this s#!t because go complains when math not used
  _ = math.Inf  

  func_Q(tid)

}
