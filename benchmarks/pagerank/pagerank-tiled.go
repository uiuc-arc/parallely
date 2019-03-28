package main

import (
  "os"
  "fmt"
  "io/ioutil"
  "strings"
  "math/rand"
  "time"
  "strconv"
)

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func pagerank_func(iterations int, W [][]int, inlinks []int, outlinks []int,
	myfirstnode, mylastnode int, datachannel, reschannel chan []float64,
	datasigchannel chan bool, ressigchannel chan bool, ackchannel chan bool, redo bool){
  r := 0.15
  d := 0.85

  for myiteration := 0; myiteration < iterations; myiteration++{
    <- datasigchannel
    pageranks := <- datachannel
    mypageranks := make([]float64, mylastnode-myfirstnode)
    for node := range mypageranks{
      mypageranks[node] = r
      for k := 0; k<inlinks[myfirstnode+node]; k++ {
        neighbor := W[myfirstnode+node][k]
        //fmt.Println(myfirstnode+node,myiteration,neighbor,len(pageranks),len(outlinks))
        mypageranks[node] += d * pageranks[neighbor]/float64(outlinks[neighbor])
      }
    }

		if redo {
			pass := false
			for !pass {
				failure := rand.Float64()
				if failure >= 0.01 {
					ressigchannel <- true		
					reschannel <- mypageranks
				} else {
					ressigchannel <- false
				}
				pass = <- ackchannel 
			}
		} else {
			failure := rand.Float64()
			if failure >= 0.01 {
				ressigchannel <- true		
				reschannel <- mypageranks
			} else {
				ressigchannel <- false
			}
		}
  }
}


func main() {
  data_bytes, _ := ioutil.ReadFile(os.Args[1])
  num_nodes, _ := strconv.Atoi(os.Args[2])
  iterations, _ := strconv.Atoi(os.Args[3])
  num_threads, _ := strconv.Atoi(os.Args[4])
	redoin, _ := strconv.Atoi(os.Args[5])
	outfile := os.Args[6]

	var redo bool = redoin != 0

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

  channels := make([]chan []float64, num_threads)
  for i := range channels {
    channels[i] = make(chan []float64, 10)
  }
  reschannels := make([]chan []float64, num_threads)
  for i := range reschannels {
    reschannels[i] = make(chan []float64, 10)
  }
  sigchannels := make([]chan bool, num_threads)
  for i := range sigchannels {
    sigchannels[i] = make(chan bool, 10)
  }
  ressigchannels := make([]chan bool, num_threads)
  for i := range ressigchannels {
    ressigchannels[i] = make(chan bool, 10)
  }

	ackchannels := make([]chan bool, num_threads)
  for i := range ressigchannels {
    ackchannels[i] = make(chan bool, 10)
  }

  nodesPerThread := num_nodes/num_threads

  for i := range channels {
		mystart := i*nodesPerThread
		myend := max((i+1)*nodesPerThread, num_nodes)
    go pagerank_func(iterations, W, inlinks, outlinks, mystart, myend,
			channels[i], reschannels[i], sigchannels[i],
			ressigchannels[i], ackchannels[i], redo)
  }

  //fmt.Println("Starting the iterations")
  startTime := time.Now()
  for iter:=0; iter < iterations; iter++{
    fmt.Println("Iteration : ", iter)
    results := make([]float64, num_nodes)
    copy(results, pagerank)
    for i := range channels {
      sigchannels[i] <- true
      pagerankcopy := make([]float64, num_nodes)
      copy(pagerankcopy, pagerank)
      channels[i] <- pagerankcopy
    }
		
    for i := range channels {
			mystart := i*nodesPerThread
			myend := max((i+1)*nodesPerThread, num_nodes)
      pass := <- ressigchannels[i]

			if redo {
				for !pass {
					ackchannels[i] <- false
					pass = <- ressigchannels[i]
				}
				ackchannels[i] <- true			
			}

			if pass {
				tile := <- reschannels[i]			
				copy(results[mystart:myend], tile)
			}
    }
    pagerank = results
  }
	
	end := time.Now()
	elapsed := end.Sub(startTime).Nanoseconds()
	fmt.Println("Elapsed time :", elapsed)

  f, _ := os.Create(outfile)
  defer f.Close()

  for i := range pagerank{
    f.WriteString(fmt.Sprintln(pagerank[i]))
  }
}
