package main

import (
  "os"
  "fmt"
  "io/ioutil"
  "strings"
  "time"
  "strconv"
	"runtime/pprof"
	"log"
	"runtime"
	."dynfloats"
	// "encoding/json"
)


// var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
// var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxArray(a []int) int {
	max := a[0]
	for x := range a {
		if a[x] > max {
			max = x
		}
	}
	return max
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var dynamic_error map[string]int

func pagerank_func(iterations int, W [][]int, inlinks []int, outlinks []int,
	myfirstnode, mylastnode int, inchannel, outchannel chan float64,
	// inchannel_reliability, outchannel_reliability chan []float32,
	datasigchannel_in chan bool, datasigchannel_out chan bool, size int){

  r := 0.15
  d := 0.85

	pageranks := make([]float64, size)

  for myiteration := 0; myiteration < iterations; myiteration++{
    // <- datasigchannel_in

		// fmt.Println(myfirstnode, mylastnode, size)
    for j := 0; j < size; j++ {
    	pageranks[j] = <- inchannel
    }
		// fmt.Println("Got data")

    mypageranks := make([]float64, mylastnode-myfirstnode)

    for node := range mypageranks{
      mypageranks[node] = r
      for k := 0; k<inlinks[myfirstnode+node]; k++ {
        neighbor := W[myfirstnode+node][k]

				temp0 := float64(outlinks[neighbor])
				temp1 := d * pageranks[neighbor]
				temp2 := temp1 / temp0
				
        mypageranks[node] = temp2 + mypageranks[node]
      }
    }

		for node := range mypageranks{
			DynSendFloat(datasigchannel_out, outchannel, mypageranks[node], 0.00001)
		}
  }
}


func main() {
  data_bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	
  num_nodes, _ := strconv.Atoi(os.Args[2])
  iterations, _ := strconv.Atoi(os.Args[3])
  num_threads, _ := strconv.Atoi(os.Args[4])
	debug, _ := strconv.Atoi(os.Args[5])
	outfile := os.Args[6]

	// failed_once := false
	
	// var cpuprofile = flag.String("cpuprofile", "", "profile cpu")
	// var memprofile = flag.String("memprofile", "", "profile mem")
	// var ip = flag.String("flagname", 1234, "help message for flagname")

	// var redo bool = redoin != 0

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

	fmt.Printf("Max number of degrees : %d\n", maxArray(inlinks))

  //fmt.Println("Finished populating the data structures")

  channels_main_threads := make([]chan float64, num_threads)
  for i := range channels_main_threads {
    channels_main_threads[i] = make(chan float64, 1)
  }
  channels_threads_main := make([]chan float64, num_threads)
  for i := range channels_threads_main {
    channels_threads_main[i] = make(chan float64, 1)
  }

	
	// dynamic_channels_main_threads := make([]chan DynRelyFloat, num_threads)
  // for i := range dynamic_channels_main_threads {
  //   dynamic_channels_main_threads[i] = make(chan DynRelyFloat, 10)
  // }
	// dynamic_channels_threads_main := make([]chan []float32, num_threads)
  // for i := range dynamic_channels_main_threads {
  //   dynamic_channels_threads_main[i] = make(chan []float32, 10)
  // }	
	
  sigchannels_in := make([]chan bool, num_threads)
  for i := range sigchannels_in {
    sigchannels_in[i] = make(chan bool, 1)
  }
	sigchannels_out := make([]chan bool, num_threads)
  for i := range sigchannels_out {
		sigchannels_out[i] = make(chan bool, 1)
  }

		if debug != 0 {
		fmt.Println("Start profiling")
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

  nodesPerThread := num_nodes/num_threads

  for i := range channels_main_threads {
		mystart := i*nodesPerThread
		myend := (i+1)*nodesPerThread
		if i == num_threads - 1 {
			myend = max((i+1)*nodesPerThread, num_nodes)
		}
		go pagerank_func(iterations, W, inlinks, outlinks, mystart, myend,
			channels_main_threads[i], channels_threads_main[i],
			// dynamic_channels_main_threads[i], dynamic_channels_threads_main[i],
			sigchannels_in[i], sigchannels_out[i], num_nodes)
	}

  //fmt.Println("Starting the iterations")
  startTime := time.Now()
	
  for iter:=0; iter < iterations; iter++{
    fmt.Println("Iteration : ", iter)
    // results := make([]DynRelyFloat, num_nodes)
    // copy(results, pagerank)
		
    for i := range channels_main_threads {
      for j := 0; j < num_nodes; j++ {
      	channels_main_threads[i] <- pagerank[j]
      }
    }
		
    for i := range channels_main_threads {
			mystart := i*nodesPerThread
			myend := (i+1)*nodesPerThread
			if i == num_threads -1 {
				myend = max((i+1)*nodesPerThread, num_nodes)
			}

			for k := mystart; k < myend; k++ {			
				pass := <- sigchannels_out[i]
				if pass {
					pagerank[k] = <- channels_threads_main[i]
				}
				// else {
				// 	failed_once = true
				// }
			}
    }

		// f, _ := os.Create(fmt.Sprintf("_iter_%d.txt", iter))
		// f.WriteString(fmt.Sprintln(failed_once))

		// fileJson, _ := json.MarshalIndent(pagerank, "", " ")
		// f.WriteString(string(fileJson))

		// for i := range pagerank{

		// 	f.WriteString(fmt.Sprintln(pagerank[i]))
		// }
		// f.Close()
			
  }
	
	end := time.Now()
	elapsed := end.Sub(startTime).Nanoseconds()
	fmt.Println("Elapsed time :", elapsed)

	if debug != 0 {
		f, err := os.Create("mem.prof")
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
  f, _ := os.Create(outfile)
  defer f.Close()

  for i := range pagerank{
    f.WriteString(fmt.Sprintln(pagerank[i]))
  }
}
