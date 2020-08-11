package main

import "fmt"
import "parallely"
import "time"
import "os"
import  "io/ioutil"
import  "strings"
import  "strconv"

var Num_threads int
var Edges [10909200]int
var Inlinks [1090920]int
var Outlinks [1090920]int
var PagerankGlobal [1090920]float64
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

var Q = []int {1,2,3,4,5,6,7,8};


func func_0() {
  defer parallely.Wg.Done();
  var DynMap [1290921]float64;
  _ = DynMap;
  var pageranks [1090920]float64;
	parallely.InitDynArray(0, 1090920, DynMap[:]);
	var newPagerank float64;
	DynMap[1090920] = 1;
	var slice [200000]float64;
	parallely.InitDynArray(1090921, 200000, DynMap[:]);
	var mystart int;
	var myend int;
	var i int;
	var j int;
	var lastthread int;
	var mysize int;
	pageranks=PagerankGlobal;
	parallely.InitDynArray(0, 1090920, DynMap[:]);
	i = 0;
	for _, q := range(Q) {
		mystart = i*NodesPerThread;
		myend = (i+1)*NodesPerThread;
		lastthread = parallely.ConvBool(i==(Num_threads-1));
		if lastthread != 0 {
			myend = Num_nodes;
		}
		parallely.SendInt(mystart, 0, q);
		parallely.SendInt(myend, 0, q);
		i = i+1;
	}
	for __temp_0 := 0; __temp_0 < 10; __temp_0++ {
		for _, q := range(Q) {
			parallely.SendDynFloat64ArrayO1(pageranks[:], 0, q, DynMap[:], 0);
		}
		i = 0;
		for _, q := range(Q) {
			mystart = i*NodesPerThread;
			myend = (i+1)*NodesPerThread;
			lastthread = parallely.ConvBool(i==(Num_threads-1));
			if lastthread != 0 {
				myend = Num_nodes;
			}
			mysize = myend-mystart;
			j = 0;
			parallely.ReceiveDynFloat64ArrayO1(slice[:], 0, q, DynMap[:], 1090921);
			for __temp_1 := 0; __temp_1 < mysize; __temp_1++ {
				_temp_index_1 := j;
				newPagerank=slice[_temp_index_1];
				DynMap[1090920] = DynMap[1090921 + _temp_index_1];
				_temp_index_2 := mystart+j;
				pageranks[_temp_index_2]=newPagerank;
				DynMap[0 + _temp_index_2] = DynMap[1090920];
				j = j+1;
			}
			i = i+1;
		}
		parallely.DumpDynMap(DynMap[:], "dynmap" + strconv.Itoa(__temp_0));
	}

	fmt.Println("----------------------------");

	fmt.Println("Spec checkarray(pageranks, 0.99): ", parallely.CheckArray(0, 0.99, 1090920, DynMap[:]));

	fmt.Println("----------------------------");

	PagerankGlobal = pageranks;

  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  defer parallely.Wg.Done();
  var DynMap [1290923]float64;
  _ = DynMap;
  q := tid;
	var edges [10909200]int;
	var inlinks [1090920]int;
	var outlinks [1090920]int;
	var pageranks [1090920]float64;
	parallely.InitDynArray(0, 1090920, DynMap[:]);
	var inlink int;
	var neighbor int;
	var outN int;
	var outNf float64;
	var current float64;
	DynMap[1090920] = 1;
	var newPagerank [200000]float64;
	parallely.InitDynArray(1090921, 200000, DynMap[:]);
	var nodeInlinks int;
	var i int;
	var mystart int;
	var myend int;
	var cur int;
	DynMap[1290921] = 1;
	var temp0 float64;
	DynMap[1290922] = 1;
	var mysize int;
	edges = Edges;
	inlinks = Inlinks;
	outlinks = Outlinks;
	parallely.ReceiveInt(&mystart, tid, 0);
	parallely.ReceiveInt(&myend, tid, 0);
	for __temp_2 := 0; __temp_2 < 10; __temp_2++ {
		parallely.ReceiveDynFloat64ArrayO1(pageranks[:], tid, 0, DynMap[:], 0);
		inlink = 0;
		mysize = myend-mystart;
		i = 0;
		for __temp_3 := 0; __temp_3 < mysize; __temp_3++ {
			cur = mystart+i;
			DynMap[1290921] = 1;
			_temp_index_1 := cur;
			nodeInlinks=inlinks[_temp_index_1];
			_temp_index_2 := i;
			newPagerank[_temp_index_2]=0.15;
			DynMap[1090921 + _temp_index_2] = 1;
			_temp_index_3 := i;
			temp0=newPagerank[_temp_index_3];
			DynMap[1290922] = DynMap[1090921 + _temp_index_3];
			for __temp_4 := 0; __temp_4 < nodeInlinks; __temp_4++ {
				_temp_index_4 := cur*10+inlink;
				neighbor=edges[_temp_index_4];
				_temp_index_5 := neighbor;
				outN=outlinks[_temp_index_5];
				outNf=convertToFloat(outN);
				_temp_index_6 := neighbor;
				current=pageranks[_temp_index_6];
				DynMap[1090920] = DynMap[0 + _temp_index_6];
				temp0 = temp0+0.85*current/outNf;
				DynMap[1290922] = parallely.Max(0.0, DynMap[1090920] + DynMap[1290922] - float64(1));
				inlink = inlink+1;
			}
			temp0 = parallely.RandchoiceFloat64(float32(0.999999999), temp0, -1);
			DynMap[1290922] = DynMap[1290922] * 0.999999999;
			_temp_index_7 := i;
			newPagerank[_temp_index_7]=temp0;
			DynMap[1090921 + _temp_index_7] = DynMap[1290922];
			i = i+1;
		}
		parallely.SendDynFloat64ArrayO1(newPagerank[:], tid, 0, DynMap[:], 1090921);
	}

  fmt.Println("Ending thread : ", q);
}

func main() {
	fmt.Println("Starting main thread");

  Num_threads = 9;
	
	parallely.InitChannels(9);
	
  data_bytes, err := ioutil.ReadFile("../../inputs/roadNet-PA.txt")
  if err != nil {
		fmt.Println("[ERROR] Input does not exist")
		os.Exit(-1)
  }
  
  Num_nodes = 1090920 // strconv.Atoi(os.Args[2])
  Num_edges = Num_nodes * 10

  fmt.Println("Starting reading the file")
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
	node := 0
  max_degree := 0
  
  fmt.Println("Populating the data structures")
  for i := 1; i<len(data_str_array)-1 ; i++ {
    elements := strings.Fields(data_str_array[i])
    index1, _ := strconv.Atoi(elements[0])
    index2, _ := strconv.Atoi(elements[1])

    Edges[(index2 * 10) + Inlinks[index2]] = index1
    Inlinks[index2]++
    Outlinks[index1]++
		// fmt.Println("---------------")
    if Inlinks[index2]>max_degree {
			max_degree = Inlinks[index2]
			node =  index2
		}
  }

	fmt.

	fmt.Println("Max Degree : ", node, max_degree)

  fmt.Println("Number of worker threads: ", Num_threads);
  fmt.Println("Number of nodes: ", len(PagerankGlobal));
  fmt.Println("Size of Inlinks: ", len(Inlinks));

  fmt.Println("Starting the iterations")
  startTime := time.Now()

	go func_0();
	for _, index := range Q {
		go func_Q(index);
	}


	fmt.Println("Main thread waiting for others to finish");  
	parallely.Wg.Wait()
  elapsed := time.Since(startTime)
	parallely.PrintMemory()
	fmt.Println("Done!");
  fmt.Println("Elapsed time : ", elapsed.Nanoseconds());

  f, _ := os.Create("output.txt")
  defer f.Close()

  for i := range PagerankGlobal{
    f.WriteString(fmt.Sprintln(PagerankGlobal[i]))
  }
}
