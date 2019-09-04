package main

import (
  "os"
  "fmt"
  "io/ioutil"
  "strings"
  "math/rand"
  "time"
  "strconv"
  "parallely"
)

func main() {
  rand.Seed(time.Now().UTC().UnixNano())
  
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
  
  r := 0.15
  d := 0.85

  flag := false
  overallflag := false

  //fmt.Println("Starting the iterations")
  //startTime := time.Now()
  for iter:=0; iter < iterations; iter++ {
    results := make([]float64, num_nodes)
    copy(results, pagerank)
    for node:=0; node < num_nodes; node++ {
      newpagerank := r
      for k := 0; k<inlinks[node]; k++ {
        neighbor := W[node][k]
        contribution := d * pagerank[neighbor]/float64(outlinks[neighbor])
        flag = false
        contribution = parallely.RandchoiceFlagFloat64(0.999, contribution, 0, &flag)
        if flag {
          flag = false
          contribution = parallely.RandchoiceFlagFloat64(0.9999, contribution, 0, &flag)
        }
        overallflag = overallflag || flag
        newpagerank += contribution
      }
      results[node] = newpagerank
    }
    pagerank = results
  }
  //elapsed := time.Since(startTime)
  //fmt.Println(elapsed)
  if overallflag {
    fmt.Println(1)
  } else {
    fmt.Println(0)
  }

  f, _ := os.Create("output.txt")
  defer f.Close()

  for i := range pagerank{
    f.WriteString(fmt.Sprintln(pagerank[i]))
  }
}
