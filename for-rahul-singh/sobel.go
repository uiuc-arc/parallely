// to build: go build sobel.go
// to run: ./sobel ArrayDim inputFile outputFile
// e.g., ./sobel 100 input-100x100.txt output-100x100.txt

package main

import (
  "bufio"
  "fmt"
  "os"
  "strconv"
  "time"
)

var ArrayDim, ArraySize int
var Image, Output []int

// read input file
func readInput() {
  inFile, _ := os.Open(os.Args[2])
  defer inFile.Close()
  scanner := bufio.NewScanner(inFile)
  for i := 0; i < ArraySize; i++ {
    scanner.Scan()
    Image[i], _ = strconv.Atoi(scanner.Text())
  }
}

// write output file
func writeOutput() {
  outFile, _ := os.Create(os.Args[3])
  defer outFile.Close()
  for i := 0; i < ArraySize; i++ {
    _, _ = outFile.WriteString(fmt.Sprintf("%d\n", Output[i]))
  }
}

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
  ArrayDim, _ = strconv.Atoi(os.Args[1])
  ArraySize = ArrayDim*ArrayDim

  Image = make([]int, ArraySize)
  Output = make([]int, ArraySize)

  fmt.Println("Reading input file")

  readInput()

  fmt.Println("Starting computation")

  startTime := time.Now()

  sobel()

  elapsed := time.Since(startTime)
  fmt.Println("Elapsed time (nanoseconds):", elapsed.Nanoseconds());

  fmt.Println("Writing output file")

  writeOutput()
}