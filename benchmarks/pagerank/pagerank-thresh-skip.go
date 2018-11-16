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

func pagerank_func(iterations int, W [][]int, inlinks []int, outlinks []int, mynode int, datachannel chan []float64, reschannel chan float64, datasigchannel chan bool, ressigchannel chan bool){
  r := 0.15
  d := 0.85

  for myiteration := 0; myiteration < iterations; myiteration++{
    gotwork := <- datasigchannel
    if gotwork {
      pageranks := <- datachannel
      mypagerank := r
      for k := 0; k<inlinks[mynode]; k++ {
        neighbor := W[mynode][k]
        //fmt.Println(mynode,myiteration,neighbor,len(pageranks),len(outlinks))
        mypagerank += d * pageranks[neighbor]/float64(outlinks[neighbor])
      }
      if math.Abs(mypagerank-pageranks[mynode]) >= 0.01 {
        ressigchannel <- true
        reschannel <- mypagerank
      } else {
        ressigchannel <- false
      }
    } else {
      ressigchannel <- false
    }
  }
}

func main() {
  data_bytes, _ := ioutil.ReadFile(os.Args[1])
  num_nodes, _ := strconv.Atoi(os.Args[2])
  iterations, _ := strconv.Atoi(os.Args[3])

  //fmt.Println("Starting reading the file")
  data_string := string(data_bytes)
  data_str_array := strings.Split(data_string, "\n")

  //fmt.Println("Setting up the data structures")
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

  //fmt.Println("Populating the data structures")
  for i := 1; i<len(data_str_array)-1 ; i++ {
    elements := strings.Fields(data_str_array[i])
    index1, _ := strconv.Atoi(elements[0])
    index2, _ := strconv.Atoi(elements[1])

    W[index2][inlinks[index2]] = index1
    inlinks[index2]++
    outlinks[index1]++
  }

  //fmt.Println("Finished populating the data structures")

  channels := make([]chan []float64, num_nodes)
  for i := range channels {
    channels[i] = make(chan []float64, 1)
  }
  reschannels := make([]chan float64, num_nodes)
  for i := range reschannels {
    reschannels[i] = make(chan float64, 1)
  }
  sigchannels := make([]chan bool, num_nodes)
  for i := range sigchannels {
    sigchannels[i] = make(chan bool, 1)
  }
  ressigchannels := make([]chan bool, num_nodes)
  for i := range ressigchannels {
    ressigchannels[i] = make(chan bool, 1)
  }

  for i := range channels {
    go pagerank_func(iterations, W, inlinks, outlinks, i, channels[i], reschannels[i], sigchannels[i], ressigchannels[i])
  }

  toSkip := make([]bool, num_nodes)
  dropped := 0
  //fmt.Println("Starting the iterations")
  startTime := time.Now()
  for iter:=0; iter < iterations; iter++{
    //fmt.Println("Iteration : ", iter)
    results := make([]float64, num_nodes)
    copy(results, pagerank)
    for i := range channels {
      if !toSkip[i] {
        sigchannels[i] <- true
        pagerankcopy := make([]float64, num_nodes)
        copy(pagerankcopy, pagerank)
        channels[i] <- pagerankcopy
      } else {
        sigchannels[i] <- false
      }
    }
    for i := range channels {
      gotresult := <- ressigchannels[i]
      if gotresult {
        results[i] = <- reschannels[i]
      } else {
        dropped += 1
        toSkip[i] = true
      }
    }
    pagerank = results
  }
  elapsed := time.Since(startTime)
  fmt.Println(elapsed)
  fmt.Println(dropped)

  f, _ := os.Create("output.txt")
  defer f.Close()

  for i := range pagerank{
    f.WriteString(fmt.Sprintln(pagerank[i]))
  }
}
