package main

import (
  "math"
  "fmt"
  "dieseldist"
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

var Num_threads int
var Iterations int
var Sensors [1024]float64
var Sensorshumid [1024]float64
var CenterIds [8] int

var Q = []int {1,2,3,4,5,6,7,8};


func func_0() {
  dieseldist.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldist.WaitForWorkers(Num_threads)
  var DynMap [2092]dieseldist.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var datatemp [1024]float64;
dieseldist.InitDynArray(0, 1024, DynMap[:]);
var datahumid [1024]float64;
dieseldist.InitDynArray(1024, 1024, DynMap[:]);
var centersTemp [8]float64;
dieseldist.InitDynArray(2048, 8, DynMap[:]);
var centersHumid [8]float64;
dieseldist.InitDynArray(2056, 8, DynMap[:]);
var centerSlice [8]float64;
dieseldist.InitDynArray(2064, 8, DynMap[:]);
var tempcentersTemp [8]float64;
dieseldist.InitDynArray(2072, 8, DynMap[:]);
var tempcentersHumid [8]float64;
dieseldist.InitDynArray(2080, 8, DynMap[:]);
var i int;
var temp int;
var temp0 float64;
var tempf0 float64;
var tempf float64;
DynMap[2088] = dieseldist.ProbInterval{1, 0};
var tempf1 float64;
DynMap[2089] = dieseldist.ProbInterval{1, 0};
var tempf2 float64;
DynMap[2090] = dieseldist.ProbInterval{1, 0};
var temp1 float64;
DynMap[2091] = dieseldist.ProbInterval{1, 0};
i = 0;
for __temp_0 := 0; __temp_0 < 1024; __temp_0++ {
 _temp_index_1 := i;
tempf0=Sensors[_temp_index_1];
tempf1=tempf0;
DynMap[2089] = dieseldist.ProbInterval{1.5, 1.0};
_temp_index_2 := i;
datatemp[_temp_index_2]=tempf1;
DynMap[0 + _temp_index_2] = DynMap[2089];
_temp_index_3 := i;
tempf0=Sensorshumid[_temp_index_3];
tempf1=tempf0;
DynMap[2089] = dieseldist.ProbInterval{2.0, 1.0};
_temp_index_4 := i;
datahumid[_temp_index_4]=tempf1;
DynMap[1024 + _temp_index_4] = DynMap[2089];
 }
i = 0;
for __temp_1 := 0; __temp_1 < 8; __temp_1++ {
 _temp_index_5 := i;
temp=CenterIds[_temp_index_5];
_temp_index_6 := temp;
tempf=datatemp[_temp_index_6];
DynMap[2088] = DynMap[0 + _temp_index_6];
_temp_index_7 := i;
centersTemp[_temp_index_7]=tempf;
DynMap[2048 + _temp_index_7] = DynMap[2088];
_temp_index_8 := temp;
tempf=datahumid[_temp_index_8];
DynMap[2088] = DynMap[1024 + _temp_index_8];
_temp_index_9 := 1;
centersHumid[_temp_index_9]=tempf;
DynMap[2056 + _temp_index_9] = DynMap[2088];
i = i+1;
 }
 dieseldist.StartTiming() ;
for _, q := range(Q) {
 dieseldist.SendDynFloat64ArrayO1(datatemp[:], 0, q, DynMap[:], 0);
dieseldist.SendDynFloat64ArrayO1(datahumid[:], 0, q, DynMap[:], 1024);
 }
for __temp_2 := 0; __temp_2 < Iterations; __temp_2++ {
 for _, q := range(Q) {
 dieseldist.SendDynFloat64ArrayO1(centersTemp[:], 0, q, DynMap[:], 2048);
dieseldist.SendDynFloat64ArrayO1(centersHumid[:], 0, q, DynMap[:], 2056);
 }
temp0 = 0.0;
i = 0;
for __temp_3 := 0; __temp_3 < 8; __temp_3++ {
 _temp_index_10 := i;
tempcentersTemp[_temp_index_10]=temp0;
DynMap[2072 + _temp_index_10] = dieseldist.ProbInterval{1, 0};
_temp_index_11 := i;
tempcentersHumid[_temp_index_11]=temp0;
DynMap[2080 + _temp_index_11] = dieseldist.ProbInterval{1, 0};
i = i+1;
 }
for _, q := range(Q) {
 dieseldist.ReceiveDynFloat64ArrayO1(centerSlice[:], 0, q, DynMap[:], 2064);
i = 0;
for __temp_4 := 0; __temp_4 < 8; __temp_4++ {
 _temp_index_12 := i;
tempf=tempcentersTemp[_temp_index_12];
DynMap[2088] = DynMap[2072 + _temp_index_12];
_temp_index_13 := i;
tempf1=centerSlice[_temp_index_13];
DynMap[2089] = DynMap[2064 + _temp_index_13];
DynMap[2090].Reliability = DynMap[2091].Reliability + DynMap[2088].Reliability - 1.0;
DynMap[2090].Delta = DynMap[2088].Delta + DynMap[2091].Delta;
tempf2 = tempf+temp1;
_temp_index_14 := i;
tempcentersTemp[_temp_index_14]=tempf2;
DynMap[2072 + _temp_index_14] = DynMap[2090];
i = i+1;
 }
dieseldist.ReceiveDynFloat64ArrayO1(centerSlice[:], 0, q, DynMap[:], 2064);
i = 0;
for __temp_5 := 0; __temp_5 < 8; __temp_5++ {
 _temp_index_15 := i;
tempf=tempcentersHumid[_temp_index_15];
DynMap[2088] = DynMap[2080 + _temp_index_15];
_temp_index_16 := i;
tempf1=centerSlice[_temp_index_16];
DynMap[2089] = DynMap[2064 + _temp_index_16];
DynMap[2090].Reliability = DynMap[2091].Reliability + DynMap[2088].Reliability - 1.0;
DynMap[2090].Delta = DynMap[2088].Delta + DynMap[2091].Delta;
tempf2 = tempf+temp1;
_temp_index_17 := i;
tempcentersHumid[_temp_index_17]=tempf2;
DynMap[2080 + _temp_index_17] = DynMap[2090];
i = i+1;
 }
 }
i = 0;
for __temp_6 := 0; __temp_6 < 8; __temp_6++ {
 _temp_index_18 := i;
tempf1=tempcentersTemp[_temp_index_18];
DynMap[2089] = DynMap[2072 + _temp_index_18];
DynMap[2088].Reliability = DynMap[2089].Reliability;
DynMap[2088].Delta =  DynMap[2089].Delta / math.Abs(8.0);
tempf = tempf1/8.0;
_temp_index_19 := i;
centersTemp[_temp_index_19]=tempf;
DynMap[2048 + _temp_index_19] = DynMap[2088];
_temp_index_20 := i;
tempf1=tempcentersHumid[_temp_index_20];
DynMap[2089] = DynMap[2080 + _temp_index_20];
DynMap[2088].Reliability = DynMap[2089].Reliability;
DynMap[2088].Delta =  DynMap[2089].Delta / math.Abs(8.0);
tempf = tempf1/8.0;
_temp_index_21 := i;
tempcentersHumid[_temp_index_21]=tempf;
DynMap[2080 + _temp_index_21] = DynMap[2088];
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
  var DynMap [3123]dieseldist.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
var datatemp [1024]float64;
dieseldist.InitDynArray(0, 1024, DynMap[:]);
var datahumid [1024]float64;
dieseldist.InitDynArray(1024, 1024, DynMap[:]);
var centersTemp [8]float64;
dieseldist.InitDynArray(2048, 8, DynMap[:]);
var centersHumid [8]float64;
dieseldist.InitDynArray(2056, 8, DynMap[:]);
var tempcentersTemp [8]float64;
dieseldist.InitDynArray(2064, 8, DynMap[:]);
var tempcentersHumid [8]float64;
dieseldist.InitDynArray(2072, 8, DynMap[:]);
var countcenters [8]int;
dieseldist.InitDynArray(2080, 8, DynMap[:]);
var assigned [1024]int;
dieseldist.InitDynArray(2088, 1024, DynMap[:]);
var mystart int;
var myend int;
var perthread int;
var mypoints int;
var i int;
var k int;
var temp0 float64;
DynMap[3112] = dieseldist.ProbInterval{1, 0};
var mindist float64;
DynMap[3113] = dieseldist.ProbInterval{1, 0};
var mincenter int;
DynMap[3114] = dieseldist.ProbInterval{1, 0};
var temp4 float64;
DynMap[3115] = dieseldist.ProbInterval{1, 0};
var temp1 float64;
DynMap[3116] = dieseldist.ProbInterval{1, 0};
var temp2 float64;
DynMap[3117] = dieseldist.ProbInterval{1, 0};
var temp23 float64;
DynMap[3118] = dieseldist.ProbInterval{1, 0};
var temp24 float64;
DynMap[3119] = dieseldist.ProbInterval{1, 0};
var tempi int;
DynMap[3120] = dieseldist.ProbInterval{1, 0};
var data1 float64;
DynMap[3121] = dieseldist.ProbInterval{1, 0};
var center1 float64;
DynMap[3122] = dieseldist.ProbInterval{1, 0};
perthread = 1024/8;
mystart = (q-1)*perthread;
myend = mystart+perthread;
mypoints = myend-mystart;
dieseldist.ReceiveDynFloat64ArrayO1(datatemp[:], tid, 0, DynMap[:], 0);
dieseldist.ReceiveDynFloat64ArrayO1(datahumid[:], tid, 0, DynMap[:], 1024);
for __temp_7 := 0; __temp_7 < Iterations; __temp_7++ {
 dieseldist.ReceiveDynFloat64ArrayO1(centersTemp[:], tid, 0, DynMap[:], 2048);
dieseldist.ReceiveDynFloat64ArrayO1(centersHumid[:], tid, 0, DynMap[:], 2056);
DynMap[3112] = dieseldist.ProbInterval{1, 0};
temp0 = 0.0;
i = 0;
for __temp_8 := 0; __temp_8 < 8; __temp_8++ {
 _temp_index_1 := i;
tempcentersTemp[_temp_index_1]=temp0;
DynMap[2064 + _temp_index_1] = DynMap[3112];
_temp_index_2 := i;
tempcentersHumid[_temp_index_2]=temp0;
DynMap[2072 + _temp_index_2] = DynMap[3112];
i = i+1;
 }
i = mystart;
for __temp_9 := 0; __temp_9 < mypoints; __temp_9++ {
 DynMap[3113] = dieseldist.ProbInterval{1, 0};
mindist = 1000000.0;
DynMap[3114] = dieseldist.ProbInterval{1, 0};
mincenter = 0;
k = 0;
for __temp_10 := 0; __temp_10 < 8; __temp_10++ {
 _temp_index_3 := i;
data1=datatemp[_temp_index_3];
DynMap[3121] = DynMap[0 + _temp_index_3];
_temp_index_4 := k;
center1=centersTemp[_temp_index_4];
DynMap[3122] = DynMap[2048 + _temp_index_4];
DynMap[3112].Reliability = DynMap[3121].Reliability + DynMap[3122].Reliability - 1.0;
DynMap[3112].Delta = DynMap[3121].Delta + DynMap[3122].Delta;
temp0 = data1-center1;
_temp_index_5 := i;
data1=datahumid[_temp_index_5];
DynMap[3121] = DynMap[1024 + _temp_index_5];
_temp_index_6 := k;
center1=centersHumid[_temp_index_6];
DynMap[3122] = DynMap[2056 + _temp_index_6];
DynMap[3116].Reliability = DynMap[3121].Reliability + DynMap[3122].Reliability - 1.0;
DynMap[3116].Delta = DynMap[3121].Delta + DynMap[3122].Delta;
temp1 = data1-center1;
DynMap[3118].Reliability = DynMap[3116].Reliability;
DynMap[3118].Delta = math.Abs(float64(temp1)) * DynMap[3116].Delta + math.Abs(float64(temp1)) * DynMap[3116].Delta + DynMap[3116].Delta*DynMap[3116].Delta;
temp23 = temp1*temp1;
DynMap[3119].Reliability = DynMap[3117].Reliability;
DynMap[3119].Delta = math.Abs(float64(temp2)) * DynMap[3117].Delta + math.Abs(float64(temp2)) * DynMap[3117].Delta + DynMap[3117].Delta*DynMap[3117].Delta;
temp24 = temp2*temp2;
DynMap[3115].Reliability = DynMap[3118].Reliability + DynMap[3119].Reliability - 1.0;
DynMap[3115].Delta = DynMap[3119].Delta + DynMap[3118].Delta;
temp4 = temp24+temp23;
mindist = dieseldist.DynCondFloat64GeqFloat64(mindist, temp4, DynMap[:], 3113, 3115, temp4, mindist, 3115, 3113, 3113);
tempi=k;
DynMap[3120] = dieseldist.ProbInterval{1.0, 0.0};
mincenter = dieseldist.DynCondFloat64GeqInt(mindist, temp4, DynMap[:], 3113, 3115, tempi, mincenter, 3120, 3114, 3114);
k = k+1;
 }
_temp_index_7 := i;
assigned[_temp_index_7]=mincenter;
DynMap[2088 + _temp_index_7] = DynMap[3114];
_temp_index_8 := mincenter;
tempi=countcenters[_temp_index_8];
DynMap[3120] = DynMap[2080 + _temp_index_8];
_temp_index_9 := mincenter;
countcenters[_temp_index_9]=tempi+1;
DynMap[2080 + _temp_index_9] = DynMap[3120];
i = i+1;
 }
i = mystart;
for __temp_11 := 0; __temp_11 < mypoints; __temp_11++ {
 _temp_index_10 := i;
tempi=assigned[_temp_index_10];
DynMap[3120] = DynMap[2088 + _temp_index_10];
_temp_index_11 := tempi;
temp1=tempcentersTemp[_temp_index_11];
DynMap[3116] = DynMap[2064 + _temp_index_11];
_temp_index_12 := i;
temp2=datatemp[_temp_index_12];
DynMap[3117] = DynMap[0 + _temp_index_12];
DynMap[3115].Reliability = DynMap[3117].Reliability + DynMap[3116].Reliability - 1.0;
DynMap[3115].Delta = DynMap[3116].Delta + DynMap[3117].Delta;
temp4 = temp1+temp2;
_temp_index_13 := tempi;
tempcentersTemp[_temp_index_13]=temp0;
DynMap[2064 + _temp_index_13] = DynMap[3112];
_temp_index_14 := tempi;
temp1=tempcentersHumid[_temp_index_14];
DynMap[3116] = DynMap[2072 + _temp_index_14];
_temp_index_15 := i;
temp2=datahumid[_temp_index_15];
DynMap[3117] = DynMap[1024 + _temp_index_15];
DynMap[3115].Reliability = DynMap[3117].Reliability + DynMap[3116].Reliability - 1.0;
DynMap[3115].Delta = DynMap[3116].Delta + DynMap[3117].Delta;
temp4 = temp1+temp2;
_temp_index_16 := tempi;
tempcentersHumid[_temp_index_16]=temp0;
DynMap[2072 + _temp_index_16] = DynMap[3112];
i = i+1;
 }
dieseldist.SendDynFloat64ArrayO1(tempcentersTemp[:], tid, 0, DynMap[:], 2064);
dieseldist.SendDynFloat64ArrayO1(tempcentersHumid[:], tid, 0, DynMap[:], 2072);
 }

  fmt.Println("Ending thread : ", q);
}

func main() {
	Iterations = 20

  fmt.Println("Starting main thread");
  Num_threads = 9;
	
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

	func_0();


}
