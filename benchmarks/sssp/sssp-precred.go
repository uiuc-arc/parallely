package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"math"
	// "time"
	"strconv"
)

func sliceTof32(slice []float64) []float32 {
  length := len(slice)
  result := make([]float32, length)
  for i := 0; i < length; i++ {
    result[i] = float32(slice[i])
  }
  return result
}

func sliceTof64(slice []float32) []float64 {
  length := len(slice)
  result := make([]float64, length)
  for i := 0; i < length; i++ {
    result[i] = float64(slice[i])
  }
  return result
}

func sssp_func(iterations int, W [][]int, inlinks []int, outlinks []int, mystart, myend int, channel chan []float32){
	// For now running a max number of iterations
	// Might not cover everything or end up running longer than needed
	for myiteration := 0; myiteration < iterations; myiteration++{
	  distance32 := <- channel
	  distance := sliceTof64(distance32)
		for j := mystart; j<myend; j++ {
			for k := 0; k<inlinks[j]; k++ {
				neighbor := W[j][k]

				if distance[j] > distance[neighbor] + 1 {
					distance[j] = distance[neighbor] + 1
				}
				
			}
		}
		result := sliceTof32(distance[mystart:myend])
		channel <- result
	}
}

func main() {
	argsWithoutProg := os.Args[1:]

	data_bytes, _ := ioutil.ReadFile(argsWithoutProg[0])
	num_nodes, _ := strconv.Atoi(argsWithoutProg[1])
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
	channels := make([]chan []float32, num_threads)
	for i := range channels {
		channels[i] = make(chan []float32, 100)
	}
	iterations := 10

	// func pagerank(iterations int, W [][]int, inlinks []int, outlinks []int, mystart, myend int, channel chan []float64){
	for i := range channels {
		t_start := (num_nodes/num_threads) * i
		t_end := (num_nodes/num_threads) * (i+1)
		go sssp_func(iterations, W, inlinks, outlinks, t_start, t_end, channels[i])
	}
	
	fmt.Println("Starting the iterations")
	for iter:=0; iter < iterations; iter++{
		fmt.Println("Iteration : ", iter)
		for i := range channels {
			t_start := (num_nodes/num_threads) * i
			t_end := (num_nodes/num_threads) * (i+1)
			distance32 := sliceTof32(distance)
			channels[i] <- distance32
			results32 := <- channels[i]
			results := sliceTof64(results32)
			// fmt.Println(i, len(results), t_start, t_end, t_end-t_start)
			
			k := 0
			for j:=t_start; j<t_end; j++ {
				distance[j] = results[k]
				k++
			}
		}
	}

	f, _ := os.Create("output.txt")
	defer f.Close()
	
	for i := range distance{
		f.WriteString(fmt.Sprintln(distance[i]))	
	}
}
