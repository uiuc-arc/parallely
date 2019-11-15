package main

import "fmt"
import "time"

func integrate(input chan int, output chan float32) {
  sum := float32(0)
  threadId := <- input
  x := Min + IntervalPerThread*float32(threadId)
  for i:=0; i<DivisionsPerThread; i++ {
    sum += x*x*Delta
    x += Delta
  }
  output <- sum
}

func main() {
  inputChans := make([]chan int, Threads)
  for i := range inputChans {
    inputChans[i] = make(chan int, 1)
  }
  outputChans := make([]chan float32, Threads)
  for i := range inputChans {
    outputChans[i] = make(chan float32, 1)
  }

  for i:=0; i<Threads; i++ {
    go integrate(inputChans[i], outputChans[i])
  }

  startTime := time.Now()

  for i:=0; i<Threads; i++ {
    inputChans[i] <- i
  }
  sum := float32(0)
  for i:=0; i<Threads; i++ {
    temp := <- outputChans[i]
    sum += temp
  }

  elapsed := time.Since(startTime)
  fmt.Println(elapsed)
}
