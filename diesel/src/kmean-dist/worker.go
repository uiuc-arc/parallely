package main

import (
  "math"
  "fmt"
  "strconv"
  "os"    
  "dieseldist"
)

func Min(x, y int) int {
  if x < y {
    return x
  }
  return y
}

func Max(x, y int) int {
  if x > y {
    return x
  }
  return y
}

func Idx(i, j, width int) int {
  return i*width+j
}

func floorInt(input float64) int {
	return int(math.Floor(input));
}

func ceilInt(input float64) int {
	return int(math.Ceil(input));
}

func convertToFloat(x int) float64 {
	return float64(x)
}

var Num_threads int
var Iterations int
var Sensors [1024]float64
var Sensorshumid [1024]float64
var CenterIds [8] int

var Q = []int {1,2,3,4,5,6,7,8};


func func_0() {
  dieseldist.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldist.WaitForWorkers(Num_threads)
  var DynMap [0]dieseldist.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var datatemp [1024]float64;
var datahumid [1024]float64;
var centersTemp [8]float64;
var centersHumid [8]float64;
var centerSlice [8]float64;
var tempcentersTemp [8]float64;
var tempcentersHumid [8]float64;
var i int;
var temp int;
var temp0 float64;
var tempf0 float64;
var tempf float64;
var tempf1 float64;
var tempf2 float64;
var temp1 float64;
i = 0;
for __temp_0 := 0; __temp_0 < 1024; __temp_0++ {
 _temp_index_1 := i;
tempf0=Sensors[_temp_index_1];
tempf1=tempf0;
_temp_index_2 := i;
datatemp[_temp_index_2]=tempf1;
_temp_index_3 := i;
tempf0=Sensorshumid[_temp_index_3];
tempf1=tempf0;
_temp_index_4 := i;
datahumid[_temp_index_4]=tempf1;
 }
i = 0;
for __temp_1 := 0; __temp_1 < 8; __temp_1++ {
 _temp_index_5 := i;
temp=CenterIds[_temp_index_5];
_temp_index_6 := temp;
tempf=datatemp[_temp_index_6];
_temp_index_7 := i;
centersTemp[_temp_index_7]=tempf;
_temp_index_8 := temp;
tempf=datahumid[_temp_index_8];
_temp_index_9 := 1;
centersHumid[_temp_index_9]=tempf;
i = i+1;
 }
 dieseldist.StartTiming() ;
for _, q := range(Q) {
 dieseldist.SendFloat64Array(datatemp[:], 0, q);
dieseldist.SendFloat64Array(datahumid[:], 0, q);
 }
for __temp_2 := 0; __temp_2 < Iterations; __temp_2++ {
 for _, q := range(Q) {
 dieseldist.SendFloat64Array(centersTemp[:], 0, q);
dieseldist.SendFloat64Array(centersHumid[:], 0, q);
 }
temp0 = 0.0;
i = 0;
for __temp_3 := 0; __temp_3 < 8; __temp_3++ {
 _temp_index_10 := i;
tempcentersTemp[_temp_index_10]=temp0;
_temp_index_11 := i;
tempcentersHumid[_temp_index_11]=temp0;
i = i+1;
 }
for _, q := range(Q) {
 dieseldist.ReceiveFloat64Array(centerSlice[:], 0, q);
i = 0;
for __temp_4 := 0; __temp_4 < 8; __temp_4++ {
 _temp_index_12 := i;
tempf=tempcentersTemp[_temp_index_12];
_temp_index_13 := i;
tempf1=centerSlice[_temp_index_13];
tempf2 = tempf+temp1;
_temp_index_14 := i;
tempcentersTemp[_temp_index_14]=tempf2;
i = i+1;
 }
dieseldist.ReceiveFloat64Array(centerSlice[:], 0, q);
i = 0;
for __temp_5 := 0; __temp_5 < 8; __temp_5++ {
 _temp_index_15 := i;
tempf=tempcentersHumid[_temp_index_15];
_temp_index_16 := i;
tempf1=centerSlice[_temp_index_16];
tempf2 = tempf+temp1;
_temp_index_17 := i;
tempcentersHumid[_temp_index_17]=tempf2;
i = i+1;
 }
 }
i = 0;
for __temp_6 := 0; __temp_6 < 8; __temp_6++ {
 _temp_index_18 := i;
tempf1=tempcentersTemp[_temp_index_18];
tempf = tempf1/8.0;
_temp_index_19 := i;
centersTemp[_temp_index_19]=tempf;
_temp_index_20 := i;
tempf1=tempcentersHumid[_temp_index_20];
tempf = tempf1/8.0;
_temp_index_21 := i;
tempcentersHumid[_temp_index_21]=tempf;
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
var datatemp [1024]float64;
var datahumid [1024]float64;
var centersTemp [8]float64;
var centersHumid [8]float64;
var tempcentersTemp [8]float64;
var tempcentersHumid [8]float64;
var countcenters [8]int;
var assigned [1024]int;
var mystart int;
var myend int;
var perthread int;
var mypoints int;
var i int;
var k int;
var temp0 float64;
var mindist float64;
var mincenter int;
var temp4 float64;
var temp1 float64;
var temp2 float64;
var temp23 float64;
var temp24 float64;
var tempi int;
var data1 float64;
var center1 float64;
perthread = 1024/8;
mystart = (q-1)*perthread;
myend = mystart+perthread;
mypoints = myend-mystart;
dieseldist.ReceiveFloat64Array(datatemp[:], tid, 0);
dieseldist.ReceiveFloat64Array(datahumid[:], tid, 0);
for __temp_7 := 0; __temp_7 < Iterations; __temp_7++ {
 dieseldist.ReceiveFloat64Array(centersTemp[:], tid, 0);
dieseldist.ReceiveFloat64Array(centersHumid[:], tid, 0);
temp0 = 0.0;
i = 0;
for __temp_8 := 0; __temp_8 < 8; __temp_8++ {
 _temp_index_1 := i;
tempcentersTemp[_temp_index_1]=temp0;
_temp_index_2 := i;
tempcentersHumid[_temp_index_2]=temp0;
i = i+1;
 }
i = mystart;
for __temp_9 := 0; __temp_9 < mypoints; __temp_9++ {
 mindist = 1000000.0;
mincenter = 0;
k = 0;
for __temp_10 := 0; __temp_10 < 8; __temp_10++ {
 _temp_index_3 := i;
data1=datatemp[_temp_index_3];
_temp_index_4 := k;
center1=centersTemp[_temp_index_4];
temp0 = data1-center1;
_temp_index_5 := i;
data1=datahumid[_temp_index_5];
_temp_index_6 := k;
center1=centersHumid[_temp_index_6];
temp1 = data1-center1;
temp23 = temp1*temp1;
temp24 = temp2*temp2;
temp4 = temp24+temp23;
if mindist>=temp4 { mindist=temp4 } else { mindist=mindist };
tempi=k;
if mindist>=temp4 { mincenter=tempi } else { mincenter=mincenter };
k = k+1;
 }
_temp_index_7 := i;
assigned[_temp_index_7]=mincenter;
_temp_index_8 := mincenter;
tempi=countcenters[_temp_index_8];
_temp_index_9 := mincenter;
countcenters[_temp_index_9]=tempi+1;
i = i+1;
 }
i = mystart;
for __temp_11 := 0; __temp_11 < mypoints; __temp_11++ {
 _temp_index_10 := i;
tempi=assigned[_temp_index_10];
_temp_index_11 := tempi;
temp1=tempcentersTemp[_temp_index_11];
_temp_index_12 := i;
temp2=datatemp[_temp_index_12];
temp4 = temp1+temp2;
_temp_index_13 := tempi;
tempcentersTemp[_temp_index_13]=temp0;
_temp_index_14 := tempi;
temp1=tempcentersHumid[_temp_index_14];
_temp_index_15 := i;
temp2=datahumid[_temp_index_15];
temp4 = temp1+temp2;
_temp_index_16 := tempi;
tempcentersHumid[_temp_index_16]=temp0;
i = i+1;
 }
dieseldist.SendFloat64Array(tempcentersTemp[:], tid, 0);
dieseldist.SendFloat64Array(tempcentersHumid[:], tid, 0);
 }

  fmt.Println("Ending thread : ", q);
}

func main() {
	Iterations = 20

	tid, _ := strconv.Atoi(os.Args[1])
  
  fmt.Println("Starting thread: " + os.Args[1]);
  Num_threads = 9;
	
  func_Q(tid)

}
