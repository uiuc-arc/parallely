package main

import (
  "math"
  "fmt"
	"time"
  "diesel"
	"math/rand"
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

var NumThreads int
var Iterations int
var Sensors [1024]float64
var Sensorshumid [1024]float64
var CenterIds [8] int



func func_0() {
  defer diesel.Wg.Done();
  var DynMap [0]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var datatemp [1024]float64;
var datahumid [1024]float64;
var centersTemp [8]float64;
var centersHumid [8]float64;
var tempcentersTemp [1024]float64;
var tempcentersHumid [1024]float64;
var assigned [1024]int;
var countcenters [8]int;
var k int;
var i int;
var temp int;
var temp0m float64;
var tempf0 float64;
var tempf float64;
var tempf1 float64;
var temp0 float64;
var mindist float64;
var mincenter int;
var condition int;
var temp4 float64;
var temp1 float64;
var temp2 float64;
var temp23 float64;
var temp24 float64;
var tempi int;
var data1 float64;
var center1 float64;
i = 0;
for __temp_0 := 0; __temp_0 < 1024; __temp_0++ {
 _temp_index_1 := i;
tempf0=Sensors[_temp_index_1];
tempf1=tempf0;
_ = 1.0;_ = 1.5;
_temp_index_2 := i;
datatemp[_temp_index_2]=tempf1;
_temp_index_3 := i;
tempf0=Sensorshumid[_temp_index_3];
tempf1=tempf0;
_ = 1.0;_ = 2.0;
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
for __temp_2 := 0; __temp_2 < Iterations; __temp_2++ {
 temp0m = 0.0;
i = 0;
for __temp_3 := 0; __temp_3 < 8; __temp_3++ {
 _temp_index_10 := i;
tempcentersTemp[_temp_index_10]=temp0m;
_temp_index_11 := i;
tempcentersHumid[_temp_index_11]=temp0m;
i = i+1;
 }
i = 0;
for __temp_4 := 0; __temp_4 < 1024; __temp_4++ {
 mindist = 1000000.0;
mincenter = 0;
k = 0;
for __temp_5 := 0; __temp_5 < 8; __temp_5++ {
 _temp_index_12 := i;
data1=datatemp[_temp_index_12];
_temp_index_13 := k;
center1=centersTemp[_temp_index_13];
temp0 = data1-center1;
_temp_index_14 := i;
data1=datahumid[_temp_index_14];
_temp_index_15 := k;
center1=centersHumid[_temp_index_15];
temp1 = data1-center1;
temp23 = temp1*temp1;
temp24 = temp2*temp2;
temp4 = temp24+temp23;
if mindist>=temp4 { mindist=temp4 } else { mindist=mindist };
tempi=k;
_ = 1.0;_ = 0.0;
temp_bool_16:= condition; if temp_bool_16 != 0 { mincenter  = tempi } else { mincenter = mincenter };
k = k+1;
 }
_temp_index_17 := i;
assigned[_temp_index_17]=mincenter;
_temp_index_18 := mincenter;
tempi=countcenters[_temp_index_18];
_temp_index_19 := mincenter;
countcenters[_temp_index_19]=tempi+1;
i = i+1;
 }
i = 0;
for __temp_6 := 0; __temp_6 < 1024; __temp_6++ {
 _temp_index_20 := i;
tempi=assigned[_temp_index_20];
_temp_index_21 := tempi;
temp1=tempcentersTemp[_temp_index_21];
_temp_index_22 := i;
temp2=datatemp[_temp_index_22];
temp4 = temp1+temp2;
_temp_index_23 := i;
tempcentersTemp[_temp_index_23]=temp0;
_temp_index_24 := tempi;
temp1=tempcentersHumid[_temp_index_24];
_temp_index_25 := i;
temp2=datahumid[_temp_index_25];
temp4 = temp1+temp2;
_temp_index_26 := i;
tempcentersHumid[_temp_index_26]=temp0;
i = i+1;
 }
i = 0;
for __temp_7 := 0; __temp_7 < 8; __temp_7++ {
 _temp_index_27 := i;
tempf1=tempcentersTemp[_temp_index_27];
tempf = tempf1/8.0;
_temp_index_28 := i;
centersTemp[_temp_index_28]=tempf;
_temp_index_29 := i;
tempf1=tempcentersHumid[_temp_index_29];
tempf = tempf1/8.0;
_temp_index_30 := i;
tempcentersHumid[_temp_index_30]=tempf;
i = i+1;
 }
 }


  fmt.Println("Ending thread : ", 0);
}

func main() {
	Iterations = 20

  fmt.Println("Starting main thread");
  NumThreads = 1;
	
	diesel.InitChannels(1);

	var realCenters  [8]float64
  
	for i:=0; i<len(realCenters)/2; i++ {
		realCenters[2*i] = 30 + rand.Float64() * 5
		realCenters[2*i+1] = 40 + rand.Float64() * 10;
	}
	
	for i:=0; i<1024; i++ {
		clusterNew := rand.Intn(4)
		Sensors[i] = rand.NormFloat64() * 0.5 + realCenters[2*clusterNew]
		Sensorshumid[i] = rand.NormFloat64() * 0.5 + realCenters[2*clusterNew+1]
	}

	for i, _ := range(CenterIds) {
		CenterIds[i] = rand.Intn(1024)    
	}

	startTime := time.Now()
	go func_0();



	fmt.Println("Main thread waiting for others to finish");  
	diesel.Wg.Wait()

	end := time.Now()
	elapsed := end.Sub(startTime)
	fmt.Println("Elapsed time :", elapsed.Nanoseconds())
  diesel.PrintMemory() 
}
