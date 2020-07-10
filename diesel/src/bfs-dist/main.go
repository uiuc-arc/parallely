package main

import "fmt"
import "diesel"
import "os"
import  "io/ioutil"
import  "strings"
import  "strconv"

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
  diesel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  diesel.WaitForWorkers(Num_threads)
  var DynMap [46683]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var distance [36682]int;
diesel.InitDynArray(0, 36682, DynMap[:]);
var slice [10000]int;
diesel.InitDynArray(36682, 10000, DynMap[:]);
var newDist int;
DynMap[46682] = diesel.ProbInterval{1, 0};
var mystart int;
var myend int;
var i int;
var j int;
var lastthread int;
var mysize int;
distance=DistGlobal;
diesel.CopyDynArray(0, 0, 36682, DynMap[:]);
i = 0;
for _, q := range(Q) {
 mystart = i*NodesPerThread;
myend = (i+1)*NodesPerThread;
lastthread = diesel.ConvBool(i==Num_threads-1);
if lastthread != 0 {
 myend = Num_nodes;
 }
diesel.SendInt(mystart, 0, q);
diesel.SendInt(myend, 0, q);
i = i+1;
 }
 diesel.StartTiming() ;
for __temp_0 := 0; __temp_0 < 10; __temp_0++ {
 for _, q := range(Q) {
 diesel.SendDynIntArrayO1(distance[:], 0, q, DynMap[:], 0);
 }
i = 0;
for _, q := range(Q) {
 mystart = i*NodesPerThread;
myend = (i+1)*NodesPerThread;
lastthread = diesel.ConvBool(i==Num_threads-1);
if lastthread != 0 {
 myend = Num_nodes;
 }
diesel.ReceiveDynIntArrayO1(slice[:], 0, q, DynMap[:], 36682);
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
 diesel.EndTiming() ;
DistGlobal=distance;
diesel.CopyDynArray(0, 0, 36682, DynMap[:]);


  diesel.CleanupMain()
  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  diesel.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  diesel.PingMain()
  var DynMap [46687]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
var edges [36682000]int;
var inlinks [36682]int;
var distances [36682]int;
diesel.InitDynArray(0, 36682, DynMap[:]);
var distance int;
DynMap[36682] = diesel.ProbInterval{1, 0};
var newDistance [10000]int;
diesel.InitDynArray(36683, 10000, DynMap[:]);
var condition int;
DynMap[46683] = diesel.ProbInterval{1, 0};
var inlink int;
var neighbor int;
var nodeInlinks int;
var i int;
var mystart int;
var myend int;
var cur int;
var temp int;
DynMap[46684] = diesel.ProbInterval{1, 0};
var temp1 int;
DynMap[46685] = diesel.ProbInterval{1, 0};
var temp2 int;
DynMap[46686] = diesel.ProbInterval{1, 0};
var mysize int;
diesel.ReceiveInt(&mystart, tid, 0);
diesel.ReceiveInt(&myend, tid, 0);
edges = Edges;
inlinks = Inlinks;
for __temp_2 := 0; __temp_2 < 10; __temp_2++ {
 diesel.ReceiveDynIntArrayO1(distances[:], tid, 0, DynMap[:], 0);
mysize = myend-mystart;
i = 0;
for __temp_3 := 0; __temp_3 < mysize; __temp_3++ {
 cur = mystart+i;
_temp_index_1 := cur;
nodeInlinks=inlinks[_temp_index_1];
_temp_index_2 := cur;
distance=distances[_temp_index_2];
DynMap[36682] = DynMap[0 + _temp_index_2];
inlink = 0;
for __temp_4 := 0; __temp_4 < nodeInlinks; __temp_4++ {
 _temp_index_3 := cur*1000+inlink;
neighbor=edges[_temp_index_3];
_temp_index_4 := neighbor;
temp=distances[_temp_index_4];
DynMap[46684] = DynMap[0 + _temp_index_4];
DynMap[46683].Reliability = DynMap[46684].Reliability;
condition = diesel.ConvBool(temp==1);
DynMap[46685] = diesel.ProbInterval{1, 0};
temp1 = 1;
DynMap[46686] = diesel.ProbInterval{1, 0};
temp2 = 2;
temp_bool_5:= condition; if temp_bool_5 != 0 { distance  = temp1 } else { distance = temp2 };
if temp_bool_5 != 0 {DynMap[36682].Reliability  = DynMap[46683].Reliability * DynMap[46685].Reliability} else { DynMap[36682].Reliability = DynMap[46683].Reliability * DynMap[46686].Reliability};
inlink = inlink+1;
 }
DynMap[46684].Reliability = DynMap[36682].Reliability;
DynMap[46684].Delta = DynMap[36682].Delta;
temp = distance;
_temp_index_6 := i;
newDistance[_temp_index_6]=temp;
DynMap[36683 + _temp_index_6] = DynMap[46684];
i = i+1;
 }
diesel.SendDynIntArrayO1(newDistance[:], tid, 0, DynMap[:], 36683);
 }

  diesel.CleanupMain()
  fmt.Println("Ending thread : ", q);
}

func main() {
	fmt.Println("Starting main thread");

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

  fmt.Println("Number of worker threads: ", Num_threads);
  fmt.Println("Number of nodes: ", len(DistGlobal));
  fmt.Println("Size of Inlinks: ", len(Inlinks));

  // fmt.Println("Starting the iterations")
  // startTime := time.Now()

	func_0();


	// fmt.Println("Main thread waiting for others to finish");  
	// diesel.Wg.Wait()
  // elapsed := time.Since(startTime)

	fmt.Println("Done!");
  // fmt.Println("Elapsed time : ", elapsed.Nanoseconds());

  f, _ := os.Create("output.txt")
  defer f.Close()

  for i := range DistGlobal{
    f.WriteString(fmt.Sprintln(DistGlobal[i]))
  }
}
