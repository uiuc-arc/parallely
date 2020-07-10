package main

import (
	"dieseldist"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// import "time"
// import "os"

var Num_threads int
var Edges [6258600]int
var Inlinks [62586]int
var Outlinks [62586]int
var DistGlobal [62586]int
var Num_nodes int
var Num_edges int
var NodesPerThread int

func max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func convertToFloat(x int) float64 {
	return float64(x)
}

var Q = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func func_0() {
	dieseldist.InitQueues(11, "amqp://guest:guest@localhost:5672/")

	var DynMap [72587]dieseldist.ProbInterval
	var my_chan_index int
	_ = my_chan_index
	_ = DynMap
	var distance [62586]int
	dieseldist.InitDynArray(0, 62586, DynMap[:])
	var slice [10000]int
	dieseldist.InitDynArray(62586, 10000, DynMap[:])
	// var newDist int
	DynMap[72586] = dieseldist.ProbInterval{1, 0}
	var mystart int
	var myend int
	var i int
	// var j int
	var lastthread int

	// var mysize int
	distance = DistGlobal

	fmt.Println("In thread 0: ", distance[0])

	dieseldist.InitDynArray(0, 62586, DynMap[:])

	i = 0
	for _, q := range Q {
		mystart = i * NodesPerThread
		myend = (i + 1) * NodesPerThread
		lastthread = dieseldist.ConvBool(i == Num_threads-1)
		if lastthread != 0 {
			myend = Num_nodes
		}
		dieseldist.SendInt(mystart, 0, q)
		dieseldist.SendInt(myend, 0, q)
		i = i + 1
	}

	fmt.Println("In thread 0: Sending the distance array")

	// for __temp_0 := 0; __temp_0 < 10; __temp_0++ {
	for _, q := range Q {
		dieseldist.SendIntArray(distance[:], 0, q)
	}
	// i = 0
	for _, q := range Q {
		// 	mystart = i * NodesPerThread
		// 	myend = (i + 1) * NodesPerThread
		// 	lastthread = diesel.ConvBool(i == Num_threads-1)
		// 	if lastthread != 0 {
		// 		myend = Num_nodes
		// 	}
		dieseldist.ReceiveDynIntArray(slice[:], 0, q, DynMap[:], 62586)
		fmt.Println(slice[0])
		fmt.Println(DynMap[62586])
		// 	mysize = myend - mystart
		// 	j = 0
		// 	for __temp_1 := 0; __temp_1 < mysize; __temp_1++ {
		// 		_temp_index_1 := j
		// 		newDist = slice[_temp_index_1]
		// 		DynMap[72586] = DynMap[62586+_temp_index_1]
		// 		_temp_index_2 := mystart + j
		// 		distance[_temp_index_2] = newDist
		// 		DynMap[0+_temp_index_2] = DynMap[72586]
		// 		j = j + 1
		// 	}
		// 	i = i + 1
		// }
		// diesel.PrintWorstElement(DynMap[:], 0, 62586)
	}
	// distglobal = distance

	dieseldist.Cleanup()

	fmt.Println("Ending thread : ", 0)
}

func main() {
	fmt.Println("Starting main thread")

	Num_threads = 11

	data_bytes, _ := ioutil.ReadFile("../../inputs/p2p-Gnutella31.txt")
	Num_nodes = 62586
	Num_edges = Num_nodes * 1000

	fmt.Println("Starting reading the file")
	data_string := string(data_bytes)
	data_str_array := strings.Split(data_string, "\n")

	fmt.Println("Setting up the data structures")

	for i := range Inlinks {
		Inlinks[i] = 0
		Outlinks[i] = 0
		DistGlobal[i] = 0
	}
	DistGlobal[0] = 1

	NodesPerThread = Num_nodes / Num_threads

	fmt.Println("Populating the data structures")
	for i := 1; i < len(data_str_array)-1; i++ {
		elements := strings.Fields(data_str_array[i])
		index1, _ := strconv.Atoi(elements[0])
		index2, _ := strconv.Atoi(elements[1])

		Edges[(index2*100)+Inlinks[index2]] = index1
		Inlinks[index2]++
		Outlinks[index1]++
	}

	fmt.Println("Number of worker threads: ", Num_threads)
	fmt.Println("Number of nodes: ", len(DistGlobal))
	fmt.Println("Size of Inlinks: ", len(Inlinks))

	fmt.Println("Starting the iterations")
	// startTime := time.Now()

	func_0()

	// go func_0()
	// for _, index := range Q {
	// 	go func_Q(index)
	// }

	// fmt.Println("Main thread waiting for others to finish")
	// diesel.Wg.Wait()
	// elapsed := time.Since(startTime)

	// fmt.Println("Done!")
	// fmt.Println("Elapsed time : ", elapsed.Nanoseconds())
	// diesel.PrintMemory()
	// f, _ := os.Create("output.txt")
	// defer f.Close()

	// for i := range DistGlobal {
	// 	f.WriteString(fmt.Sprintln(DistGlobal[i]))
	// }
}
