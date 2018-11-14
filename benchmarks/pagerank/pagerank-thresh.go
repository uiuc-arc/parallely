package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"math"
	"time"
	"strconv"
)

func pagerank_func(iterations int, W [][]int, inlinks []int, outlinks []int, mystart, myend int, channel chan []float64){
	r := 0.15
	d := 0.85
	// mywork := myend-mystart
	
	for myiteration := 0; myiteration < iterations; myiteration++{
	  maxDiff := 0.0
		pageranks := <- channel
		for j := mystart; j<myend; j++ {
		  original := pageranks[j]
			pageranks[j] = r
			for k := 0; k<inlinks[j]; k++ {
				neighbor := W[j][k]
				pageranks[j] = pageranks[j] + d * pageranks[neighbor]/float64(outlinks[neighbor])
			}
			diff := math.Abs(original-pageranks[j])
			if diff>maxDiff {
			  maxDiff = diff
			}
		}
		var toSend []float64
		if maxDiff<0.1 {
		  //fmt.Println("below threshold",maxDiff)
		  toSend = make([]float64,0)
		} else {
		  toSend = pageranks[mystart:myend]
		}
		channel <- toSend
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
	pagerank := make([]float64, num_nodes)
	
	for i := range W{
		W[i] = make([]int, num_nodes)
		inlinks[i] = 0
		outlinks[i] = 0
		pagerank[i] = 0.15
	}

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
	for i := range channels {
		channels[i] = make(chan []float64, 100)
	}
	iterations := 10

	// func pagerank(iterations int, W [][]int, inlinks []int, outlinks []int, mystart, myend int, channel chan []float64){
	for i := range channels {
		t_start := (num_nodes/num_threads) * i
		t_end := (num_nodes/num_threads) * (i+1)
		go pagerank_func(iterations, W, inlinks, outlinks, t_start, t_end, channels[i])
	}
	
	startTime := time.Now()
	fmt.Println("Starting the iterations")
	for iter:=0; iter < iterations; iter++{
		fmt.Println("Iteration : ", iter)
		for i := range channels {
			t_start := (num_nodes/num_threads) * i
			t_end := (num_nodes/num_threads) * (i+1)
			channels[i] <- pagerank
			results := <- channels[i]
			fmt.Println(i, len(results), t_start, t_end, t_end-t_start)
			if len(results)>0 {
			  k := 0
			  for j:=t_start; j<t_end; j++ {
				  pagerank[j] = results[k]
				  k++
			  }
			}
		}
	}
	elapsed := time.Since(startTime)
	fmt.Println(elapsed)

	f, _ := os.Create("output.txt")
	defer f.Close()

	
	for i := range pagerank{
		f.WriteString(fmt.Sprintln(pagerank[i]))	
	}
}
