package main

import (
  "bufio"
  "os"
  "strconv"
  "fmt"
  "time"
)

func adder(dataChannel chan []int, resultChannel chan int) {
  data := <- dataChannel
  sum := 0
  for i := range data {
    sum += data[i]
  }
  resultChannel <- sum
}

func main() {
  iFile := os.Args[1]
  nums, err := strconv.Atoi(os.Args[2])
  if err != nil {
    fmt.Println(err)
  }
  numThreads, err := strconv.Atoi(os.Args[3])
  if err != nil {
    fmt.Println(err)
  }
  file, err := os.Open(iFile)
  if err != nil {
    fmt.Println(err)
  }
  defer file.Close()
  scanner := bufio.NewScanner(file)
  data := make([]int, nums)
  for i := range data {
    scanner.Scan()
    data[i], err = strconv.Atoi(scanner.Text())
    if err != nil {
      fmt.Println(err)
    }
  }

  dataChannels := make([]chan []int, numThreads)
  resChannels := make([]chan int, numThreads)
  for i := range dataChannels {
		dataChannels[i] = make(chan []int, 1)
		resChannels[i] = make(chan int, 1)
		go adder(dataChannels[i],resChannels[i])
	}
	startTime := time.Now()
	dataPerThread := nums/numThreads
	sum := 0
	for i := range dataChannels {
	  start := i*dataPerThread
	  var end int
	  if i==numThreads-1 {
	    end = nums
	  } else {
	    end = (i+1)*dataPerThread
	  }
	  dataChannels[i] <- data[start:end]
	  partial := <- resChannels[i]
	  sum += partial
	}
	elapsed := time.Since(startTime)
	fmt.Println(sum)
	fmt.Println(elapsed)
}
