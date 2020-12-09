// to build: go build sort-rec.go
// to run: ./sort-rec ArraySize InputFile OutputFile
// e.g., ./sort-rec 100 input-100.txt output-100.txt

package main

import (
  "bufio"
  "fmt"
  "os"
  "strconv"
  "time"
)

var ArraySize int
var Input, Output []int

// read input file
func readInput() {
  inFile, _ := os.Open(os.Args[2])
  defer inFile.Close()
  scanner := bufio.NewScanner(inFile)
  for i := 0; i < ArraySize; i++ {
    scanner.Scan()
    Input[i], _ = strconv.Atoi(scanner.Text())
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

// selection sort
func selSort() {
  for i := 0; i < ArraySize; i++ {
    Output[i] = Input[i]
  }
  for i := ArraySize-1; i > 0; i-- {
    maxVal := Output[0]
    maxPos := 0
    for j := 1; j <= i; j++ {
      if Output[j] > maxVal {
        maxVal = Output[j]
	maxPos = j
      }
    }
    if maxPos != i {
      temp := Output[maxPos]
      Output[maxPos] = Output[i]
      Output[i] = temp
    }
  }
}

// check that array is sorted
// does NOT guarantee all elements are present in same quantities as input
func checkSort() bool {
  sorted := true
  for i := 1; i < ArraySize; i++ {
    if Output[i] < Output[i-1] {
      sorted = false
      break
    }
  }
  return sorted
}

func main() {
  ArraySize, _ = strconv.Atoi(os.Args[1])

  Input = make([]int, ArraySize)
  Output = make([]int, ArraySize)

  fmt.Println("Reading input file")

  readInput()

  fmt.Println("Starting computation")

  startTime := time.Now()

  // recover from errors in first call
  selSort()
  if !checkSort() {
    selSort()
  }

  elapsed := time.Since(startTime)
  fmt.Println("Elapsed time (nanoseconds):", elapsed.Nanoseconds());

  fmt.Println("Writing output file")

  writeOutput()
}