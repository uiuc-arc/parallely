package main

import "fmt"
import "dieseldist"
import "os"
import "math"
import  "io/ioutil"
import  "strings"
import  "strconv"

var Num_threads int
var Edges [6258600]int
var Inlinks [62586]int
var Outlinks [62586]int
var PagerankGlobal [62586]float64
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
	
var Q = []int {1,2,3,4,5};


func func_0() {
  dieseldist.InitQueues(Num_threads, "amqp://guest:guest@localhost:5672/")
  dieseldist.WaitForWorkers(Num_threads)
  var DynMap [0]dieseldist.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var pageranks [62586]float64;
var newPagerank float64;
var slice [15000]float64;
var mystart int;
var myend int;
var i int;
var j int;
var lastthread int;
var mysize int;
pageranks = PagerankGlobal;
 dieseldist.StartTiming() ;
i = 0;
for _, q := range(Q) {
 mystart = i*NodesPerThread;
myend = (i+1)*NodesPerThread;
lastthread = dieseldist.ConvBool(i==(Num_threads-1));
if lastthread != 0 {
 myend = Num_nodes;
 }
dieseldist.SendInt(mystart, 0, q);
dieseldist.SendInt(myend, 0, q);
i = i+1;
 }
for __temp_0 := 0; __temp_0 < 10; __temp_0++ {
 for _, q := range(Q) {
 dieseldist.SendFloat64Array(pageranks[:], 0, q);
 }
i = 0;
for _, q := range(Q) {
 mystart = i*NodesPerThread;
myend = (i+1)*NodesPerThread;
lastthread = dieseldist.ConvBool(i==(Num_threads-1));
if lastthread != 0 {
 myend = Num_nodes;
 }
mysize = myend-mystart;
j = 0;
dieseldist.ReceiveFloat64Array(slice[:], 0, q);
for __temp_1 := 0; __temp_1 < mysize; __temp_1++ {
 _temp_index_1 := j;
newPagerank=slice[_temp_index_1];
_temp_index_2 := mystart+j;
pageranks[_temp_index_2]=newPagerank;
j = j+1;
 }
i = i+1;
 }
 }
 dieseldist.EndTiming() ;
PagerankGlobal = pageranks;


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
var edges [6258600]int;
var inlinks [62586]int;
var outlinks [62586]int;
var pageranks [62586]float64;
var inlink int;
var neighbor int;
var outN int;
var outNf float64;
var current float64;
var newPagerank [15000]float64;
var nodeInlinks int;
var i int;
var mystart int;
var myend int;
var cur int;
var temp0 float64;
var mysize int;
edges = Edges;
inlinks = Inlinks;
outlinks = Outlinks;
dieseldist.ReceiveInt(&mystart, tid, 0);
dieseldist.ReceiveInt(&myend, tid, 0);
for __temp_2 := 0; __temp_2 < 10; __temp_2++ {
 dieseldist.ReceiveFloat64Array(pageranks[:], tid, 0);
inlink = 0;
mysize = myend-mystart;
i = 0;
for __temp_3 := 0; __temp_3 < mysize; __temp_3++ {
 cur = mystart+i;
_temp_index_1 := cur;
nodeInlinks=inlinks[_temp_index_1];
_temp_index_2 := i;
newPagerank[_temp_index_2]=0.15;
_temp_index_3 := i;
temp0=newPagerank[_temp_index_3];
for __temp_4 := 0; __temp_4 < nodeInlinks; __temp_4++ {
 _temp_index_4 := cur*100+inlink;
neighbor=edges[_temp_index_4];
_temp_index_5 := neighbor;
outN=outlinks[_temp_index_5];
outNf=convertToFloat(outN);
_temp_index_6 := neighbor;
current=pageranks[_temp_index_6];
temp0 = temp0+0.85*current/outNf;
inlink = inlink+1;
 }
_temp_index_7 := i;
newPagerank[_temp_index_7]=temp0;
i = i+1;
 }
dieseldist.SendFloat64Array(newPagerank[:], tid, 0);
 }

  fmt.Println("Ending thread : ", q);
}

func main() {
	tid, _ := strconv.Atoi(os.Args[1])	
	fmt.Println("Starting worker thread: ", tid)

  Num_threads = 6;
	
  data_bytes, err := ioutil.ReadFile("../../inputs/p2p-Gnutella31.txt")
  if err != nil {
     fmt.Println("[ERROR] Input does not exist")
     os.Exit(-1)
  }
  
  Num_nodes = int(math.Abs(62586)) // strconv.Atoi(os.Args[2])
  Num_edges = 6258600

  // fmt.Println("Starting reading the file")
  data_string := string(data_bytes)
  data_str_array := strings.Split(data_string, "\n")

  fmt.Println("Setting up the data structures")
  // Edges = make([]int, Num_nodes*1000)
  // Inlinks = make([]int, Num_nodes)
  // Outlinks = make([]int, Num_nodes)
  // PagerankGlobal = make([]float64, Num_nodes)

  for i := range Inlinks{
    Inlinks[i] = 0
    Outlinks[i] = 0
    PagerankGlobal[i] = 0.15
  }

  NodesPerThread = Num_nodes/Num_threads;

  // fmt.Println("Populating the data structures")
  for i := 1; i<len(data_str_array)-1 ; i++ {
    elements := strings.Fields(data_str_array[i])
    index1, _ := strconv.Atoi(elements[0])
    index2, _ := strconv.Atoi(elements[1])

    Edges[(index2 * 100) + Inlinks[index2]] = index1
    Inlinks[index2]++
    Outlinks[index1]++
		// fmt.Println("---------------")
  }

  // fmt.Println("Number of worker threads: ", Num_threads);
  // fmt.Println("Number of nodes: ", len(PagerankGlobal));
  // fmt.Println("Size of Inlinks: ", len(Inlinks));

	func_Q(tid);

}
