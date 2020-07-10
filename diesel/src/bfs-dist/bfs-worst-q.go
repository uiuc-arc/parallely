package main

import "fmt"
import "os"
import "dieseldist"
// import "io/ioutil"
// import"strings"
import "strconv"

var Num_threads int
var Edges [6258600]int
var Inlinks [62586]int
var Outlinks [62586]int
var DistGlobal [62586]int
var Num_nodes int
var Num_edges int
var NodesPerThread int

// func max(x, y int) int {
// 	if x > y 
// 		return x
// 	} else {
// 		return y
// 	}
// }

// func min(x, y int) int {
// 	if x < y 
// 		return x
// 	} else {
// 		return y
// 	}
// }

func convertToFloat(x int) float64 {
	return float64(x)
}

func func_q(tid int) {
	var DynMap [72592]dieseldist.ProbInterval;
	dieseldist.InitQueues(11, "amqp://guest:guest@localhost:5672/")
	var distances [62586]int;
	// diesel.InitDynArray(0, 62586, DynMap[:]);
	var mystart int
	var myend int

	var newDistance [10000]int;
	dieseldist.InitDynArray(62587, 10000, DynMap[:]);
	
	dieseldist.ReceiveInt(&mystart, tid, 0);
	dieseldist.ReceiveInt(&myend, tid, 0);

	// for __temp_2 := 0; __temp_2 < 10; __temp_2++ {
	dieseldist.ReceiveIntArray(distances[:], tid, 0);
	// }

	DynMap[62587] = dieseldist.ProbInterval{0.999999, 0.9999}
	newDistance[0] = 1234
	dieseldist.SendDynIntArray(newDistance[:], tid, 0, DynMap[:], 62587);

	fmt.Println(tid, mystart, myend, distances[0])

	dieseldist.Cleanup()
 
	fmt.Println(" Ending thread :", tid)
}    
  
func main() {
	tid, _ := strconv.Atoi(os.Args[1])
	
	fmt.Println("Starting worker thread: ", tid)

	Num_threads = 11
  
	// Num_nodes := 62586
	// Num_edges = Num_nodes * 1000

	func_q(tid)  
}
