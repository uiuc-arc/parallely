// to build: go build sobel.go
// to run: go run sobel.go  OR  ./sobel (only after building)

package main

import (
  "fmt"
  "math/rand"
  "time"
)

// modify to change input size
const ArrayDim = 100

// do not modify this - based on above
const ArraySize = ArrayDim*ArrayDim
var Image,Output [ArraySize]int

// index calculator for 2D arrays
func idx(i, j int) int {
  return i*ArrayDim+j
}

// sobel kernel
func sobel() {
  for i := 1; i < ArrayDim-1; i++ {
    for j := 1; j < ArrayDim-1; j++ {
      etl := Image[idx(i-1,j-1)]
      etc := Image[idx(i-1,j)]
      etr := Image[idx(i-1,j+1)]
      ell := Image[idx(i+1,j-1)]
      elc := Image[idx(i+1,j)]
      elr := Image[idx(i+1,j+1)]
      point := etl+2*etc+etr-ell-2*elc-elr
      Output[idx(i,j)] = point
    }
  }
}

func main() {
  fmt.Println("Generating random input of size",ArraySize)

  rand.Seed(time.Now().UTC().UnixNano())
  for i := 0; i < ArraySize; i++ {
    Image[i] = rand.Intn(256)
  }

  fmt.Println("Starting computation")

  startTime := time.Now()

  sobel()

  elapsed := time.Since(startTime)
  fmt.Println("Elapsed time (nanoseconds):", elapsed.Nanoseconds());
}