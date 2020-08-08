package main

import "fmt"
import "diesel"
import "time"
import "os"
import  "io/ioutil"
import  "strings"
import  "strconv"

var Num_threads int
var Edges [811400]int
var Inlinks [8114]int
var Outlinks [8114]int
var DistGlobal [8114]int
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
  defer diesel.Wg.Done();
  var DynMap [8118]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var distance [8114]int;
diesel.InitDynArray(0, 8114, DynMap[:]);
var i int;
var cur int;
var inlink int;
var nodeInlinks int;
var neighbor int;
var tempdist int;
DynMap[8114] = diesel.ProbInterval{1, 0};
var temp1 int;
DynMap[8115] = diesel.ProbInterval{1, 0};
var temp int;
DynMap[8116] = diesel.ProbInterval{1, 0};
var condition int;
DynMap[8117] = diesel.ProbInterval{1, 0};
distance=DistGlobal;
diesel.InitDynArray(0, 8114, DynMap[:]);
for __temp_0 := 0; __temp_0 < 10; __temp_0++ {
 i = 0;
for __temp_1 := 0; __temp_1 < 8114; __temp_1++ {
 cur = i;
_temp_index_1 := cur;
nodeInlinks=Inlinks[_temp_index_1];
DynMap[8114] = diesel.ProbInterval{1, 0};
tempdist = 0;
inlink = 0;
for __temp_2 := 0; __temp_2 < nodeInlinks; __temp_2++ {
 _temp_index_2 := cur*100+inlink;
neighbor=Edges[_temp_index_2];
_temp_index_3 := neighbor;
temp=distance[_temp_index_3];
DynMap[8116] = DynMap[0 + _temp_index_3];
DynMap[8117].Reliability = DynMap[8116].Reliability;
condition = diesel.ConvBool(temp==1);
DynMap[8115] = diesel.ProbInterval{1, 0};
temp1 = 1;
temp_bool_4:= condition; if temp_bool_4 != 0 { tempdist  = temp1 } else { tempdist = tempdist };
if temp_bool_4 != 0 {DynMap[8114].Reliability  = DynMap[8117].Reliability * DynMap[8115].Reliability} else { DynMap[8114].Reliability = DynMap[8117].Reliability * DynMap[8114].Reliability};
inlink = inlink+1;
 }
_temp_index_5 := i;
distance[_temp_index_5]=tempdist;
DynMap[0 + _temp_index_5] = DynMap[8114];
i = i+1;
 }
 }


  fmt.Println("Ending thread : ", 0);
}

func main() {
	fmt.Println("Starting main thread");

  Num_threads = 1;
	
	diesel.InitChannels(1);
	
  data_bytes, err := ioutil.ReadFile("../../inputs/p2p-Gnutella09.txt")

  if err != nil {
     fmt.Println("[ERROR] Input does not exist")
     os.Exit(-1)
  }

  Num_nodes = 8114 // strconv.Atoi(os.Args[2])
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

    if index1 >= Num_nodes || index2 >= Num_nodes {
       continue;
    }

    Edges[(index2 * 100) + Inlinks[index2]] = index1
    Inlinks[index2]++
    Outlinks[index1]++
		// fmt.Println("---------------")
  }

	max := 0
	for i := range Inlinks{
		if Inlinks[i] > max {
			max = Inlinks[i]
		}
	}

  fmt.Println("Number of worker threads: ", Num_threads);
  fmt.Println("Number of nodes: ", len(DistGlobal));
  fmt.Println("Maximum degree: ", max);

  fmt.Println("Starting the iterations")
  startTime := time.Now()

	go func_0();


	fmt.Println("Main thread waiting for others to finish");  
	diesel.Wg.Wait()
  elapsed := time.Since(startTime)

	fmt.Println("Done!");
  fmt.Println("Elapsed time : ", elapsed.Nanoseconds());
	diesel.PrintMemory()
  f, _ := os.Create("output.txt")
  defer f.Close()

  for i := range DistGlobal{
    f.WriteString(fmt.Sprintln(DistGlobal[i]))
  }
}
