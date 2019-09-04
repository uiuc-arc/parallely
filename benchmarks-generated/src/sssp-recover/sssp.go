package main

import (
  "os"
  "fmt"
  "io/ioutil"
  "strings"
  "math/rand"
  "strconv"
  "time"
  "parallely"
)

func main() {
  rand.Seed(time.Now().UTC().UnixNano())
  argsWithoutProg := os.Args[1:]

  data_bytes, _ := ioutil.ReadFile(argsWithoutProg[0])
  num_nodes, _ := strconv.Atoi(argsWithoutProg[1])
  //outfile := argsWithoutProg[2]
  // num_edges, _ := strconv.Atoi(argsWithoutProg[2])

  //fmt.Println("Starting reading the file")
  data_string := string(data_bytes)
  data_str_array := strings.Split(data_string, "\n")
  
  //fmt.Println("Setting up the data structures")
  W := make([][]int, num_nodes)
  inlinks := make([]int, num_nodes)
  outlinks := make([]int, num_nodes)
  distances := make([]int, num_nodes)
  
  for i := range W{
    W[i] = make([]int, num_nodes)
    inlinks[i] = 0
    outlinks[i] = 0
    distances[i] = 0
  }
  distances[0] = 1

  //fmt.Println("Populating the data structures")
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

  iterations := 10
  overallflag := false

  for i := 0; i <iterations; i++ {
    for node := 0; node < num_nodes; node++ {
      neighbors := inlinks[node]
      dist := distances[node]
      for j := 0; j < neighbors; j++ {
        neighbor := W[node][j]
        ndist := distances[neighbor]
        flag := false
        ndist = parallely.RandchoiceFlag(0.999, ndist, 0, &flag)
        if flag {
          flag = false
          ndist = parallely.RandchoiceFlag(0.9999, ndist, 0, &flag)
        }
        overallflag = overallflag || flag
        if dist > ndist+1 {
          dist = ndist+1
        }
      }
      distances[node] = dist
    }
  }
  if overallflag {
    fmt.Println(1)
  } else {
    fmt.Println(0)
  }

  /*f, _ := os.Create(outfile)
  defer f.Close()
  
  for i := range distances{
    f.WriteString(fmt.Sprintln(distances[i]))  
  }*/
}
