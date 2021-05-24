package main

import (
	"parallely"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var Num_threads int
var Edges [8114000]int
var Inlinks [8114]int
var Outlinks [8114]int

var VisitedGlobal [8114]int
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

var Q = []process{1, 2, 3, 4, 5, 6, 7, 8}

func func_0() {
	// defer diesel.Wg.Done()
	/*approx*/ var visited [8114]int
	/*approx*/ var tempVisited int
	/*approx*/ var slice [1200]int
	var mystart int
	var myend int
	var i int
	var iter int
	
	var j int
	var temp int
	var lastthread int
	var mysize int
	
	visited = visitedGlobal
	i = 0
	for _, q := range Q {
		mystart = i * NodesPerThread
		myend = (i + 1) * NodesPerThread
		temp = Num_threads - 1
		lastthread = (i == temp)		
		if lastthread != 0 {
			myend = Num_nodes
		}		
		send(q, mystart)
		send(q, myend)
		i = i + 1
	}
	
	for iter := 0; iter < 10; iter++ /*maxiterations=10*/ {
		for _, q := range Q {
			send(q, visited)
		}
		i = 0
		for _, q := range Q {
			mystart = i * NodesPerThread
			myend = (i + 1) * NodesPerThread
			temp = Num_threads - 1
			lastthread = (i == temp)
			if lastthread != 0 {
				myend = Num_nodes
			}
			mysize = myend - mystart
			j = 0
			slice = receive(q)
			for j := 0; __temp_1 < mysize; __temp_1++ /*maxiterations=10*/ {
				tempVisited = slice[j]
				temp = mystart + j
				visited[temp] = tempVisited
				j = j + 1
			}
			i = i + 1
		}
	}
	visitedGlobal = visited
}

func func_Q(q int) {
	// defer diesel.Wg.Done()
	var edges [8114000]int
	var inlinks [8114]int
	var outlinks [8114]int
	/*approx*/ var visited [8114]float64
	
	var inlink int
	var neighbor int
	/*approx*/ var current int
	/*approx*/ var newVisited [1200]float64
	var nodeInlinks int
	var i int
	var mystart int
	var myend int
	var cur int
	var temp int
	/*approx*/ var atemp int
	/*approx*/ var temp0 int
	var mysize int

	var iter int
	
	edges = Edges
	inlinks = Inlinks
	outlinks = Outlinks
	
	mystart = receive(func_0)
	myend = receive(func_0)
	
	for iter := 0; iter < 10; iter++ /*maxiterations=10*/ {
		visited = receive(func_0)
		mysize = myend - mystart		
		for i := 0; i < mysize; i++ /*maxiterations=10*/ {
			cur = mystart + i
			nodeInlinks = inlinks[cur]			
			temp0 = 0
			
			for inlink := 0; inlink < nodeInlinks; inlink++ /*maxiterations=10*/ {
				temp = cur*1000 + inlink
				neighbor = edges[temp]
				outN = outlinks[neighbor]
				outNf = convertToFloat(outN)
				current = visited[neighbor]
				temp0 = current || temp0
				inlink = inlink + 1
			}
			visited[i] = temp0
			i = i + 1
		}
		send(func_0, visited) /*reliability=0.999*/
	}
}

func main() {
	fmt.Println("Starting main thread")

	Num_threads = 9

	diesel.InitChannels(9)

	data_bytes, err := ioutil.ReadFile("../../inputs/p2p-Gnutella09.txt")
	if err != nil {
		fmt.Println("[ERROR] Input does not exist")
		os.Exit(-1)
	}

	Num_nodes = int(math.Abs(8114))
	Num_edges = Num_nodes * 1000

	fmt.Println("Starting reading the file")
	data_string := string(data_bytes)
	data_str_array := strings.Split(data_string, "\n")

	fmt.Println("Setting up the data structures")

	for i := range Inlinks {
		Inlinks[i] = 0
		Outlinks[i] = 0
		VisitedGlobal[i] = 0
	}
	VisitedGlobal[0] = 1

	NodesPerThread = Num_nodes / Num_threads

	fmt.Println("Populating the data structures")
	for i := 1; i < len(data_str_array)-1; i++ {
		elements := strings.Fields(data_str_array[i])
		index1, _ := strconv.Atoi(elements[0])
		index2, _ := strconv.Atoi(elements[1])

		Edges[(index2*1000)+Inlinks[index2]] = index1
		Inlinks[index2]++
		Outlinks[index1]++
	}

	fmt.Println("Number of worker threads: ", Num_threads)
	fmt.Println("Number of nodes: ", len(PagerankGlobal))
	fmt.Println("Size of Inlinks: ", len(Inlinks))

	fmt.Println("Starting the iterations")
	startTime := time.Now()

	parallely.LaunchThread(0, func_0)
	parallely.LaunchThreadGroup(Q, func_Q, "q")

	fmt.Println("Main thread waiting for others to finish")
	diesel.Wg.Wait()
	elapsed := time.Since(startTime)

	fmt.Println("Done!")
	fmt.Println("Elapsed time : ", elapsed.Nanoseconds())

	// f, _ := os.Create("output.txt")
	// defer f.Close()

	// for i := range PagerankGlobal {
	// 	f.WriteString(fmt.Sprintln(PagerankGlobal[i]))
	// }
}
