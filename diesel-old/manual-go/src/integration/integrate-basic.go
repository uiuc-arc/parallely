package main

import "fmt"
import "time"

const Max = float32(11)
const Min = float32(1)
const Delta = float32(0.00390625)
const Threads = 10
const IntervalPerThread = (Max-Min)/float32(Threads)
const DivisionsPerThread = int(IntervalPerThread/Delta)

func calc(x float32) float32 {
  return x*x
}

func integrate(input chan int, output chan float32) {
  sum := float32(0)
  threadId := <- input
  x := Min + IntervalPerThread*float32(threadId)
  for i:=0; i<DivisionsPerThread; i++ {
    sum += calc(x)*Delta
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
