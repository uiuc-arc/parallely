package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"math"
	"time"
	"strconv"
	"math/rand"
	."dynfloats"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func sssp_func(iterations int, W [][]int, inlinks []int, outlinks []int,
	mystart, myend int, channel_in, channel_out chan DynRelyInt, 
	sigchannel_in, sigchannel_out chan bool, tid, num_nodes int){

	// randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	distance := make([]DynRelyInt, num_nodes)
	
	// For now running a max number of iterations
	// Might not cover everything or end up running longer than needed
	for myiteration := 0; myiteration < iterations; myiteration++{
		for i := 0; i < num_nodes; i++ {
			distance[i] = <- channel_in
		}
		
		for j := mystart; j<myend; j++ {
			for k := 0; k<inlinks[j]; k++ {
				neighbor := W[j][k]
				if EqualconstIntRely(distance[neighbor], 1) {
					distance[j] = DynRelyInt{1, 1}
				}
			}
		}
		for i := mystart; i < myend; i++ {
			DynSendDynInt(sigchannel_out, channel_out, distance[i], 0.00001) 
		}
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	argsWithoutProg := os.Args[1:]

	data_bytes, _ := ioutil.ReadFile(argsWithoutProg[0])
	num_nodes, _ := strconv.Atoi(argsWithoutProg[1])
	outfile := argsWithoutProg[2]
	// num_edges, _ := strconv.Atoi(argsWithoutProg[2])

	fmt.Println("Starting reading the file")
	data_string := string(data_bytes)
	data_str_array := strings.Split(data_string, "\n")
	
	fmt.Println("Setting up the data structures")
	W := make([][]int, num_nodes)
	inlinks := make([]int, num_nodes)
	outlinks := make([]int, num_nodes)
	distance := make([]DynRelyInt, num_nodes)
	
	for i := range W{
		W[i] = make([]int, num_nodes)
		inlinks[i] = 0
		outlinks[i] = 0
		distance[i] = DynRelyInt{math.MaxInt32, 1}
	}
	distance[0] = DynRelyInt{1, 1}

	fmt.Println("Populating the data structures")
	// Nodes: 62586 Edges: 147892
	for i := 1; i<len(data_str_array)-1 ; i++ {
		// fmt.Println(data_str_array[i])
		elements := strings.Fields(data_str_array[i])
		// fmt.Println(elements[0])
		index1, _ := strconv.Atoi(elements[0])
		index2, _ := strconv.Atoi(elements[1])
		
		W[index2][inlinks[index2]] = index1
		inlinks[index2]++
		outlinks[index1]++		
	}

	fmt.Println("Finished populating the data structures")

	num_threads := 4
	channels_main_threads := make([]chan DynRelyInt, num_threads)
	channels_threads_main := make([]chan DynRelyInt, num_threads)
	
	sigchannels_in := make([]chan bool, num_threads)
	sigchannels_out := make([]chan bool, num_threads)	
	
	for i := range channels_main_threads {
		channels_main_threads[i] = make(chan DynRelyInt, 100)
		channels_threads_main[i] = make(chan DynRelyInt, 100)
		sigchannels_in[i] = make(chan bool, 100)
		sigchannels_out[i] = make(chan bool, 100)
	}
	iterations := 20

	start := time.Now()

	for i := range channels_main_threads {
		t_start := (num_nodes/num_threads) * i
		t_end := (num_nodes/num_threads) * (i+1)
		if i == num_threads {
			t_end = max(num_nodes, (num_nodes/num_threads) * (i+1))
		}
		go sssp_func(iterations, W, inlinks, outlinks, t_start, t_end, channels_main_threads[i],
			channels_threads_main[i], sigchannels_in[i], sigchannels_out[i], i, num_nodes)
	}

	k := 0
	
	fmt.Println("Starting the iterations")
	// failed_once := false
	for iter:=0; iter < iterations; iter++{
		fmt.Println("Iteration : ", iter)
		for i := range channels_main_threads {
			t_start := (num_nodes/num_threads) * i
			t_end := (num_nodes/num_threads) * (i+1)
			if i == num_threads {
				t_end = max(num_nodes, (num_nodes/num_threads) * (i+1))
			}

			for j:=0; j<num_nodes; j++ {
				channels_main_threads[i] <- distance[j]
			}
			// fmt.Println(i, len(results), t_start, t_end, t_end-t_start)
			
			for j:=t_start; j<t_end; j++ {
				pass := <- sigchannels_out[i]
				if pass {
					distance[j] = <- channels_threads_main[i]
				}
				// else {
				// 	fmt.Println(j)
				// 	failed_once = true
				// }
			}
		}

		// f, _ := os.Create(fmt.Sprintf("_iter_%d.txt", iter))
		// f.WriteString(fmt.Sprintln(failed_once))
		// for i := range distance{
		// 	f.WriteString(fmt.Sprintln(distance[i]))
		// }
		// f.Close()
	}

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Println("Retries :", k)
	fmt.Println("Elapsed time :", elapsed.Nanoseconds())

	fmt.Println(distance[:5])
	fmt.Println(outfile)
	f, _ := os.Create(outfile)
	defer f.Close()
	
	for i := range distance{
		f.WriteString(fmt.Sprintln(distance[i]))	
	}
}
