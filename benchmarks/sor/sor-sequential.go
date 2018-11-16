package main

import (
  "os"
  "strconv"
  "math/rand"
  "time"
  "fmt"
)

func Idx(i, j, width int) int {
  return i*width+j
}

func main() {
  if len(os.Args)!=5 {
    fmt.Println("Usage:\nsor-sequential.go height width iterations omega")
    os.Exit(0)
  }
  height, _ := strconv.Atoi(os.Args[1])
  width, _ := strconv.Atoi(os.Args[2])
  iterations, _ := strconv.Atoi(os.Args[3])
  omega, _ := strconv.ParseFloat(os.Args[4],64)

  randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

  array := make([]float64, height*width)
  for i := 0; i < height*width; i++ {
    array[i] = randGen.Float64()
  }

  startTime := time.Now()

  result := make([]float64, height*width)
  for iter := 0; iter < iterations; iter++ {
    for i := 1; i < height-1; i++ {
      for j := 1; j < width-1; j++ {
        up := array[Idx(i-1,j,width)]
        down := array[Idx(i+1,j,width)]
        left := array[Idx(i,j-1,width)]
        right := array[Idx(i,j+1,width)]
        center := array[Idx(i,j,width)]
        result[Idx(i,j,width)] = omega/4.0*(up+down+left+right) + (1.0-omega)*center
      }
    }
    array = result
  }

  elapsed := time.Since(startTime)
  fmt.Println(elapsed)
}
