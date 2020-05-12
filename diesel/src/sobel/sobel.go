package main

import "fmt"
import "math/rand"
import "time"
import "math"
import "diesel"

var Input_array [10000]float64
var Num_threads int

var Q = []int {1,2,3,4,5,6,7,8,9,10};


func func_0() {
  defer diesel.Wg.Done();
  var DynMap [0]diesel.ProbInterval;
  _ = DynMap;
  var output [10000]float64;
var i int;
var idx int;
var temp1 float64;
var temp2 float32;
var temp3 float32;
var input32 [10000]float32;
var slice [1000]float32;
var elem float64;
i = 0;
for __temp_0 := 0; __temp_0 < 10000; __temp_0++ {
 _temp_index_1 := i;
temp1=Input_array[_temp_index_1];
temp2 = float32(temp1);
_temp_index_2 := i;
input32[_temp_index_2]=temp2;
i = i+1;
 }
for _, q := range(Q) {
 diesel.SendFloat32Array(input32[:], 0, q);
 }
for _, q := range(Q) {
 diesel.ReceiveFloat32Array(slice[:], 0, q);
i = 0;
for __temp_1 := 0; __temp_1 < 1000; __temp_1++ {
 idx = (q-1)*100+i;
_temp_index_3 := i;
temp3=slice[_temp_index_3];
elem = float64(temp3);
_temp_index_4 := idx;
output[_temp_index_4]=elem;
i = i+1;
 }
 }


  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  defer diesel.Wg.Done();
  var DynMap [0]diesel.ProbInterval;
  _ = DynMap;
  q := tid;
var array [10000]float32;
var slice [1000]float32;
var i int;
var j int;
var k int;
var conditional int;
var point float32;
var temp1 float32;
var temp2 float32;
var temp3 float32;
var temp4 float32;
var temp5 float32;
var temp6 float32;
var temp7 float32;
var temp8 float32;
diesel.ReceiveFloat32Array(array[:], tid, 0);
i = (q-1)*10;
k = 0;
for __temp_2 := 0; __temp_2 < 10; __temp_2++ {
 j = 1;
conditional = diesel.ConvBool((i<99)&&(i>0));
if conditional != 0 {
 for __temp_3 := 0; __temp_3 < 98; __temp_3++ {
 _temp_index_1 := i*100+j-101;
temp1=array[_temp_index_1];
_temp_index_2 := i*100+j-100;
temp2=array[_temp_index_2];
_temp_index_3 := i*100+j-99;
temp3=array[_temp_index_3];
_temp_index_4 := i*100+j+99;
temp4=array[_temp_index_4];
_temp_index_5 := i*100+j+100;
temp5=array[_temp_index_5];
_temp_index_6 := i*100+j+101;
temp6=array[_temp_index_6];
temp7 = temp2+temp2;
temp8 = temp5+temp5;
point = temp1+temp7;
point = point+temp3;
point = point-temp4;
point = point-temp8;
point = point-temp6;
_temp_index_7 := k;
slice[_temp_index_7]=point;
j = j+1;
k = k+1;
 }
 } else {
 for __temp_4 := 0; __temp_4 < 98; __temp_4++ {
 _temp_index_8 := i*100+j;
temp7=array[_temp_index_8];
_temp_index_9 := k;
slice[_temp_index_9]=temp7;
j = j+1;
k = k+1;
 }
 }
i = i+1;
 }
diesel.SendFloat32Array(slice[:], tid, 0);

  fmt.Println("Ending thread : ", q);
}

func main() {
	fmt.Println("Starting main thread");

  Num_threads = 11;
	
	diesel.InitChannels(11);

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
	diesel.Wg.Wait()
  elapsed := time.Since(startTime)
  
  fmt.Println("Elapsed time : ", elapsed.Nanoseconds())
}
