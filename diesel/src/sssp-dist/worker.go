package main

import "fmt"
import "dieseldistrel"
import "os"
import  "strconv"
import  "io/ioutil"
import  "strings"

var Num_threads int
var Edges [36682000]int
var Inlinks [36682]int
var Outlinks [36682]int
var DistGlobal [36682]int
var Num_nodes int
var Num_edges int
var NodesPerThread int

func max(x, y int) int {
	if x > y {
		return x;
	} else {
		return y;
	}
}

func min(x, y int) int {
	if x < y {
		return x;
	} else {
		return y;
	}
}

func convertToFloat(x int) float64 {
	return float64(x)
}
	
var Q = []int {1,2,3,4,5,6,7,8,9,10};


func func_0() {
  dieseldistrel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistrel.WaitForWorkers(Num_threads)
  var DynMap [46683]float32;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var distance [36682]int;
dieseldistrel.InitDynArray(0, 36682, DynMap[:]);
var slice [10000]int;
dieseldistrel.InitDynArray(36682, 10000, DynMap[:]);
var newDist int;
DynMap[46682] = float32(1.0);
var mystart int;
var myend int;
var i int;
var j int;
var lastthread int;
var mysize int;
distance=DistGlobal;
dieseldistrel.CopyDynArray(0, 0, 36682, DynMap[:]);
i = 0;
for _, q := range(Q) {
 mystart = i*NodesPerThread;
myend = (i+1)*NodesPerThread;
lastthread = dieseldistrel.ConvBool(i==Num_threads-1);
if lastthread != 0 {
 myend = Num_nodes;
 }
dieseldistrel.SendInt(mystart, 0, q);
dieseldistrel.SendInt(myend, 0, q);
i = i+1;
 }
 dieseldistrel.StartTiming() ;
for __temp_0 := 0; __temp_0 < 10; __temp_0++ {
 for _, q := range(Q) {
 dieseldistrel.SendDynIntArrayO1(distance[:], 0, q, DynMap[:], 0);
 }
i = 0;
for _, q := range(Q) {
 mystart = i*NodesPerThread;
myend = (i+1)*NodesPerThread;
lastthread = dieseldistrel.ConvBool(i==Num_threads-1);
if lastthread != 0 {
 myend = Num_nodes;
 }
dieseldistrel.ReceiveDynIntArrayO1(slice[:], 0, q, DynMap[:], 36682);
mysize = myend-mystart;
j = 0;
for __temp_1 := 0; __temp_1 < mysize; __temp_1++ {
 _temp_index_1 := j;
newDist=slice[_temp_index_1];
DynMap[46682] = DynMap[36682 + _temp_index_1];
_temp_index_2 := mystart+j;
distance[_temp_index_2]=newDist;
DynMap[0 + _temp_index_2] = DynMap[46682];
j = j+1;
 }
i = i+1;
 }
 }
 dieseldistrel.EndTiming() ;
DistGlobal=distance;
dieseldistrel.CopyDynArray(0, 0, 36682, DynMap[:]);


  dieseldistrel.CleanupMain()
  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  dieseldistrel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldistrel.PingMain(tid)
  var DynMap [46686]float32;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
var edges [36682000]int;
var inlinks [36682]int;
var distances [36682]int;
dieseldistrel.InitDynArray(0, 36682, DynMap[:]);
var distance int;
DynMap[36682] = float32(1.0);
var newDistance [10000]int;
dieseldistrel.InitDynArray(36683, 10000, DynMap[:]);
var condition int;
DynMap[46683] = float32(1.0);
var inlink int;
var neighbor int;
var nodeInlinks int;
var i int;
var mystart int;
var myend int;
var cur int;
var temp int;
DynMap[46684] = float32(1.0);
var temp1 int;
DynMap[46685] = float32(1.0);
var mysize int;
dieseldistrel.ReceiveInt(&mystart, tid, 0);
dieseldistrel.ReceiveInt(&myend, tid, 0);
edges = Edges;
inlinks = Inlinks;
for __temp_2 := 0; __temp_2 < 10; __temp_2++ {
 dieseldistrel.ReceiveDynIntArrayO1(distances[:], tid, 0, DynMap[:], 0);
mysize = myend-mystart;
i = 0;
for __temp_3 := 0; __temp_3 < mysize; __temp_3++ {
 cur = mystart+i;
_temp_index_1 := cur;
nodeInlinks=inlinks[_temp_index_1];
inlink = 0;
_temp_index_2 := cur;
distance=distances[_temp_index_2];
DynMap[36682] = DynMap[0 + _temp_index_2];
for __temp_4 := 0; __temp_4 < nodeInlinks; __temp_4++ {
 _temp_index_3 := cur*1000+inlink;
neighbor=edges[_temp_index_3];
_temp_index_4 := neighbor;
temp=distances[_temp_index_4];
DynMap[46684] = DynMap[0 + _temp_index_4];
DynMap[46683] = DynMap[36682] + DynMap[46684] - 1.0;
condition = dieseldistrel.ConvBool(distance>temp+1);
DynMap[46684] = DynMap[46684];
temp = temp+1;
DynMap[46685] = 1;
temp1 = 0;
temp_bool_5:= condition; if temp_bool_5 != 0 { distance  = temp } else { distance = temp1 };
if temp_bool_5 != 0 {DynMap[36682]  = DynMap[46683] * DynMap[46684]} else { DynMap[36682] = DynMap[46683] * DynMap[46685]};
inlink = inlink+1;
 }
_temp_index_6 := i;
newDistance[_temp_index_6]=distance;
DynMap[36683 + _temp_index_6] = DynMap[36682];
i = i+1;
 }
dieseldistrel.SendDynIntArrayO1(newDistance[:], tid, 0, DynMap[:], 36683);
 }

  fmt.Println("Ending thread : ", q);
}

func main() {
	tid, _ := strconv.Atoi(os.Args[1])	
	fmt.Println("Starting worker thread: ", tid)
  
  Num_threads = 11;

  data_bytes, _ := ioutil.ReadFile("../../inputs/p2p-Gnutella30.txt")
  Num_nodes = 36682 // strconv.Atoi(os.Args[2])
  Num_edges = Num_nodes * 1000

  fmt.Println("Starting reading the file")
  data_string := string(data_bytes)
  data_str_array := strings.Split(data_string, "\n")

  fmt.Println("Setting up the data structures")

  for i := range Inlinks{
    Inlinks[i] = 0
    Outlinks[i] = 0
    DistGlobal[i] = 0
  }
  DistGlobal[0] = 1

  NodesPerThread = Num_nodes/Num_threads;

  fmt.Println("Populating the data structures")
  for i := 1; i<len(data_str_array)-1 ; i++ {
    elements := strings.Fields(data_str_array[i])
    index1, _ := strconv.Atoi(elements[0])
    index2, _ := strconv.Atoi(elements[1])

    Edges[(index2 * 1000) + Inlinks[index2]] = index1
    Inlinks[index2]++
    Outlinks[index1]++
		// fmt.Println("---------------")
  }

	func_Q(tid)

}
