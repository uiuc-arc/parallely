package main

import "fmt"
import "diesel"
import "time"
import "os"
import  "io/ioutil"
import  "strings"
import  "strconv"

var Num_threads int
var Edges [108809200]int
var Inlinks [1088092]int
var Outlinks [1088092]int
var PagerankGlobal [1088092]float64
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
  defer diesel.Wg.Done();
  var DynMap [1288093]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  var pageranks [1088092]float64;
	diesel.InitDynArray(0, 1088092, DynMap[:]);
	var newPagerank float64;
	DynMap[1088092] = diesel.ProbInterval{1, 0};
	var slice [200000]float64;
	diesel.InitDynArray(1088093, 200000, DynMap[:]);
	var mystart int;
	var myend int;
	var i int;
	var j int;
	var lastthread int;
	var mysize int;
	pageranks=PagerankGlobal;
	diesel.InitDynArray(0, 1088092, DynMap[:]);
	i = 0;
	for _, q := range(Q) {
		mystart = i*NodesPerThread;
		myend = (i+1)*NodesPerThread;
		lastthread = diesel.ConvBool(i==(Num_threads-1));
		if lastthread != 0 {
			myend = Num_nodes;
		}
		diesel.SendInt(mystart, 0, q);
		diesel.SendInt(myend, 0, q);
		i = i+1;
	}
	for __temp_0 := 0; __temp_0 < 10; __temp_0++ {
		for _, q := range(Q) {
			diesel.SendDynFloat64ArrayO1(pageranks[:], 0, q, DynMap[:], 0);
		}
		i = 0;
		for _, q := range(Q) {
			mystart = i*NodesPerThread;
			myend = (i+1)*NodesPerThread;
			lastthread = diesel.ConvBool(i==(Num_threads-1));
			if lastthread != 0 {
				myend = Num_nodes;
			}
			mysize = myend-mystart;
			j = 0;
			diesel.ReceiveDynFloat64ArrayO1(slice[:], 0, q, DynMap[:], 1088093);
			for __temp_1 := 0; __temp_1 < mysize; __temp_1++ {
				_temp_index_1 := j;
				newPagerank=slice[_temp_index_1];
				DynMap[1088092] = DynMap[1088093 + _temp_index_1];
				_temp_index_2 := mystart+j;
				pageranks[_temp_index_2]=newPagerank;
				DynMap[0 + _temp_index_2] = DynMap[1088092];
				j = j+1;
			}
			i = i+1;
		}
	}

	fmt.Println("----------------------------");

	fmt.Println("Spec checkarray(pageranks, 0.99): ", diesel.CheckArray(0, 0.99, 1088092, DynMap[:]));

	fmt.Println("----------------------------");

	PagerankGlobal = pageranks;


  fmt.Println("Ending thread : ", 0);
}
func func_Q(tid int) {
  defer diesel.Wg.Done();
  var DynMap [1288095]diesel.ProbInterval;
  var my_chan_index int;
  _ = my_chan_index;
  _ = DynMap;
  q := tid;
	var edges [108809200]int;
	var inlinks [1088092]int;
	var outlinks [1088092]int;
	var pageranks [1088092]float64;
	diesel.InitDynArray(0, 1088092, DynMap[:]);
	var inlink int;
	var neighbor int;
	var outN int;
	var outNf float64;
	var current float64;
	DynMap[1088092] = diesel.ProbInterval{1, 0};
	var newPagerank [200000]float64;
	diesel.InitDynArray(1088093, 200000, DynMap[:]);
	var nodeInlinks int;
	var i int;
	var mystart int;
	var myend int;
	var cur int;
	DynMap[1288093] = diesel.ProbInterval{1, 0};
	var temp0 float64;
	DynMap[1288094] = diesel.ProbInterval{1, 0};
	var mysize int;
	edges = Edges;
	inlinks = Inlinks;
	outlinks = Outlinks;
	diesel.ReceiveInt(&mystart, tid, 0);
	diesel.ReceiveInt(&myend, tid, 0);
	for __temp_2 := 0; __temp_2 < 10; __temp_2++ {
		diesel.ReceiveDynFloat64ArrayO1(pageranks[:], tid, 0, DynMap[:], 0);
		inlink = 0;
		mysize = myend-mystart;
		i = 0;
		for __temp_3 := 0; __temp_3 < mysize; __temp_3++ {
			DynMap[1288093] = diesel.ProbInterval{1, 0};
			cur = mystart+i;
			_temp_index_1 := cur;
			nodeInlinks=inlinks[_temp_index_1];
			_temp_index_2 := i;
			newPagerank[_temp_index_2]=0.15;
			DynMap[1088093 + _temp_index_2] = diesel.ProbInterval{1, 0};
			_temp_index_3 := i;
			temp0=newPagerank[_temp_index_3];
			DynMap[1288094] = DynMap[1088093 + _temp_index_3];
			for __temp_4 := 0; __temp_4 < nodeInlinks; __temp_4++ {
				_temp_index_4 := cur*100+inlink;
				neighbor=edges[_temp_index_4];
				_temp_index_5 := neighbor;
				outN=outlinks[_temp_index_5];
				outNf=convertToFloat(outN);
				_temp_index_6 := neighbor;
				current=pageranks[_temp_index_6];
				DynMap[1088092] = DynMap[0 + _temp_index_6];
				DynMap[1288094].Reliability = DynMap[1088092].Reliability + DynMap[1288094].Reliability - 1.0;
				DynMap[1288094].Delta = DynMap[1288094].Delta + DynMap[1088092].Delta;
				temp0 = temp0+0.85*current/outNf;
				inlink = inlink+1;
			}
			_temp_index_7 := i;
			newPagerank[_temp_index_7]=temp0;
			DynMap[1088093 + _temp_index_7] = DynMap[1288094];
			i = i+1;
		}
		diesel.SendDynFloat64ArrayO1(newPagerank[:], tid, 0, DynMap[:], 1088093);
	}

  fmt.Println("Ending thread : ", q);
}

func main() {
	fmt.Println("Starting main thread");

  Num_threads = 9;
	
	diesel.InitChannels(9);
	
  data_bytes, err := ioutil.ReadFile("../../inputs/roadNet-PA.txt")
  if err != nil {
		fmt.Println("[ERROR] Input does not exist")
		os.Exit(-1)
  }
  
  Num_nodes = 1088092 // strconv.Atoi(os.Args[2])
  Num_edges = Num_nodes * 100

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
	diesel.Wg.Wait()
  elapsed := time.Since(startTime)
	diesel.PrintMemory()
	fmt.Println("Done!");
  fmt.Println("Elapsed time : ", elapsed.Nanoseconds());

  f, _ := os.Create("output.txt")
  defer f.Close()

  for i := range PagerankGlobal{
    f.WriteString(fmt.Sprintln(PagerankGlobal[i]))
  }
}
