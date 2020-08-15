package main

import "fmt"
import "dieseldist"
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
  dieseldist.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldist.WaitForWorkers(Num_threads)
  var DynMap [0]dieseldist.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var distance [36682]int;
var slice [10000]int;
var newDist int;
var mystart int;
var myend int;
var i int;
var j int;
var lastthread int;
var mysize int;
distance = DistGlobal;
i = 0;
for _, q := range(Q) {
 mystart = i*NodesPerThread;
myend = (i+1)*NodesPerThread;
lastthread = dieseldist.ConvBool(i==Num_threads-1);
if lastthread != 0 {
 myend = Num_nodes;
 }
dieseldist.SendInt(mystart, 0, q);
dieseldist.SendInt(myend, 0, q);
i = i+1;
 }
 dieseldist.StartTiming() ;
for __temp_0 := 0; __temp_0 < 10; __temp_0++ {
 for _, q := range(Q) {
 dieseldist.SendIntArray(distance[:], 0, q);
 }
i = 0;
for _, q := range(Q) {
 mystart = i*NodesPerThread;
myend = (i+1)*NodesPerThread;
lastthread = dieseldist.ConvBool(i==Num_threads-1);
if lastthread != 0 {
 myend = Num_nodes;
 }
dieseldist.ReceiveIntArray(slice[:], 0, q);
mysize = myend-mystart;
j = 0;
for __temp_1 := 0; __temp_1 < mysize; __temp_1++ {
 _temp_index_1 := j;
newDist=slice[_temp_index_1];
_temp_index_2 := mystart+j;
distance[_temp_index_2]=newDist;
j = j+1;
 }
i = i+1;
 }
 }
 dieseldist.EndTiming() ;
DistGlobal = distance;


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
var edges [36682000]int;
var inlinks [36682]int;
var distances [36682]int;
var distance int;
var newDistance [10000]int;
var condition int;
var inlink int;
var neighbor int;
var nodeInlinks int;
var i int;
var mystart int;
var myend int;
var cur int;
var temp int;
var temp1 int;
var temp2 int;
var mysize int;
dieseldist.ReceiveInt(&mystart, tid, 0);
dieseldist.ReceiveInt(&myend, tid, 0);
edges = Edges;
inlinks = Inlinks;
for __temp_2 := 0; __temp_2 < 10; __temp_2++ {
 dieseldist.ReceiveIntArray(distances[:], tid, 0);
mysize = myend-mystart;
i = 0;
for __temp_3 := 0; __temp_3 < mysize; __temp_3++ {
 cur = mystart+i;
_temp_index_1 := cur;
nodeInlinks=inlinks[_temp_index_1];
_temp_index_2 := cur;
distance=distances[_temp_index_2];
inlink = 0;
for __temp_4 := 0; __temp_4 < nodeInlinks; __temp_4++ {
 _temp_index_3 := cur*1000+inlink;
neighbor=edges[_temp_index_3];
_temp_index_4 := neighbor;
temp=distances[_temp_index_4];
condition = dieseldist.ConvBool(temp==1);
temp1 = 1;
temp2 = 2;
temp_bool_5:= condition; if temp_bool_5 != 0 { distance  = temp1 } else { distance = temp2 };
inlink = inlink+1;
 }
temp = distance;
_temp_index_6 := i;
newDistance[_temp_index_6]=temp;
i = i+1;
 }
dieseldist.SendIntArray(newDistance[:], tid, 0);
 }

  dieseldist.CleanupMain()
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
