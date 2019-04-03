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
	"math/bits"
)

func sssp_func(iterations int, W [][]int, inlinks []int, outlinks []int,
	mystart, myend int, channel chan []float64, result_channel chan []float64,
	sigchannel chan bool, ackchannel chan bool, tid int){

	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// For now running a max number of iterations
	// Might not cover everything or end up running longer than needed
	for myiteration := 0; myiteration < iterations; myiteration++{
		distance := <- channel
		for j := mystart; j<myend; j++ {
			for k := 0; k<inlinks[j]; k++ {
				neighbor := W[j][k]

				if distance[j] > distance[neighbor] + 1 {
					distance[j] = distance[neighbor] + 1
				}


			}
		}
		for i := mystart; i < myend; i++ {
			// result_channel <- distance[i]
			parity := bits.OnesCount64(math.Float64bits(distance[i]))

			redo := false
			for !redo {
				// Failing with prob
				if randGen.Float64()<0.001 {
					// fmt.Println("Message failed")
					// int_val := math.Float64bits(-1)
					result_channel <- []float64{-1.0, float64(parity)}
					redo = <- ackchannel
				} else {
					// int_val := math.Float64bits(distance[i])
					result_channel <-  []float64{distance[i], float64(parity)}
					redo = <- ackchannel
				}				
			}
		}
	}
}

func main() {
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
	distance := make([]float64, num_nodes)
	
	for i := range W{
		W[i] = make([]int, num_nodes)
		inlinks[i] = 0
		outlinks[i] = 0
		distance[i] = math.MaxInt32
	}
	distance[0] = 1

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
	channels := make([]chan []float64, num_threads)
	sigchannels := make([]chan bool, num_threads)
	ackchannels := make([]chan bool, num_threads)
	reschannels := make([]chan []float64, num_threads)
	
	for i := range channels {
		channels[i] = make(chan []float64, 100)
		reschannels[i] = make(chan []float64, 100)
		sigchannels[i] = make(chan bool, 100)
		ackchannels[i] = make(chan bool, 100)
	}
	iterations := 1000

	// func pagerank(iterations int, W [][]int, inlinks []int, outlinks []int, mystart, myend int, channel chan []float64){
	start := time.Now()
	for i := range channels {
		t_start := (num_nodes/num_threads) * i
		t_end := (num_nodes/num_threads) * (i+1)
		go sssp_func(iterations, W, inlinks, outlinks, t_start, t_end, channels[i],
			reschannels[i], sigchannels[i], ackchannels[i], i)
	}

	k := 0
	
	fmt.Println("Starting the iterations")
	for iter:=0; iter < iterations; iter++{
		fmt.Println("Iteration : ", iter)
		for i := range channels {
			t_start := (num_nodes/num_threads) * i
			t_end := (num_nodes/num_threads) * (i+1)
			channels[i] <- distance
			// fmt.Println(i, len(results), t_start, t_end, t_end-t_start)
			
			for j:=t_start; j<t_end; j++ {
				result := <- reschannels[i]
				parity := float64(bits.OnesCount64(math.Float64bits(result[0])))
				for result[1] != parity {
					fmt.Println("Failed")
					ackchannels[i] <- false
					result = <- reschannels[i]
					parity = float64(bits.OnesCount64(math.Float64bits(result[0])))
					k++
				}
				// Parity check
				// parity := bits.OnesCount64(result[0])
				ackchannels[i] <- true
				distance[j] = float64(result[0])
			}
		}
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
